package item

type MusicItem struct {
	baseItem
}

func FileToMusic(fItem *FileItem) *MusicItem {
	return &MusicItem{baseItem{JsonDocument{FileStats: fItem}}}
}

func (f *MusicItem) AddMetadata() error {
	err := f.addMediaInfo()
	if err != nil {
		return err
	}

	// todo: web metadata

	return nil
}
