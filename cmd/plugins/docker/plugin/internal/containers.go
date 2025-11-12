package internal

import (
	"github.com/bgrewell/dtac-agent/cmd/plugins/docker/plugin/internal/utilities"
)

type listContainerOptions struct {
	all     bool
	size    bool
	limit   int
	since   string
	before  string
	filters map[string][]string
}

// ListContainerOptions is a typedef that defines the return value for container option functions
type ListContainerOptions func(containerOptions *listContainerOptions)

// WithContainerShowAll is an options function to set the all value
func WithContainerShowAll(all bool) ListContainerOptions {
	return func(o *listContainerOptions) {
		o.all = all
	}
}

// WithContainerShowSize is an options function to set the size value
func WithContainerShowSize(size bool) ListContainerOptions {
	return func(o *listContainerOptions) {
		o.size = size
	}
}

// WithContainerLimit is an options function to set the limit value
func WithContainerLimit(limit int) ListContainerOptions {
	return func(o *listContainerOptions) {
		o.limit = limit
	}
}

// WithContainerSince is an options function to set the since value
func WithContainerSince(since string) ListContainerOptions {
	return func(o *listContainerOptions) {
		o.since = since
	}
}

// WithContainerBefore is an options function to set the before value
func WithContainerBefore(before string) ListContainerOptions {
	return func(o *listContainerOptions) {
		o.before = before
	}
}

// WithContainerFilters is an options function to set the filters options
func WithContainerFilters(filters map[string][]string) ListContainerOptions {
	return func(o *listContainerOptions) {
		o.filters = filters
	}
}

// ParseListContainerOptions is a helper function to parse parameters into a list of container options
func ParseListContainerOptions(parameters map[string][]string) (options []ListContainerOptions, err error) {
	// No parameters are required so we can ignore the errors
	all, err := utilities.ExtractBool(parameters, "all")
	if err == nil {
		options = append(options, WithContainerShowAll(all))
	}

	digests, err := utilities.ExtractBool(parameters, "size")
	if err == nil {
		options = append(options, WithContainerShowSize(digests))
	}

	limit, err := utilities.ExtractInt(parameters, "limit")
	if err == nil {
		options = append(options, WithContainerLimit(limit))
	}

	since, err := utilities.ExtractString(parameters, "since")
	if err == nil {
		options = append(options, WithContainerSince(since))
	}

	before, err := utilities.ExtractString(parameters, "before")
	if err == nil {
		options = append(options, WithContainerBefore(before))
	}

	return options, nil
}
