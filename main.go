package main

import (
    "os"

    "github.com/go-kit/kit/endpoint"

    "github.com/go-kit/kit/log"

    "net/http"

    kitprometheus "github.com/go-kit/kit/metrics/prometheus"
    httptransport "github.com/go-kit/kit/transport/http"
    stdprometheus "github.com/prometheus/client_golang/prometheus"
    promhttp "github.com/roman-vynar/client_golang/prometheus/promhttp"
)

func main() {
    logger := log.NewLogfmtLogger(os.Stderr)

    fieldKeys := []string{"method", "error"}
    requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
        Namespace: "my_group",
        Subsystem: "string_service",
        Name:      "request_count",
        Help:      "Number of requests received.",
    }, fieldKeys)
    requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
        Namespace: "my_group",
        Subsystem: "string_service",
        Name:      "request_latency_microseconds",
        Help:      "Total duration of requests in microseconds.",
    }, fieldKeys)
    countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
        Namespace: "my_group",
        Subsystem: "string_service",
        Name:      "count_result",
        Help:      "The result of each count method.",
    }, []string{}) // no fields here

    var svc StringService
    svc = stringService{}
    svc = loggingMiddleware{logger, svc}
    svc = instrumentingMiddleware{requestCount, requestLatency, countResult, svc}

    var uppercase endpoint.Endpoint
    uppercase = makeUppercaseEndpoint(svc)

    var count endpoint.Endpoint
    count = makeCountEndpoint(svc)

    uppercaseHandler := httptransport.NewServer(
        uppercase,
        decodeUppercaseRequest,
        encodeResponse,
    )

    countHandler := httptransport.NewServer(
        count,
        decodeCountRequest,
        encodeResponse,
    )

    http.Handle("/uppercase", uppercaseHandler)
    http.Handle("/count", countHandler)
    http.Handle("/metrics", promhttp.Handler())
    logger.Log("msg", "HTTP", "addr", ":8080")
    logger.Log("err", http.ListenAndServe(":8080", nil))
}
