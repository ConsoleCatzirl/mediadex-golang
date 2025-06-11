package item

type Family string

const (
	MovieFamily  Family = "movies"
	MusicFamily         = "music"
	SeriesFamily        = "series"
)

func (f *Family) String() string {
	return string(*f)
}
