package item

import "github.com/junlicn/yami"

type JsonDocument struct {
	FileStats  *FileItem   `json:"fileinfo"`
	MediaStats *yami.Media `json:"mediainfo"`

	DiscogsInfo interface{} `json:"discogs"` // todo: upstream type
	ImdbInfo    interface{} `json:"imdb"`    // todo: upstream type
	TvdbInfo    interface{} `json:"tvdb"`    // todo: upstream type
}

type ArangoDocument struct {
	*JsonDocument
	Key string `json:"_key"`
}

func (j *JsonDocument) Arango() *ArangoDocument {
	newDoc := &ArangoDocument{
		JsonDocument: j,
		Key:          j.FileStats.Checksum,
	}
	return newDoc
}

type OpenSearchDocument struct {
	*JsonDocument
	ID string `json:"_id"`
}

func (j *JsonDocument) OpenSearch() *OpenSearchDocument {
	newDoc := &OpenSearchDocument{
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
