#!/usr/bin/env bash

go build -o cmd/server/server cmd/server/*.go
go build -o cmd/agent/agent cmd/agent/*.go

increment="$1"

if [ $increment = "1" ]; then
    ./metricstest -test.v -test.run=^TestIteration1$ \
                  -binary-path=cmd/server/server
fi

if [ "$increment" = "2" ]; then
    ./metricstest -test.v -test.run=^TestIteration2[AB]$ \
                -source-path=. \
                -agent-binary-path=cmd/agent/agent
fi

if [ "$increment" = "3" ]; then
    ./metricstest -test.v -test.run=^TestIteration3[AB]*$ \
                  -source-path=. \
                  -agent-binary-path=cmd/agent/agent \
                  -binary-path=cmd/server/server
fi

if [ "$increment" = "4" ]; then
    ./metricstest -test.v -test.run=^TestIteration4$ \
                -agent-binary-path=cmd/agent/agent \
                -binary-path=cmd/server/server \
                -server-port=8080 \
                -source-path=.
fi

if [ "$increment" = "5" ]; then
    ./metricstest -test.v -test.run=^TestIteration5$ \
                -agent-binary-path=cmd/agent/agent \
                -binary-path=cmd/server/server \
                -server-port=8888 \
                -source-path=.
fi

if [ "$increment" = "6" ]; then
    ./metricstest -test.v -test.run=^TestIteration6$ \
                -agent-binary-path=cmd/agent/agent \
                -binary-path=cmd/server/server \
                -server-port=8888 \
                -source-path=.
fi

if [ "$increment" = "7" ]; then
    ./metricstest -test.v -test.run=^TestIteration7$ \
                  -agent-binary-path=cmd/agent/agent \
                  -binary-path=cmd/server/server \
                  -server-port=8888 \
                  -source-path=.
fi

if [ "$increment" = "8" ]; then
    ./metricstest -test.v -test.run=^TestIteration8$ \
                  -agent-binary-path=cmd/agent/agent \
                  -binary-path=cmd/server/server \
                  -server-port=8888 \
                  -source-path=.
fi
