package retrieval

import (
	"errors"
	"sort"

	"github.com/postfinance/kubewire/pkg/report"
	"k8s.io/apimachinery/pkg/api/meta"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Resources retrieves all API resources, the result is sorted
func Resources(config *rest.Config) ([]report.Resource, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// Discover resources
	res, err := clientset.Discovery().ServerResources()
	if err != nil {
		return nil, err
	}

	// Convert them
	rep := []report.Resource{}
	for _, v := range res {
		for _, resource := range v.APIResources {
			res := report.Resource{
				GroupVersion: v.GroupVersion,
				Kind:         resource.Kind,
				Namespaced:   resource.Namespaced,
				Name:         resource.Name,
				Listable:     sliceContains(resource.Verbs, "list"),
			}
			rep = append(rep, res)
		}
	}

	sort.Sort(report.ResourceSort(rep))

	return rep, nil
}

// ResourceObjects retrieves all API resource objects which are global or
// in the list of provided namespaces, the result is sorted
func ResourceObjects(config *rest.Config, namespaces []string) ([]report.ResourceObject, error) {
	// Create client
	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// Get resources
	resources, err := Resources(config)
	if err != nil {
		return nil, err
	}

	rep := []report.ResourceObject{}
	for _, resource := range resources {
		if !resource.Listable {
			// do not try to scrape non listable objects
			continue
		}

		// Parse Group/Version
		gv, err2 := schema.ParseGroupVersion(resource.GroupVersion)
		if err2 != nil {
			return nil, err
		}

		// Create Resource Interface
		rif := dyn.Resource(gv.WithResource(resource.Name))

		// If namespaced
		if resource.Namespaced {
			for _, ns := range namespaces {
				li, err := rif.Namespace(ns).List(meta_v1.ListOptions{})
				if err != nil {
					return nil, err
				}

				items, err := ExtractRuntimeObjectList(li, resource)
				if err != nil {
					return nil, err
				}

				rep = append(rep, items...)
			}
		} else {
			li, err := rif.List(meta_v1.ListOptions{})
			if err != nil {
				return nil, err
			}

			items, err := ExtractRuntimeObjectList(li, resource)
			if err != nil {
				return nil, err
			}

			rep = append(rep, items...)
		}

	}

	// Sort results
	sort.Sort(report.ResourceObjectSort(rep))
	return rep, err
}

// ExtractRuntimeObjectList extracts the list from a raw listing result and converts
// it to a ResourceObject slice
func ExtractRuntimeObjectList(li runtime.Object, resource report.Resource) ([]report.ResourceObject, error) {
	rep := []report.ResourceObject{}

	// Parse list
	items, err := meta.ExtractList(li)
	if err != nil {
		return nil, err
	}

	// Iterate through objects
	for _, item := range items {
		unstructured, ok := item.(runtime.Unstructured)
		if !ok {
			return nil, errors.New("assertion to runtime.Unstructured failed")
		}

		metadata, err := meta.Accessor(unstructured)
		if err != nil {
			return nil, err
		}

		// Create a new ResourceObject
		obj := report.ResourceObject{
			GroupVersion: resource.GroupVersion,
			Name:         metadata.GetName(),
			Namespace:    metadata.GetNamespace(),
			Resource:     resource.Name,
		}

		rep = append(rep, obj)
	}

	return rep, nil
}

func sliceContains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}

	return false
}
