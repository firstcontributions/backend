package reputation

import (
	"strings"

	"github.com/firstcontributions/backend/internal/models/usersstore"
)

var fileExtensionBadgeMap = map[string]string{
	"go":   "Go",
	"py":   "Python",
	"c":    "C",
	"cpp":  "C++",
	"js":   "JavaScript",
	"ts":   "TypeScript",
	"java": "Java",
	"rb":   "Ruby",
	"rs":   "Rust",
	"sh":   "Unix Sell",
	"php":  "PHP",
	"html":  "HTML",
	"css":  "CSS",
	"md":  "Markdown",
}

type BadgeMap struct {
	data map[string]*usersstore.Badge
}

func (b *BadgeMap) Add(path string, additions int) {
	// TODO: ignore vendor files like node_modules
	filesplit := strings.Split(path, ".")
	if len(filesplit) < 2 {
		return
	}
	name, ok := fileExtensionBadgeMap[filesplit[1]]
	if !ok {
		return
	}
	if b.data[name] == nil {
		b.data[name] = &usersstore.Badge{
			DisplayName: name,
		}
	}
	b.data[name].Points += int64(additions)
}

func (b *BadgeMap) ToBadges() []*usersstore.Badge {
	badges := []*usersstore.Badge{}

	for _, badge := range b.data {
		badges = append(badges, badge)
	}
	return badges
}

func BadgeMapFromBadges(badges []*usersstore.Badge) *BadgeMap {
	b := &BadgeMap{
		data: map[string]*usersstore.Badge{},
	}
	for _, bd := range badges {
		b.data[bd.DisplayName] = bd
	}
	return b
}
