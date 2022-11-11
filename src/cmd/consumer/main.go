package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
	"runtime"
	"sync"
	"flag"

	"github.com/timandy/routine"
	"github.com/juliocesarscheidt/golang-rabbitmq-worker/internal/dto"
	"github.com/juliocesarscheidt/golang-rabbitmq-worker/internal/order/infra/database"
	"github.com/juliocesarscheidt/golang-rabbitmq-worker/internal/order/usecase"
	"github.com/juliocesarscheidt/golang-rabbitmq-worker/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

var qtyWorkersToConsume = flag.Int("qtyWorkersToConsume", runtime.NumCPU(), "Quantity of workers to consume")

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error recovered", r)
		}
	}()

	flag.Parse()
	fmt.Println("qtyWorkersToConsume", *qtyWorkersToConsume)

	var wg sync.WaitGroup

	// db, err := sql.Open("sqlite3", ":memory:")
	db, err := sql.Open("sqlite3", "./orders.db")
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()
	repository := database.NewOrderRepository(db)
	uc := usecase.CalculateFinalPriceUseCase{OrderRepository: repository}

	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		fmt.Println(err)
	}

	defer ch.Close()
	out := make(chan amqp.Delivery, 100) // 100 => buffer size
	go rabbitmq.Consume(ch, out)

	for idx := 1; idx <= *qtyWorkersToConsume; idx++ {
		wg.Add(1)
		go worker(out, &uc, idx)
	}

	wg.Wait()
}

func worker(deliveryMsg <-chan amqp.Delivery, uc *usecase.CalculateFinalPriceUseCase, workerID int) {
	// msg := <-deliveryMsg
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error recovered", r)
		}
	}()

	for msg := range deliveryMsg {
		fmt.Printf("Goroutine Number :: %v | Routine GoId :: %v\n", runtime.NumGoroutine(), routine.Goid())

		var inputDTO dto.OrderInputDTO
		err := json.Unmarshal(msg.Body, &inputDTO)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(fmt.Sprintf("Consuming order %s", inputDTO.ID))

		outputDTO, err := uc.Execute(inputDTO)
		if err != nil {
			fmt.Println(err)
		}
		msg.Ack(true)

		fmt.Printf("Worker %d has processed order %s\n", workerID, outputDTO.ID)

		time.Sleep(250 * time.Millisecond)
	}
}
