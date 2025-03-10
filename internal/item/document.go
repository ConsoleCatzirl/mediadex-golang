package item

import "github.com/junlicn/yami"

type JsonDocument struct {
	FileStats  *FileItem   `json:"fileinfo"`
	MediaStats *yami.Media `json:"mediainfo"` // todo: *yami.Media

	DiscogsInfo interface{} `json:"discogs"` // todo: upstream type
	ImdbInfo    interface{} `json:"imdb"`    // todo: upstream type
	TvdbInfo    interface{} `json:"tvdb"`    // todo: upstream type
}
