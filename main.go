package main

import (
    "os"

    "github.com/go-kit/kit/endpoint"
    "github.com/go-kit/kit/log"

    "net/http"

    httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
    logger := log.NewLogfmtLogger(os.Stderr)

    var svc StringService
    svc = stringService{}
    svc = loggingMiddleware{logger, svc}

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
}
