package report

import "fmt"

// DiffReport defines the type for a single diffing result where A is the
// old and B is the new value of Element
type DiffReport struct {
	Element string
	A       string
	B       string
}

func (r DiffReport) String() string {
	return fmt.Sprintf("Element: %s, A: %s, B: %s", r.Element, r.A, r.B)
}

// Keyer defines an interface which can uniquely identify and compare types
// together in order to diff them
type Keyer interface {
	Key() string
	Compare(b interface{}) []DiffReport
}

// AnnotateDiffReports annotates every DiffReport with a prefix
func AnnotateDiffReports(r []DiffReport, prefix string) {
	for i := range r {
		r[i].Element = prefix + r[i].Element
	}
}

// DiffReports creates a DiffReport between a and b.
// Resources and ResourceObjects from a and b must be sorted
func DiffReports(a Report, b Report) []DiffReport {
	ret := []DiffReport{}

	if a.ScanStart != b.ScanStart {
		ret = append(ret, DiffReport{Element: "ScanStart", A: a.ScanStart.String(), B: b.ScanStart.String()})
	}

	if a.ScanEnd != b.ScanEnd {
		ret = append(ret, DiffReport{Element: "ScanEnd", A: a.ScanEnd.String(), B: b.ScanEnd.String()})
	}

	// Configuration
	if a.Configuration.KubewireVersion != b.Configuration.KubewireVersion {
		ret = append(ret, DiffReport{Element: "Configuration.KubewireVersion", A: a.Configuration.KubewireVersion, B: b.Configuration.KubewireVersion})
	}

	// Configuration.Namespaces
	ans := fmt.Sprintf("%v", a.Configuration.Namespaces)
	bns := fmt.Sprintf("%v", b.Configuration.Namespaces)
	if ans != bns {
		ret = append(ret, DiffReport{Element: "Configuration.Namespaces", A: ans, B: bns})
	}

	// Server
	if a.Server.Host != b.Server.Host {
		ret = append(ret, DiffReport{Element: "Server.Host", A: a.Server.Host, B: b.Server.Host})
	}
	if a.Server.Version != b.Server.Version {
		ret = append(ret, DiffReport{Element: "Server.Version", A: a.Server.Version, B: b.Server.Version})
	}

	// Resources
	cmp := Diff(resourcesToKeyer(a.Resources), resourcesToKeyer(b.Resources))
	if len(cmp) != 0 {
		AnnotateDiffReports(cmp, "Resources.")
		ret = append(ret, cmp...)
	}

	// ResourceObject
	cmp = Diff(resourceobjectsToKeyer(a.ResourceObjects), resourceobjectsToKeyer(b.ResourceObjects))
	if len(cmp) != 0 {
		AnnotateDiffReports(cmp, "ResourceObjects.")
		ret = append(ret, cmp...)
	}

	return ret
}

func resourcesToKeyer(a []Resource) []Keyer {
	t := make([]Keyer, len(a))

	for i, v := range a {
		t[i] = v
	}

	return t
}

func resourceobjectsToKeyer(a []ResourceObject) []Keyer {
	t := make([]Keyer, len(a))

	for i, v := range a {
		t[i] = v
	}

	return t
}

// rangeDiff returns a DiffReport for not existing elements over a range
func rangeDiff(x []Keyer, start int, end int, a bool) []DiffReport {
	ret := []DiffReport{}
	for _, v := range x[start:end] {
		if a {
			ret = append(ret, DiffReport{Element: "\"" + v.Key() + "\"", A: "exists", B: "does not exist"})
		} else {
			ret = append(ret, DiffReport{Element: "\"" + v.Key() + "\"", B: "exists", A: "does not exist"})
		}
	}

	return ret
}

// Diff creates a difference report betweed a and b. It assumes that the
// elements in a and b are sorted by Key()
func Diff(a, b []Keyer) []DiffReport {
	rep := []DiffReport{}

	var aIndex int
	var bIndex int

	// Iterate through a and b, stop when some of them has their bounds exceeded
outer:
	for aIndex < len(a) && bIndex < len(b) {
		// if a[aIndex] exists in b[>=bIndex]
		for bIndexNew := bIndex; bIndexNew < len(b); bIndexNew++ {
			if a[aIndex].Key() == b[bIndexNew].Key() {
				// Diff skiped elements
				if bIndex != bIndexNew {
					cmp := rangeDiff(b, bIndex, bIndexNew, false)
					rep = append(rep, cmp...)
				}

				// Append difference of both if there is one
				if cmp := a[aIndex].Compare(b[bIndexNew]); cmp != nil {
					AnnotateDiffReports(cmp, a[aIndex].Key()+".")
					rep = append(rep, cmp...)
				}

				aIndex++
				bIndex = bIndexNew + 1
				continue outer
			}

		}

		// else look if b[bIndex] exists in a[>=aIndex]
		for aIndexNew := aIndex; aIndexNew < len(a); aIndexNew++ {
			if a[aIndexNew].Key() == b[bIndex].Key() {
				// Diff skiped elements
				if aIndex != aIndexNew {
					cmp := rangeDiff(a, aIndex, aIndexNew, true)
					rep = append(rep, cmp...)
				}

				// Append difference of both if there is one
				if cmp := a[aIndexNew].Compare(b[bIndex]); cmp != nil {
					AnnotateDiffReports(cmp, a[aIndexNew].Key()+".")
					rep = append(rep, cmp...)
				}

				aIndex = aIndexNew + 1
				bIndex++
				continue outer
			}
		}

		// Gone too far
		break
	}

	// Process the cutoff of a and b
	if aIndex < len(a) {
		cmp := rangeDiff(a, aIndex, len(a), true)
		rep = append(rep, cmp...)
	}

	if bIndex < len(b) {
		cmp := rangeDiff(b, bIndex, len(b), false)
		rep = append(rep, cmp...)
	}

	return rep
}
