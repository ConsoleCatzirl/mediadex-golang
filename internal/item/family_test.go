package item

import (
	"fmt"
	"testing"
)

func TestFamilyString(t *testing.T) {
	var testTable = []struct {
		given Family
		want  string
	}{
		{FeatureFamily, "feature"},
		{EpisodeFamily, "episode"},
		{MusicFamily, "music"},
	}

	for _, test := range testTable {
		have := fmt.Sprintf("%s", test.given)
		if have != test.want {
			t.Errorf("want: '%s' ; have: '%s'", test.want, have)
		}
	}
}
