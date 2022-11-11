package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
	"flag"

	"github.com/google/uuid"
	"github.com/juliocesarscheidt/golang-rabbitmq-worker/internal/order/entity"
	amqp "github.com/rabbitmq/amqp091-go"
)

var qtyOrdersToProduce = flag.Int("qtyOrdersToProduce", 10000000, "Quantity of orders to produce")

func Publish(ch *amqp.Channel, order entity.Order) error {
	body, err := json.Marshal(order)
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("Publishing order %v", string(body)))
	err = ch.Publish(
		"",
		"orders",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func GenerateOrder() entity.Order {
	return entity.Order{
		ID:    uuid.New().String(),
		Price: rand.Float64() * 100,
		Tax:   rand.Float64() * 10,
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error recovered", r)
		}
	}()

	flag.Parse()
	fmt.Println("qtyOrdersToProduce", *qtyOrdersToProduce)

	conn, err := amqp.Dial("amqp://rabbitmq:rabbitmq@127.0.0.1:5672/")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	defer ch.Close()

	for i := 0; i < *qtyOrdersToProduce; i++ {
		o := GenerateOrder()
		Publish(ch, o)
		time.Sleep(250 * time.Millisecond)
	}
}
