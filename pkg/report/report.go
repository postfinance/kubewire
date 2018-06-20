package report

import (
	"fmt"
	"time"
)

// Report defines a cluster snapshot
type Report struct {
	ScanStart       time.Time
	ScanEnd         time.Time
	Server          Server
	Resources       []Resource
	ResourceObjects []ResourceObject
	Configuration   Configuration
}

// Configuration defines the scanning configuration used
type Configuration struct {
	Namespaces      []string
	KubewireVersion string
}

// Server holds the information about the remote kubernetes instance
type Server struct {
	Version string
	Host    string
}

// Resource defines a kubernetes api resource
type Resource struct {
	GroupVersion string
	Name         string
	Kind         string
	Namespaced   bool
	Verbs        []string
}

// Key returns a unique identifier for the Resource which can be
// used for sorting
func (a Resource) Key() string {
	// The ASCII code of space is lower than all allowed chars in kubernetes
	// resource names. So here we can use it as separator for a
	// rune wise string comparison
	group, version := SplitGroupVersionSafe(a.GroupVersion)

	return fmt.Sprintf("%s %s %s", group, version, a.Name)

}

// Compare compares two Resources, while b must be a Resource or else it will panic
func (a Resource) Compare(b interface{}) []DiffReport {
	bres, ok := b.(Resource)
	if !ok {
		// This is a programming error so it should never occure at runtime
		panic("Resource can only compare to Resource")
	}

	ret := []DiffReport{}
	if a.GroupVersion != bres.GroupVersion {
		ret = append(ret, DiffReport{Element: "GroupVersion", A: a.GroupVersion, B: bres.GroupVersion})
	}

	if a.Name != bres.Name {
		ret = append(ret, DiffReport{Element: "Name", A: a.Name, B: bres.Name})
	}

	if a.Kind != bres.Kind {
		ret = append(ret, DiffReport{Element: "Kind", A: a.Kind, B: bres.Kind})
	}

	if a.Namespaced != bres.Namespaced {
		ans := fmt.Sprintf("%t", a.Namespaced)
		bns := fmt.Sprintf("%t", bres.Namespaced)
		ret = append(ret, DiffReport{Element: "Namespaced", A: ans, B: bns})
	}

	aVerbs := fmt.Sprintf("%v", a.Verbs)
	bVerbs := fmt.Sprintf("%v", bres.Verbs)
	if aVerbs != bVerbs {
		ret = append(ret, DiffReport{Element: "Verbs", A: aVerbs, B: bVerbs})
	}

	if len(ret) == 0 {
		return nil
	}

	return ret
}

// ResourceObject represent a persistent entity in the system
type ResourceObject struct {
	GroupVersion string
	Resource     string // references Resource.Name
	Namespace    string
	Name         string
}

// Key returns a unique identifier for the ResourceObject which can be
// used for sorting
func (a ResourceObject) Key() string {
	// The ASCII code of space is lower than all allowed chars in kubernetes
	// resource/object names. So here we can use it as separator for a
	// rune wise string comparison
	group, version := SplitGroupVersionSafe(a.GroupVersion)
	return fmt.Sprintf("%s %s %s %s %s", group, version, a.Resource, a.Namespace, a.Name)
}

// Compare compares two ResourceObjects, while b must be a ResourceObject or else it will panic
func (a ResourceObject) Compare(b interface{}) []DiffReport {
	bres, ok := b.(ResourceObject)
	if !ok {
		// This is a programming error so it should never occure at runtime
		panic("ResourceObject can only compare to ResourceObject")
	}

	ret := []DiffReport{}
	if a.Name != bres.Name {
		ret = append(ret, DiffReport{Element: "Name", A: a.Name, B: bres.Name})
	}

	if a.Namespace != bres.Namespace {
		ret = append(ret, DiffReport{Element: "Namespace", A: a.Namespace, B: bres.Namespace})
	}

	if a.Resource != bres.Resource {
		ret = append(ret, DiffReport{Element: "Resource", A: a.Resource, B: bres.Resource})
	}

	if a.GroupVersion != bres.GroupVersion {
		ret = append(ret, DiffReport{Element: "GroupVersion", A: a.GroupVersion, B: bres.GroupVersion})
	}

	if len(ret) == 0 {
		return nil
	}

	return ret
}
