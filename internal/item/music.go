package item

type MusicItem struct {
	*baseItem
}

func FileToMusic(fItem *FileItem) *MusicItem {
	return &MusicItem{&baseItem{&JsonDocument{FileStats: fItem}}}
}

func JsonToMusic(jDoc *JsonDocument) *MusicItem {
	return &MusicItem{&baseItem{jDoc}}
}

func (f *MusicItem) AddMetadata() error {
	err := f.addMediaInfo()
	if err != nil {
		return err
	}

	// todo: web metadata

	return nil
}
