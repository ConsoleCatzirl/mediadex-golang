package backend

import "internal/item"

type Client interface {
	Index() error
	LookupItem(string) (item.Item, error)
}
