package internal

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/plugins/docker/plugin/internal/utilities"
)

type listImageOptions struct {
	filters map[string][]string
	all     bool
	digests bool
	filter  string
}

type ListImageOptions func(imageOption *listImageOptions)

func ParseListImageOptions(parameters map[string][]string) (options []ListImageOptions, err error) {
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
