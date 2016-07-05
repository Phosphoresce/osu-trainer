GOC=go build
GOFLAGS=-a -ldflags '-s'
CGOR=CGO_ENABLED=0

all: build

build:
	$(GOC) ot.go

run:
	go run ot.go

stat:
	$(CGOR) $(GOC) $(GOFLAGS) ot.go

fmt:
	gofmt -w .

clean:
	rm ot
