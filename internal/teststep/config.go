package teststep

import (
	"fmt"
	"strings"
)

type Config interface {
	HasConfiguration() bool
}

type configuration struct {
	raw       string
	directory string
	file      string
}

func Configuration(raw, directory string) (configuration, error) {
	var populatedConfig []string

	if raw != "" {
		populatedConfig = append(populatedConfig, fmt.Sprintf("%q", "raw"))
	}

	if directory != "" {
		populatedConfig = append(populatedConfig, fmt.Sprintf("%q", "directory"))
	}

	if len(populatedConfig) > 1 {
		return configuration{}, fmt.Errorf(
			"both %s are populated, only one configuration option is allowed",
			strings.Join(populatedConfig, " and "),
		)
	}

	return configuration{
		raw:       raw,
		directory: directory,
	}, nil
}

func (c configuration) HasConfiguration() bool {
	if c.raw != "" {
		return true
	}

	// TODO: do we need to read each file in directory and check that it's not empty.
	// TODO: do we also need to verify that it's valid TF config
	if c.directory != "" {
		return true
	}

	return false
}
