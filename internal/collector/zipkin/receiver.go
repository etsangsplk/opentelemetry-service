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

// Package zipkinreceiver wraps the functionality to start the end-point that
// receives Zipkin traces.
package zipkinreceiver

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-service/cmd/occollector/app/builder"
	"github.com/open-telemetry/opentelemetry-service/consumer"
	"github.com/open-telemetry/opentelemetry-service/receiver"
	"github.com/open-telemetry/opentelemetry-service/receiver/zipkinreceiver"
)

// Start starts the Zipkin receiver endpoint.
func Start(logger *zap.Logger, v *viper.Viper, traceConsumer consumer.TraceConsumer, asyncErrorChan chan<- error) (receiver.TraceReceiver, error) {
	rOpts, err := builder.NewDefaultZipkinReceiverCfg().InitFromViper(v)
	if err != nil {
		return nil, err
	}

	addr := ":" + strconv.FormatInt(int64(rOpts.Port), 10)
	zi, err := zipkinreceiver.New(addr, traceConsumer)
	if err != nil {
		return nil, fmt.Errorf("Failed to create the Zipkin receiver: %v", err)
	}

	if err := zi.StartTraceReception(context.Background(), asyncErrorChan); err != nil {
		return nil, fmt.Errorf("Cannot start Zipkin receiver to address %q: %v", addr, err)
	}

	logger.Info("Zipkin receiver is running.", zap.Int("port", rOpts.Port))

	return zi, nil
}
