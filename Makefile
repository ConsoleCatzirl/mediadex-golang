goFiles := $(wildcard pkg/*/*.go) $(wildcard internal/*/*.go)
modFiles := $(wildcard pkg/*/go.mod) $(wildcard internal/*/go.mod)

mediadex: cmd/mediadex.go $(goFiles) $(modFiles)
	go build $<


.PHONY: all clean mod vet test online-test
all: mod test mediadex

clean:
	rm mediadex

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
	go vet test/fakes
	go vet test/tconf
	go vet test/tworker

test: vet
	go test -v internal/item internal/runner internal/walker
	go test -v pkg/worker test/tconf


localhost-test:
	go test -v test/tworker


.PHONY: docker-start docker-stop docker-clean

docker-start:
	docker compose -p mediadex-test -f test/docker/compose.yaml up -d

docker-stop:
	docker compose -p mediadex-test -f test/docker/compose.yaml stop

docker-clean:
	docker compose -p mediadex-test -f test/docker/compose.yaml down -v
