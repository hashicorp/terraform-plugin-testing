package teststep

import (
	"fmt"
	"strings"
)

type Config interface {
	GetRaw() string
	HasConfiguration() bool
	HasRaw() bool
}

type configuration struct {
	directory string
	file      string
	raw       string
}

type ConfigurationRequest struct {
	Directory string
	Raw       string
}

func Configuration(configRequest ConfigurationRequest) (configuration, error) {
	var populatedConfig []string
	var config configuration

	if configRequest.Directory != "" {
		populatedConfig = append(populatedConfig, fmt.Sprintf("%q", "directory"))

		config = configuration{
			directory: configRequest.Directory,
		}
	}

	if configRequest.Raw != "" {
		populatedConfig = append(populatedConfig, fmt.Sprintf("%q", "raw"))

		config = configuration{
			raw: configRequest.Raw,
		}
	}

	if len(populatedConfig) > 1 {
		return configuration{}, fmt.Errorf(
			"both %s are populated, only one configuration option is allowed",
			strings.Join(populatedConfig, " and "),
		)
	}

	return config, nil
}

func (c configuration) GetRaw() string {
	return c.raw
}

func (c configuration) HasConfiguration() bool {
	if c.directory != "" {
		return true
	}

	if c.raw != "" {
		return true
	}

	return false
}

func (c configuration) HasRaw() bool {
	return c.raw != ""
}
