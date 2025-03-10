goFiles := $(wildcard pkg/*/*.go) $(wildcard internal/*/*.go)
modFiles := $(wildcard pkg/*/go.mod) $(wildcard internal/*/go.mod)

mediadex: cmd/mediadex.go $(goFiles) $(modFiles)
	go build $<

.PHONY: all mod vet test clean
all: mod test mediadex

mod:
	go mod tidy

vet:
	go vet pkg/conf
	go vet pkg/worker
	go vet internal/backend
	go vet internal/cli
	go vet internal/item
	go vet internal/runner
	go vet internal/walker
	go vet ./test

test: vet
	go test -v ./test

clean:
	rm mediadex
