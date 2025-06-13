package item

type Family string

const (
	EpisodeFamily Family = "episode"
	FeatureFamily        = "feature"
	MusicFamily          = "music"
)

func (f *Family) String() string {
	return string(*f)
}
