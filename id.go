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
// The string representation is using base32 hex (w/o padding) for better space efficiency
// when stored in that form (20 bytes). The hex variant of base32 is used to retain the
// sortable property of the id.
//
// Xid doesn't use base64 because case sensitivity and the 2 non alphanum chars may be an
// issue when transported as a string between various systems. Base36 wasn't retained either
// because 1/ it's not standard 2/ the resulting size is not predictable (not bit aligned)
// and 3/ it would not remain sortable. To validate a base32 `xid`, expect a 20 chars long,
// all lowercase sequence of `a` to `v` letters and `0` to `9` numbers (`[0-9a-v]{20}`).
//
// UUID is 16 bytes (128 bits), snowflake is 8 bytes (64 bits), xid stands in between
// with 12 bytes with a more compact string representation ready for the web and no
// required configuration or central generation server.
//
// Features:
//
//   - Size: 12 bytes (96 bits), smaller than UUID, larger than snowflake
//   - Base32 hex encoded by default (16 bytes storage when transported as printable string)
//   - Non configured, you don't need set a unique machine and/or data center id
//   - K-ordered
//   - Embedded time with 1 second precision
//   - Unicity guaranted for 16,777,216 (20 bits) unique ids per second and per host/process
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
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

// Code inspired from mgo/bson ObjectId

// ID represents a unique request id
type ID [rawLen]byte

const (
	trimLen    = 20 // len after padding removal
	encodedLen = 24 // len after base32 encoding, with padding
	decodedLen = 15 // len after base32 decoding with the padded data
	rawLen     = 12 // binary raw len
)

// ErrInvalidID is returned when trying to unmarshal an invalid ID
var ErrInvalidID = errors.New("invalid ID")

// objectIDCounter is atomically incremented when generating a new ObjectId
// using NewObjectId() function. It's used as a counter part of an id.
// This id is initialized with a random value.
var objectIDCounter = randInt()

// machineId stores machine id generated once and used in subsequent calls
// to NewObjectId function.
var machineID = readMachineID()

// readMachineId generates machine id and puts it into the machineId global
// variable. If this function fails to get the hostname, it will cause
// a runtime error.
func readMachineID() []byte {
	id := make([]byte, 3)
	if hostname, err := os.Hostname(); err == nil {
		hw := md5.New()
		hw.Write([]byte(hostname))
		copy(id, hw.Sum(nil))
	} else {
		// Fallback to rand number if machine id can't be gathered
		if _, randErr := rand.Reader.Read(id); randErr != nil {
			panic(fmt.Errorf("Cannot get hostname nor generate a random number: %v; %v", err, randErr))
		}
	}
	return id
}

// randInt generates a random uint32
func randInt() uint32 {
	b := make([]byte, 3)
	if _, err := rand.Reader.Read(b); err != nil {
		panic(fmt.Errorf("Cannot generate random number: %v;", err))
	}
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
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

// FromString reads an ID from its string representation
func FromString(id string) (ID, error) {
	i := &ID{}
	err := i.UnmarshalText([]byte(id))
	return *i, err
}

// String returns a base32 hex lowercased with no padding representation of the id (char set is 0-9, a-v).
func (id ID) String() string {
	text, _ := id.MarshalText()
	return string(text)
}

// MarshalText implements encoding/text TextMarshaler interface
func (id ID) MarshalText() ([]byte, error) {
	text := make([]byte, encodedLen)
	base32.HexEncoding.Encode(text, id[:])
	return bytes.ToLower(text[:trimLen]), nil
}

// UnmarshalText implements encoding/text TextUnmarshaler interface
func (id *ID) UnmarshalText(text []byte) error {
	if len(text) != trimLen {
		return ErrInvalidID
	}
	for _, c := range text {
		if !(c >= '0' && c <= '9') && !(c >= 'a' && c <= 'v') {
			return ErrInvalidID
		}
	}
	b := make([]byte, decodedLen)
	_, err := base32.HexEncoding.Decode(b, append(bytes.ToUpper(text), '=', '=', '=', '='))
	for i, c := range b {
		id[i] = c
		// The decoded len is larger than the actual len because of padding.
		// Stop copying data when we reach raw len.
		if i+1 == rawLen {
			break
		}
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
