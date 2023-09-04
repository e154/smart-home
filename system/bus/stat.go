package bus

import "strings"

// Stat ...
type Stat struct {
	Topic       string
	Subscribers int
}

// Stats ...
type Stats []Stat

func (s Stats) Len() int           { return len(s) }
func (s Stats) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Stats) Less(i, j int) bool { return strings.Compare(s[i].Topic, s[j].Topic) == -1 }
