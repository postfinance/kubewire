package report

import (
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ResourceSort implements sorting for Resources by Key()
type ResourceSort []Resource

func (r ResourceSort) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r ResourceSort) Less(i, j int) bool {
	return strings.Compare(r[i].Key(), r[j].Key()) == -1
}

func (r ResourceSort) Len() int {
	return len(r)
}

// ResourceObjectSort implements sorting for ResourceObjects by Key()
type ResourceObjectSort []ResourceObject

func (r ResourceObjectSort) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r ResourceObjectSort) Less(i, j int) bool {
	return strings.Compare(r[i].Key(), r[j].Key()) == -1
}

func (r ResourceObjectSort) Len() int {
	return len(r)
}

// SplitGroupVersionSafe is a helper to ensure that the sorting is done by
// group first.
// Fallback in error case:
// group: ""
// version: groupversion
func SplitGroupVersionSafe(groupversion string) (string, string) {
	gv, err := schema.ParseGroupVersion(groupversion)
	if err != nil {
		return "", groupversion
	}

	return gv.Group, gv.Version
}
