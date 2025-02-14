// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package stats

import (
	"testing"

	pb "github.com/DataDog/datadog-agent/pkg/proto/pbgo/trace"
	"github.com/stretchr/testify/assert"
)

func TestGetStatusCode(t *testing.T) {
	for _, tt := range []struct {
		in  *pb.Span
		out uint32
	}{
		{
			&pb.Span{},
			0,
		},
		{
			&pb.Span{
				Meta: map[string]string{"http.status_code": "200"},
			},
			200,
		},
		{
			&pb.Span{
				Metrics: map[string]float64{"http.status_code": 302},
			},
			302,
		},
		{
			&pb.Span{
				Meta:    map[string]string{"http.status_code": "200"},
				Metrics: map[string]float64{"http.status_code": 302},
			},
			302,
		},
		{
			&pb.Span{
				Meta: map[string]string{"http.status_code": "x"},
			},
			0,
		},
	} {
		if got := getStatusCode(tt.in); got != tt.out {
			t.Fatalf("Expected %d, got %d", tt.out, got)
		}
	}
}

func TestNewAggregation(t *testing.T) {
	for _, tt := range []struct {
		in               *pb.Span
		enablePeerSvcAgg bool
		res              Aggregation
	}{
		{
			&pb.Span{},
			false,
			Aggregation{},
		},
		{
			&pb.Span{},
			true,
			Aggregation{},
		},
		{
			&pb.Span{
				Service: "a",
				Meta:    map[string]string{"peer.service": "remote-service"},
			},
			false,
			Aggregation{BucketsAggregationKey: BucketsAggregationKey{Service: "a"}},
		},
		{
			&pb.Span{
				Service: "a",
				Meta:    map[string]string{"peer.service": "remote-service"},
			},
			true,
			Aggregation{BucketsAggregationKey: BucketsAggregationKey{Service: "a", PeerService: "remote-service"}},
		},
		{
			&pb.Span{
				Service:  "service",
				Name:     "operation",
				Resource: "resource",
				Meta: map[string]string{
					"span.kind":        "client",
					"peer.service":     "remote-service",
					"http.status_code": "200",
				},
			},
			true,
			Aggregation{
				BucketsAggregationKey: BucketsAggregationKey{
					Service:     "service",
					Name:        "operation",
					PeerService: "remote-service",
					Resource:    "resource",
					SpanKind:    "client",
					StatusCode:  200,
					Synthetics:  false,
				},
			},
		},
	} {
		assert.Equal(t, tt.res, NewAggregationFromSpan(tt.in, "", PayloadAggregationKey{}, tt.enablePeerSvcAgg))
	}
}
