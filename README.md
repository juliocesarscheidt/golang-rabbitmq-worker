# Golang Rabbitmq Worker

```bash

docker-compose up -d --build rabbitmq
docker-compose logs -f --tail 100 rabbitmq
# rabbitmq metrics
curl -H "Accept:text/plain" --url "http://localhost:15692/metrics"

docker-compose up -d prometheus
docker-compose logs -f --tail 100 prometheus
# prometheus URL: http://localhost:9090
curl -H "Accept:text/plain" --url "http://localhost:9090/api/v1/query?query=rabbitmq_queue_messages&time=$(date +%s)"
# {"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"rabbitmq_queue_messages","instance":"rabbitmq:15692","job":"rabbitmq"},"value":[1668613166,"10"]}]}}

# running producer and consumer
cd src
go mod download

# producing
go run cmd/producer/main.go -qtyOrdersToProduce 10

# consuming
go run cmd/consumer/main.go -qtyWorkersToConsume 4
go run cmd/consumer/main.go -qtyWorkersToConsume $(nproc)

```
