// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package bconnector // import "github.com/open-telemetry/opentelemetry-collector-contrib/connector/bconnector"

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/consumer"
)

func NewFactory() connector.Factory {
	return connector.NewFactory(
		"bconnector",
		createDefaultConfig,
		connector.WithTracesToLogs(createTracesToLogs, component.StabilityLevelDevelopment),
	)
}

// createTracesToLogs creates a traces to logs connector based on provided config.
func createTracesToLogs(_ context.Context, params connector.CreateSettings,
	cfg component.Config, nextConsumer consumer.Logs) (connector.Traces, error) {
	c := newConnector(params.Logger, cfg, &eventGenerator{})
	c.logConsumer = nextConsumer
	return c, nil
}
