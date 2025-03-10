package item

type Item interface {
	AddMetadata() error
	JsonDoc() (*JsonDocument, error)
}
