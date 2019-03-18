package xid_test

import (
	"testing"

	gu "github.com/google/uuid"
	"github.com/rs/xid"
	su "github.com/satori/go.uuid"
)

func BenchmarkXIDNew(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = xid.New()
		}
	})
}

func BenchmarkXIDNewString(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = xid.New().String()
		}
	})
}

func BenchmarkXIDFromString(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = xid.FromString("9m4e2mr0ui3e8a215n4g")
		}
	})
}

func BenchmarkSatoriUUIDV1(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = su.NewV1().String()
		}
	})
}

func BenchmarkSatoriUUIDV4(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = su.NewV4().String()
		}
	})
}

func BenchmarkGoogleUUID(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = gu.New()
		}
	})
}
