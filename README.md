# Globally Unique ID Generator

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/rs/xid) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/rs/xid/master/LICENSE) [![Build Status](https://travis-ci.org/rs/xid.svg?branch=master)](https://travis-ci.org/rs/xid) [![Coverage](http://gocover.io/_badge/github.com/rs/xid)](http://gocover.io/github.com/rs/xid)

Package xid is a globally unique id generator suited for web scale

Xid is using Mongo Object ID algorithm to generate globally unique ids:
https://docs.mongodb.org/manual/reference/object-id/

- 4-byte value representing the seconds since the Unix epoch,
- 3-byte machine identifier,
- 2-byte process id, and
- 3-byte counter, starting with a random value.

The binary representation of the id is compatible with Mongo 12 bytes Object IDs.
The string representation is using URL safe base64 for better space efficiency when
stored in that form (16 bytes).

UUID is 16 bytes (128 bits), snowflake is 8 bytes (64 bits), xid stands in between
with 12 bytes with a more compact string representation ready for the web and no
required configuration or central generation server.

Features:

- Size: 12 bytes (96 bits), smaller than UUID, larger than snowflake
- Base64 URL safe encoded by default (16 bytes storage when transported as printable string)
- Non configured, you don't need set a unique machine and/or data center id
- K-ordered
- Embedded time with 1 second precision
- Unicity guaranted for 16,777,216 (24 bits) unique ids per second and per host/process

Best used with [xlog](https://github.com/rs/xlog)'s
[RequestIDHandler](https://godoc.org/github.com/rs/xlog#RequestIDHandler).

References:

- http://www.slideshare.net/davegardnerisme/unique-id-generation-in-distributed-systems
- https://en.wikipedia.org/wiki/Universally_unique_identifier
- https://blog.twitter.com/2010/announcing-snowflake

## Install

    go get github.com/rs/xid

## Usage

```go
guid := xid.New()

println(guid.String())
// Output: TYjhW2D0huQoQS3J
```

Get `xid` embedded info:

```go
guid.Machine()
guid.Pid()
guid.Time()
guid.Counter()
```

## Licenses

All source code is licensed under the [MIT License](https://raw.github.com/rs/xid/master/LICENSE).
