package item

type baseItem struct {
	// the struct sent to the backend must not have any scoped functions
	jsonData JsonDocument
}

func (i *baseItem) JsonDoc() (*JsonDocument, error) {
	return &i.jsonData, nil
}

func (i *baseItem) addMediaInfo() error {
	return nil
}
