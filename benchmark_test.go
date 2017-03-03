package xid_test

import (
	"testing"

	"github.com/rs/xid"
	uuid "github.com/satori/go.uuid"
)

func BenchmarkXID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = xid.New().String()
	}
}

func BenchmarkUUIDv1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = uuid.NewV1().String()
	}
}

func BenchmarkUUIDv4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = uuid.NewV4().String()
	}
}
