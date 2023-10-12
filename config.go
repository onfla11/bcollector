// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package bconnector

import "go.opentelemetry.io/collector/component"

func createDefaultConfig() component.Config {
	return &Config{}
}

type Config struct {
}

var _ component.Config = (*Config)(nil)

// Validate checks if the processor configuration is valid
func (cfg *Config) Validate() error {
	return nil
}
