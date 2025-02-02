// Copyright The OpenTelemetry Authors
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

package otlpmetrictest // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal/otlpmetrictest"

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/metrictest"
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/metric/sdkapi"
	exportmetric "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/processor/processortest"
)

// OneRecordReader is a Reader that returns just one
// filled record. It may be useful for testing driver's metrics
// export.
func OneRecordReader() exportmetric.InstrumentationLibraryReader {
	desc := metrictest.NewDescriptor(
		"foo",
		sdkapi.CounterInstrumentKind,
		number.Int64Kind,
	)
	agg := sum.New(1)
	if err := agg[0].Update(context.Background(), number.NewInt64Number(42), &desc); err != nil {
		panic(err)
	}
	start := time.Date(2020, time.December, 8, 19, 15, 0, 0, time.UTC)
	end := time.Date(2020, time.December, 8, 19, 16, 0, 0, time.UTC)
	labels := attribute.NewSet(attribute.String("abc", "def"), attribute.Int64("one", 1))
	rec := exportmetric.NewRecord(&desc, &labels, agg[0].Aggregation(), start, end)

	return processortest.MultiInstrumentationLibraryReader(
		map[instrumentation.Library][]exportmetric.Record{
			{
				Name: "onelib",
			}: {rec},
		})
}

func EmptyReader() exportmetric.InstrumentationLibraryReader {
	return processortest.MultiInstrumentationLibraryReader(nil)
}

// FailReader is a checkpointer that returns an error during
// ForEach.
type FailReader struct{}

var _ exportmetric.InstrumentationLibraryReader = FailReader{}

// ForEach implements exportmetric.Reader. It always fails.
func (FailReader) ForEach(readerFunc func(instrumentation.Library, exportmetric.Reader) error) error {
	return fmt.Errorf("fail")
}
