package mocks

import (
	"fmt"
	"github.com/ns1labs/orb/sinks"
	"sort"
)

func sortSinks(pm sinks.PageMetadata, sks []sinks.Sink) []sinks.Sink {
	switch pm.Order {
	case "name":
		if pm.Dir == "asc" {
			sort.SliceStable(sks, func(i, j int) bool {
				return sks[i].Name.String() < sks[j].Name.String()
			})
		}
		if pm.Dir == "desc" {
			sort.SliceStable(sks, func(i, j int) bool {
				return sks[i].Name.String() > sks[j].Name.String()
			})
		}
	case "id":
		if pm.Dir == "asc" {
			sort.SliceStable(sks, func(i, j int) bool {
				return sks[i].ID < sks[j].ID
			})
		}
		if pm.Dir == "desc" {
			sort.SliceStable(sks, func(i, j int) bool {
				return sks[i].ID > sks[j].ID
			})
		}
	default:
		sort.SliceStable(sks, func(i, j int) bool {
			return sks[i].ID < sks[j].ID
		})
	}

	return sks
}

// Since mocks will store data in map, and they need to resemble the real
// identifiers as much as possible, a key will be created as combination of
// owner and their own identifiers. This will allow searching either by
// prefix or suffix.
func key(owner string, id string) string {
	return fmt.Sprintf("%s-%s", owner, id)
}
