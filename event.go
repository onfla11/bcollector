// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package bconnector // import "github.com/open-telemetry/opentelemetry-collector-contrib/connector/bconnector"

import (
	"errors"
	"fmt"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type eventGenerator struct{}

// consume takes a single trace and generate biq event
func (e *eventGenerator) generate(trace ptrace.Traces) (*plog.LogRecord, error) {
	traceID, err := getTraceID(trace)
	if err != nil {
		return nil, err
	}
	event := plog.NewLogRecord()
	body := event.Body()
	resourceSpans := trace.ResourceSpans()
	for i := 0; i < resourceSpans.Len(); i++ {
		rs := resourceSpans.At(i)
		attr := rs.Resource().Attributes()
		fmt.Println(attr)
	}
	attrMap := body.Map()
	attrMap.PutStr("test1", "val1")
	attrMap.PutStr("test2", "val2")
	attrMap.PutStr("test3", "val3")
	attrMap.PutStr("traceId", traceID.String())

	return &event, nil
}

func getTraceID(td ptrace.Traces) (pcommon.TraceID, error) {
	rss := td.ResourceSpans()
	if rss.Len() == 0 {
		return pcommon.NewTraceIDEmpty(), errors.New("no resource spans are present")
	}

	ilss := rss.At(0).ScopeSpans()
	if ilss.Len() == 0 {
		return pcommon.NewTraceIDEmpty(), errors.New("no scope spans are present")
	}

	spans := ilss.At(0).Spans()
	if spans.Len() == 0 {
		return pcommon.NewTraceIDEmpty(), errors.New("no trace id is present")
	}

	return spans.At(0).TraceID(), nil
}
