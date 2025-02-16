// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

//go:build !test

package api

import (
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/mainflux/mainflux/consumers"
)

var _ consumers.BlockingConsumer = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter  metrics.Counter
	latency  metrics.Histogram
	consumer consumers.BlockingConsumer
}

// MetricsMiddleware returns new message repository
// with Save method wrapped to expose metrics.
func MetricsMiddleware(consumer consumers.BlockingConsumer, counter metrics.Counter, latency metrics.Histogram) consumers.BlockingConsumer {
	return &metricsMiddleware{
		counter:  counter,
		latency:  latency,
		consumer: consumer,
	}
}

func (mm *metricsMiddleware) ConsumeBlocking(msgs interface{}) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "consume").Add(1)
		mm.latency.With("method", "consume").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mm.consumer.ConsumeBlocking(msgs)
}
