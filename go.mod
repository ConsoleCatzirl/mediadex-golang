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

require pkg/worker v0.0.0-00010101000000-000000000000

require (
	github.com/kr/pretty v0.2.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	pkg/conf v0.0.0-00010101000000-000000000000 // indirect
)

require (
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/junlicn/yami v0.0.0-20241205082406-718a5624f3fa // indirect
	golang.org/x/net v0.39.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
