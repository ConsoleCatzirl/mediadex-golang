package item

type EpisodeItem struct {
	*baseItem
}

func FileToEpisode(fItem *FileItem) *EpisodeItem {
	return &EpisodeItem{&baseItem{&JsonDocument{FileStats: fItem}}}
}

func JsonToEpisode(jDoc *JsonDocument) *EpisodeItem {
	return &EpisodeItem{&baseItem{jDoc}}
}

func (f *EpisodeItem) AddMetadata() error {
	err := f.addMediaInfo()
	if err != nil {
		return err
	}

	// todo: web metadata

	return nil
}
