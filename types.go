// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package bconnector

type bEvent struct {
	Data    string `json:"event.data"`
	TraceID string `json:"trace_id"`
	SpanID  string `json:"span_id"`
	//EntityMetadata entityMetadata `json:"entityMetadata"`
	Name         string `json:"event.name"`
	Domain       string `json:"event.domain"`
	Bt           string `json:"bt.id"`
	EventType    string `json:"event.type"`
	IsEvent      bool   `json:"appd.isevent"`
	AppEventType string `json:"appd.event.type"`
	EntityId     string `json:"entity_id"`
}

type entityMetadata struct {
	Type entityMetadataType `json:"type"`
}

type entityMetadataType struct {
	Name      string    `json:"name"`
	Namespace namespace `json:"namespace"`
}

type namespace struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
}
