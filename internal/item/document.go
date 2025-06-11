package item

import "github.com/junlicn/yami"

type JsonDocument struct {
	FileStats  *FileItem   `json:"fileinfo"`
	MediaStats *yami.Media `json:"mediainfo"` // todo: *yami.Media

	DiscogsInfo interface{} `json:"discogs"` // todo: upstream type
	ImdbInfo    interface{} `json:"imdb"`    // todo: upstream type
	TvdbInfo    interface{} `json:"tvdb"`    // todo: upstream type
}

type arangoDocument struct {
	*JsonDocument
	Key string `json:"_key"`
}

func (j *JsonDocument) Arango() *arangoDocument {
	newDoc := &arangoDocument{
		JsonDocument: j,
		Key:          j.FileStats.Checksum,
	}
	return newDoc
}

type openSearchDocument struct {
	*JsonDocument
	ID string `json:"_id"`
}

func (j *JsonDocument) OpenSearch() *openSearchDocument {
	newDoc := &openSearchDocument{
		JsonDocument: j,
		ID:           j.FileStats.Checksum,
	}
	return newDoc
}

func makeJsonDoc(fItem *FileItem) *JsonDocument {
	return &JsonDocument{
		FileStats: fItem,
	}
}
