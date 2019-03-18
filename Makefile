GO_VERSION:=$(shell go version)

.PHONY: all clean bench bench-all profile lint test contributors update install

all: clean install lint test bench

clean:
	rm -rf ./*.svg
	rm -rf bench
	rm -rf pprof

bench: clean
	go test -count=5 -run=NONE -bench . -benchmem

init:
	GO111MODULE=on go mod init
	GO111MODULE=on go mod vendor

profile: clean
	mkdir bench
	mkdir pprof
	go test -count=10 -run=NONE -bench . -benchmem -o pprof/test.bin -cpuprofile pprof/cpu.out -memprofile pprof/mem.out
	go tool pprof --svg pprof/test.bin pprof/mem.out > bench/mem.svg
	go tool pprof --svg pprof/test.bin pprof/cpu.out > bench/cpu.svg
	go test -count=10 -run=NONE -bench=BenchmarkXIDNewString -benchmem -o pprof/xid-test.bin -cpuprofile pprof/cpu-xid.out -memprofile pprof/mem-xid.out
	go tool pprof --svg pprof/xid-test.bin pprof/mem-xid.out > bench/mem-xid.svg
	go tool pprof --svg pprof/xid-test.bin pprof/cpu-xid.out > bench/cpu-xid.svg
	go-torch -f bench/cpu-xid-graph.svg pprof/xid-test.bin pprof/cpu-xid.out
	go-torch --alloc_objects -f bench/mem-xid-graph.svg pprof/xid-test.bin pprof/mem-xid.out

test: clean init
	go test --race -v $(go list ./... | rg -v vendor)

contributors:
	git log --format='%aN <%aE>' | sort -fu > CONTRIBUTORS
