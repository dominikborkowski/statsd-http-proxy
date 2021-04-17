# build tools
IS_GCCGO_INSTALLED=$(gccgo --version 2> /dev/null)

# build version
VERSION=`git describe --tags | awk -F'-' '{print $$1}'`
BUILD_NUMBER=`git rev-parse HEAD`
BUILD_DATE=`date +%Y-%m-%d-%H:%M`

# go compiler flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildNumber=$(BUILD_NUMBER) -X main.BuildDate=$(BUILD_DATE)"
LDFLAGS_COMPRESSED=-ldflags "-s -w -X main.Version=$(VERSION) -X main.BuildNumber=$(BUILD_NUMBER) -X main.BuildDate=$(BUILD_DATE)"

#gccgo compiler flags
GCCGOFLAGS=-gccgoflags "-march=native -O3"
GCCGOFLAGS_GOLD=-gccgoflags "-march=native -O3 -fuse-ld=gold"

# default task
default: build

deps:
    ifneq ($(GO111MODULE),on)
		export GOPATH=$(CURDIR)
		go get -v -t -d ./...
    endif

deps-gccgo: deps
    ifndef IS_GCCGO_INSTALLED
		$(error "gccgo not installed")
    endif

test: deps
	go test -cover ./...

goveralls: deps
	go get github.com/mattn/goveralls
	$(GOPATH)/bin/goveralls -service=travis-ci

# build with go compiler
build: deps
	ls $(GOPATH)
	CGO_ENABLED=0 go build -v -x -a $(LDFLAGS) -o $(CURDIR)/bin/statsd-http-proxy

# build with go compiler and link optiomizations
build-shrink: deps
	CGO_ENABLED=0 go build -v -x -a $(LDFLAGS_COMPRESSED) -o $(CURDIR)/bin/statsd-http-proxy-shrink

# build with gccgo compiler
# Require to install gccgo
build-gccgo: deps-gccgo
	CGO_ENABLED=0 go build -v -x -a -compiler gccgo $(GCCGOFLAGS) -o $(CURDIR)/bin/statsd-http-proxy-gccgo

# build with gccgo compiler and gold linker
# Require to install gccgo
build-gccgo-gold: deps-gccgo
	CGO_ENABLED=0 go build -v -x -a -compiler gccgo $(GCCGOFLAGS_GOLD) -o $(CURDIR)/bin/statsd-http-proxy-gccgo-gold

# build all
build-all: build build-shrink build-gccgo build-gccgo-gold

# clean build
clean:
	rm -rf ./bin
	go clean

# Install binary locally
install:
	cp $(CURDIR)/bin/statsd-http-proxy /usr/local/bin

# run statsd proxy in http mode
run-http:
	GOMAXPROCS=1 go run main.go --verbose --http-host=127.0.0.1 --http-port=8080 --statsd-host=127.0.0.1 --statsd-port=8125 --metric-prefix=prefix.subprefix --jwt-secret=somesecret

# run statsd proxy in http mode with profiling
run-http-prof:
	GOMAXPROCS=1 go run main.go --verbose --http-host=127.0.0.1 --http-port=8080 --statsd-host=127.0.0.1 --statsd-port=8125 --metric-prefix=prefix.subprefix --jwt-secret=somesecret -profiler-http-port=6060

# run statsd mock that listen UPD port (for monitoring that proxy sends metrics to UPD)
run-statsd-mock:
	nc -kluv localhost 8125

# show profiler results in cli
run-profiler-cli:
	go tool pprof http://localhost:6060/debug/pprof/profile

# show profiler results in web
run-profiler-web:
	go tool pprof -http=localhost:6061 http://localhost:6060/debug/pprof/profile

# show web interface with trace info
# project must be run in profiler mode through "nake run-http-prof"
run-trace:
	curl "http://localhost:6060/debug/pprof/trace" > /tmp/statsd-http-proxt-trace
	go tool trace /tmp/statsd-http-proxt-trace

run-siege:
	time siege -c 150 -r 2000 -H 'Connection: close' -H 'X-JWT-Token:eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJzdGF0c2QtcmVzdC1zZXJ2ZXIiLCJpYXQiOjE1MDY5NzI1ODAsImV4cCI6MTg4NTY2Mzc4MCwiYXVkIjoiaHR0cHM6Ly9naXRodWIuY29tL3Nva2lsL3N0YXRzZC1yZXN0LXNlcnZlciIsInN1YiI6InNva2lsIn0.sOb0ccRBnN1u9IP2jhJrcNod14G5t-jMHNb_fsWov5c' "http://127.0.0.1:8080/count/a.b.c.d POST value=42"

run-wrk:
	docker run --rm --network="host" -v $(CURDIR):/proxy williamyeh/wrk -c 1 -t 1 -d 20 -s /proxy/bench/wrk.lua http://127.0.0.1:8080/count/a.b.c.d