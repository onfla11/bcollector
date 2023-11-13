// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package bconnector // import "github.com/open-telemetry/opentelemetry-collector-contrib/connector/bconnector"

import (
	"context"

	"go.uber.org/zap"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type bconnectorImp struct {
	logConsumer consumer.Logs
	config      component.Config
	logger      *zap.Logger
	eventGen    *eventGenerator
}

// function to create a new connector
func newConnector(logger *zap.Logger, config component.Config, generator *eventGenerator) *bconnectorImp {
	logger.Info("Building bconnector connector")
	oCfg := config.(*Config)

	return &bconnectorImp{
		config:   oCfg,
		logger:   logger,
		eventGen: generator,
	}
}

// Capabilities implements the consumer interface.
func (c *bconnectorImp) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func (c *bconnectorImp) ConsumeLogs(_ context.Context, log plog.Logs) error {
	c.logger.Info("consumer logs: ", zap.Any("log", log))
	return nil
}

func (c *bconnectorImp) ConsumeTraces(ctx context.Context, trace ptrace.Traces) error {
	event, err := c.eventGen.generateLog(c.logger, trace)
	if err != nil {
		return err
	}
	if event != nil {
		c.logger.Info("ConsumeTraces: Event Generated: ", zap.Any("event", event))
	}

	err = c.logConsumer.ConsumeLogs(ctx, *event)
	if err != nil {
		return err
	}
	return nil
}

// Start implements the component.Component interface.
func (c *bconnectorImp) Start(_ context.Context, _ component.Host) error {
	c.logger.Info("Starting spanmetrics connector")
	return nil
}

// Shutdown implements the component.Component interface.
func (c *bconnectorImp) Shutdown(context.Context) error {
	c.logger.Info("Shutting down spanmetrics connector")
	return nil
}
