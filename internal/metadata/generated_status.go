// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metadata // import "github.com/open-telemetry/opentelemetry-collector-contrib/connector/bconnector/internal/metadata"

import (
	"go.opentelemetry.io/collector/component"
)

const (
	Type                  = "bconnector"
	TracesToLogsStability = component.StabilityLevelDevelopment
)
