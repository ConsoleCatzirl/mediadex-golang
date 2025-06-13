package item

import (
	"errors"
	"time"

	"github.com/junlicn/yami"
)

type baseItem struct {
	// the struct sent to the backend must not have any scoped functions
	jsonData *JsonDocument
}

func (i *baseItem) JsonDoc() (*JsonDocument, error) {
	return i.jsonData, nil
}

func (i *baseItem) addMediaInfo() error {
	mInfo, err := yami.GetMediaInfo(i.jsonData.FileStats.FullPath, time.Minute)
	if err != nil {
		return err
	}
	if mInfo.Media == nil {
		return errors.New("No mediainfo data")
	}

	i.jsonData.MediaStats = mInfo.Media
	return nil
}
