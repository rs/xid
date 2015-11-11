// Package xid is a globally unique id generator suited for web scale
//
// Xid is using Mongo Object ID algorithm to generate globally unique ids:
// https://docs.mongodb.org/manual/reference/object-id/
//
//   - 4-byte value representing the seconds since the Unix epoch,
//   - 3-byte machine identifier,
//   - 2-byte process id, and
//   - 3-byte counter, starting with a random value.
//
// The binary representation of the id is compatible with Mongo 12 bytes Object IDs.
// The string representation is using URL safe base64 for better space efficiency when
// stored in that form (16 bytes).
//
// UUID is 16 bytes (128 bits), snowflake is 8 bytes (64 bits), xid stands in between
// with 12 bytes with a more compact string representation ready for the web and no
// required configuration or central generation server.
//
// Features:
//
//   - Size: 12 bytes (96 bits), smaller than UUID, larger than snowflake
//   - Base64 URL safe encoded by default (16 bytes storage when transported as printable string)
//   - Non configured, you don't need set a unique machine and/or data center id
//   - K-ordered
//   - Embedded time with 1 second precision
//   - Unicity guaranted for 16,777,216 (24 bits) unique ids per second and per host/process
//
// Best used with xlog's RequestIDHandler (https://godoc.org/github.com/rs/xlog#RequestIDHandler).
//
// References:
//
//   - http://www.slideshare.net/davegardnerisme/unique-id-generation-in-distributed-systems
//   - https://en.wikipedia.org/wiki/Universally_unique_identifier
//   - https://blog.twitter.com/2010/announcing-snowflake
package xid

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"
)

// Code inspired from mgo/bson ObjectId

// ID represents a unique request id
type ID [rawLen]byte

const (
	encodedLen = 16
	rawLen     = 12
)

// ErrInvalidID is returned when trying to unmarshal an invalid ID
var ErrInvalidID = errors.New("invalid ID")

// objectIDCounter is atomically incremented when generating a new ObjectId
// using NewObjectId() function. It's used as a counter part of an id.
var objectIDCounter uint32

// init objectIDCounter to be a random initial value.
func init() {
	b := make([]byte, 3)
	if _, e := io.ReadFull(rand.Reader, b); e != nil {
		panic(e)
	}
	objectIDCounter = uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
}

// machineId stores machine id generated once and used in subsequent calls
// to NewObjectId function.
var machineID = readMachineID()

// readMachineId generates machine id and puts it into the machineId global
// variable. If this function fails to get the hostname, it will cause
// a runtime error.
func readMachineID() []byte {
	id := make([]byte, 3)
	hostname, err1 := os.Hostname()
	if err1 != nil {
		// Fallback to rand number if machine id can't be gathered
		_, err2 := io.ReadFull(rand.Reader, id)
		if err2 != nil {
			panic(fmt.Errorf("cannot get hostname: %v; %v", err1, err2))
		}
		return id
	}
	hw := md5.New()
	hw.Write([]byte(hostname))
	copy(id, hw.Sum(nil))
	return id
}

// New generates a globaly unique ID
func New() ID {
	var id ID
	// Timestamp, 4 bytes, big endian
	binary.BigEndian.PutUint32(id[:], uint32(time.Now().Unix()))
	// Machine, first 3 bytes of md5(hostname)
	id[4] = machineID[0]
	id[5] = machineID[1]
	id[6] = machineID[2]
	// Pid, 2 bytes, specs don't specify endianness, but we use big endian.
	pid := os.Getpid()
	id[7] = byte(pid >> 8)
	id[8] = byte(pid)
	// Increment, 3 bytes, big endian
	i := atomic.AddUint32(&objectIDCounter, 1)
	id[9] = byte(i >> 16)
	id[10] = byte(i >> 8)
	id[11] = byte(i)
	return id
}

// String returns a base64 URL safe representation of the id
func (id ID) String() string {
	return base64.URLEncoding.EncodeToString(id[:])
}

// MarshalText implements encoding/text TextMarshaler interface
func (id ID) MarshalText() (text []byte, err error) {
	text = make([]byte, encodedLen)
	base64.URLEncoding.Encode(text, id[:])
	return
}

// UnmarshalText implements encoding/text TextUnmarshaler interface
func (id *ID) UnmarshalText(text []byte) error {
	if len(text) != encodedLen {
		return ErrInvalidID
	}
	b := make([]byte, rawLen)
	_, err := base64.URLEncoding.Decode(b, text)
	for i, c := range b {
		id[i] = c
	}
	return err
}

// Time returns the timestamp part of the id.
// It's a runtime error to call this method with an invalid id.
func (id ID) Time() time.Time {
	// First 4 bytes of ObjectId is 32-bit big-endian seconds from epoch.
	secs := int64(binary.BigEndian.Uint32(id[0:4]))
	return time.Unix(secs, 0)
}

// Machine returns the 3-byte machine id part of the id.
// It's a runtime error to call this method with an invalid id.
func (id ID) Machine() []byte {
	return id[4:7]
}

// Pid returns the process id part of the id.
// It's a runtime error to call this method with an invalid id.
func (id ID) Pid() uint16 {
	return binary.BigEndian.Uint16(id[7:9])
}

// Counter returns the incrementing value part of the id.
// It's a runtime error to call this method with an invalid id.
func (id ID) Counter() int32 {
	b := id[9:12]
	// Counter is stored as big-endian 3-byte value
	return int32(uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2]))
}
