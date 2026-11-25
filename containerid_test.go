package xid

import (
	"hash/crc32"
	"os"
	"path/filepath"
	"testing"
)

func TestPIDWithContainerIDFallsBackToCgroup(t *testing.T) {
	dir := t.TempDir()
	cpuset := filepath.Join(dir, "missing-cpuset")
	cgroup := filepath.Join(dir, "cgroup")
	if err := os.WriteFile(cgroup, []byte("0::/docker/container\n"), 0644); err != nil {
		t.Fatal(err)
	}

	got := pidWithContainerID(cpuset, cgroup)
	want := os.Getpid() ^ int(crc32.ChecksumIEEE([]byte("0::/docker/container\n")))
	if got != want {
		t.Fatalf("expected pid %d, got %d", want, got)
	}
}

func TestPIDWithContainerIDPrefersCpuset(t *testing.T) {
	dir := t.TempDir()
	cpuset := filepath.Join(dir, "cpuset")
	cgroup := filepath.Join(dir, "cgroup")
	if err := os.WriteFile(cpuset, []byte("/docker/cpuset\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cgroup, []byte("0::/docker/cgroup\n"), 0644); err != nil {
		t.Fatal(err)
	}

	got := pidWithContainerID(cpuset, cgroup)
	want := os.Getpid() ^ int(crc32.ChecksumIEEE([]byte("/docker/cpuset\n")))
	if got != want {
		t.Fatalf("expected pid %d, got %d", want, got)
	}
}

func TestPIDWithContainerIDWithoutContainerFiles(t *testing.T) {
	dir := t.TempDir()
	got := pidWithContainerID(filepath.Join(dir, "missing-cpuset"), filepath.Join(dir, "missing-cgroup"))
	if got != os.Getpid() {
		t.Fatalf("expected unmodified pid %d, got %d", os.Getpid(), got)
	}
}
