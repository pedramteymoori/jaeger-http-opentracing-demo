package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

func main() {
	_, closer := configureTracer()
	defer closer.Close()

	http.HandleFunc("/", serveHTTP)
	panic(http.ListenAndServe(":8090", nil))
}

func configureTracer() (opentracing.Tracer, io.Closer) {
	cfg := jaegercfg.Configuration{
		ServiceName: "Cafebazaar-service1",
		Sampler: &jaegercfg.SamplerConfig{
			SamplingServerURL: "http://127.0.0.1:5778/sampling",
			Type:              jaeger.SamplerTypeConst,
			Param:             1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		fmt.Println(err)
	}

	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}
