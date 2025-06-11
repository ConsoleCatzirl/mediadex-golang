package item

type Family string

const (
	FeatureFamily Family = "feature"
	MusicFamily          = "music"
	EpisodeFamily        = "episode"
)

func (f *Family) String() string {
	return string(*f)
}
