package internal

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/plugins/docker/plugin/internal/utilities"
)

type listContainerOptions struct {
	all     bool
	size    bool
	limit   int
	since   string
	before  string
	filters map[string][]string
}

type ListContainerOptions func(containerOptions *listContainerOptions)

func ParseListContainerOptions(parameters map[string][]string) (options []ListContainerOptions, err error) {
	// No parameters are required so we can ignore the errors
	all, err := utilities.ExtractBool(parameters, "all")
	if err == nil {
		options = append(options, WithAll(all))
	}

	digests, err := utilities.ExtractBool(parameters, "digests")
	if err == nil {
		options = append(options, WithDigests(digests))
	}

	filter, err := utilities.ExtractString(parameters, "filter")
	if err == nil {
		options = append(options, WithFilter(filter))
	}

	return options, nil
}
