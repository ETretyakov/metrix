go build -o cmd/server/server cmd/server/*.go
go build -o cmd/agent/agent cmd/agent/*.go

./metricstest -test.v -test.run=^TestIteration1$ \
              -binary-path=cmd/server/server
./metricstest -test.v -test.run=^TestIteration2[AB]$ \
              -source-path=. \
              -agent-binary-path=cmd/agent/agent
./metricstest -test.v -test.run=^TestIteration3[AB]*$ \
              -source-path=. \
              -agent-binary-path=cmd/agent/agent \
              -binary-path=cmd/server/server
./metricstest -test.v -test.run=^TestIteration4$ \
              -agent-binary-path=cmd/agent/agent \
              -binary-path=cmd/server/server \
              -server-port=8080 \
              -source-path=.
./metricstest -test.v -test.run=^TestIteration5$ \
              -agent-binary-path=cmd/agent/agent \
              -binary-path=cmd/server/server \
              -server-port=8888 \
              -source-path=.
./metricstest -test.v -test.run=^TestIteration6$ \
              -agent-binary-path=cmd/agent/agent \
              -binary-path=cmd/server/server \
              -server-port=8888 \
              -source-path=.