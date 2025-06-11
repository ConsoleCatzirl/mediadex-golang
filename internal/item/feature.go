package item

type FeatureItem struct {
	baseItem
}

func FileToFeature(fItem *FileItem) *FeatureItem {
	return &FeatureItem{baseItem{JsonDocument{FileStats: fItem}}}
}

func (f *FeatureItem) AddMetadata() error {
	err := f.addMediaInfo()
	if err != nil {
		return err
	}

	// todo: web metadata

	return nil
}
