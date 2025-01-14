// Copyright 2019, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package jaegerreceiver wraps the functionality to start the end-point that
// receives Jaeger data sent by the jaeger-agent in jaeger.thrift format over
// TChannel and directly from clients in jaeger.thrift format over binary thrift
// protocol (HTTP transport).
// Note that the UDP transport is not supported since these protocol/transport
// are for task->jaeger-agent communication only and the receiver does not try to
// support jaeger-agent endpoints.
// TODO: add support for the jaeger proto endpoint released in jaeger 1.8package jaegerreceiver
package jaegerreceiver

import (
	"context"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-service/cmd/occollector/app/builder"
	"github.com/open-telemetry/opentelemetry-service/consumer"
	"github.com/open-telemetry/opentelemetry-service/receiver"
	"github.com/open-telemetry/opentelemetry-service/receiver/jaegerreceiver"
)

// Start starts the Jaeger receiver endpoint.
func Start(logger *zap.Logger, v *viper.Viper, traceConsumer consumer.TraceConsumer, asyncErrorChan chan<- error) (receiver.TraceReceiver, error) {
	rOpts, err := builder.NewDefaultJaegerReceiverCfg().InitFromViper(v)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	config := &jaegerreceiver.Configuration{
		CollectorThriftPort: rOpts.ThriftTChannelPort,
		CollectorHTTPPort:   rOpts.ThriftHTTPPort,
	}
	jtr, err := jaegerreceiver.New(ctx, config, traceConsumer)
	if err != nil {
		return nil, err
	}

	if err := jtr.StartTraceReception(ctx, asyncErrorChan); err != nil {
		return nil, err
	}

	logger.Info("Jaeger receiver is running.",
		zap.Int("thrift-tchannel-port", rOpts.ThriftTChannelPort),
		zap.Int("thrift-http-port", rOpts.ThriftHTTPPort))

	return jtr, nil
}
