
```bash

docker-compose up -d --build rabbitmq
docker-compose logs -f --tail 100 rabbitmq

curl -v -H "Accept:text/plain" "http://localhost:15692/metrics"

docker-compose up -d --build prometheus
docker-compose logs -f --tail 100 prometheus


cd src
go mod download

go run cmd/producer/main.go
go run cmd/consumer/main.go

```
