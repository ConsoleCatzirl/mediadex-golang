package backend

import "internal/item"

type Client interface {
	Connect() error
	Index() error
	LookupItem(string) (item.Item, error)
}
