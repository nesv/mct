package microcode

import (
	"slices"
	"strings"
)

type Action struct {
	Name string
	Args []string
}

func (a Action) String() string {
	if a.IsZero() {
		return ""
	}
	return strings.Join([]string{
		a.Name,
		strings.Join(a.Args, " "),
	}, " ")
}

func (a Action) IsZero() bool {
	return a.Name == "" && len(a.Args) == 0
}

func (a Action) Equals(other Action) bool {
	return a.Name == other.Name && slices.Compare(a.Args, other.Args) == 0
}
