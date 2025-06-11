module github.com/ConsoleCatzirl/mediadex-golang

go 1.23.0

toolchain go1.24.3

replace (
	internal/backend => ./internal/backend
	internal/cli => ./internal/cli
	internal/item => ./internal/item
	internal/runner => ./internal/runner
	internal/walker => ./internal/walker
)

replace (
	pkg/conf => ./pkg/conf
	pkg/worker => ./pkg/worker
)

require (
	internal/backend v0.0.0-00010101000000-000000000000 // indirect
	internal/cli v0.0.0-00010101000000-000000000000
	internal/item v0.0.0-00010101000000-000000000000 // indirect
	internal/runner v0.0.0-00010101000000-000000000000 // indirect
	internal/walker v0.0.0-00010101000000-000000000000 // indirect
)

require (
	pkg/conf v0.0.0-00010101000000-000000000000 // indirect
	pkg/worker v0.0.0-00010101000000-000000000000
)

require (
	github.com/arangodb/go-driver/v2 v2.1.3 // indirect
	github.com/arangodb/go-velocypack v0.0.0-20200318135517-5af53c29c67e // indirect
	github.com/dchest/siphash v1.2.3 // indirect
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/junlicn/yami v0.0.0-20241205082406-718a5624f3fa // indirect
	github.com/kkdai/maglev v0.2.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/opensearch-project/opensearch-go/v4 v4.5.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rs/zerolog v1.33.0 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
