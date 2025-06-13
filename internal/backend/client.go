package backend

import "internal/item"

type Client interface {
	Connect() error
	Index()
	LookupItem(string) (item.Item, error)
	UpsertItem(item.Item) error
}
