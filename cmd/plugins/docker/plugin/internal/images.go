package internal

import (
	"github.com/bgrewell/dtac-agent/cmd/plugins/docker/plugin/internal/utilities"
)

// ImageInfo is a structure to abstract the docker image info in a more human friendly way
type ImageInfo struct {
	ID          string            `json:"id"`
	RepoTags    []string          `json:"repo_tags"`
	Created     string            `json:"created"`
	Size        string            `json:"size"`
	VirtualSize string            `json:"virtual_size"`
	ParentID    string            `json:"parent_id"`
	RepoDigests []string          `json:"repo_digests"`
	Labels      map[string]string `json:"labels"`
}

type listImageOptions struct {
	filters map[string][]string
	all     bool
	digests bool
	filter  string
}

// ListImageOptions is a typedef for the list image options
type ListImageOptions func(imageOption *listImageOptions)

// WithImageShowAll is an options function to set the all value
func WithImageShowAll(all bool) ListImageOptions {
	return func(o *listImageOptions) {
		o.all = all
	}
}

// WithImageShowDigests is an options function to set the show digests value
func WithImageShowDigests(digests bool) ListImageOptions {
	return func(o *listImageOptions) {
		o.digests = digests
	}
}

// WithImageFilter is an options function to set the filter value
func WithImageFilter(filter string) ListImageOptions {
	return func(o *listImageOptions) {
		o.filter = filter
	}
}

// WithImageFilters is an options function to set the filters value
func WithImageFilters(filters map[string][]string) ListImageOptions {
	return func(o *listImageOptions) {
		o.filters = filters
	}
}

// ParseListImageOptions is a helper function to parse parameters into a set of list image options
func ParseListImageOptions(parameters map[string][]string) (options []ListImageOptions, err error) {
	// No parameters are required so we can ignore the errors
	all, err := utilities.ExtractBool(parameters, "all")
	if err == nil {
		options = append(options, WithImageShowAll(all))
	}

	digests, err := utilities.ExtractBool(parameters, "digests")
	if err == nil {
		options = append(options, WithImageShowDigests(digests))
	}

	filter, err := utilities.ExtractString(parameters, "filter")
	if err == nil {
		options = append(options, WithImageFilter(filter))
	}

	return options, nil
}
