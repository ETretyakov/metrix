go build -o cmd/server/server cmd/server/*.go
go build -o cmd/agent/agent cmd/agent/*.go

./metricstest-darwin-arm64 -test.v -test.run=^TestIteration1$ -binary-path=cmd/server/server