package item

type baseItem struct {
	jsonData *JsonDocument
}

func (i *baseItem) JsonDoc() (*JsonDocument, error) {
	return i.jsonData, nil
}

func (i *baseItem) addMediaInfo() error {
	return nil
}
