package usecase

import (
	"github.com/juliocesarscheidt/golang-rabbitmq-worker/internal/order/entity"
	"github.com/juliocesarscheidt/golang-rabbitmq-worker/internal/order/infra/database"
	"github.com/juliocesarscheidt/golang-rabbitmq-worker/internal/dto"

	// sqlite3
	_ "github.com/mattn/go-sqlite3"
)

type CalculateFinalPriceUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewCalculateFinalPriceUseCase(orderRepository database.OrderRepository) *CalculateFinalPriceUseCase {
	return &CalculateFinalPriceUseCase{
		OrderRepository: &orderRepository,
	}
}

func (c *CalculateFinalPriceUseCase) Execute(input dto.OrderInputDTO) (*dto.OrderOutputDTO, error) {
	order, err := entity.NewOrder(input.ID, input.Price, input.Tax)
	if err != nil {
		return nil, err
	}
	err = order.CalculateFinalPrice()
	if err != nil {
		return nil, err
	}
	err = c.OrderRepository.Save(order)
	if err != nil {
		return nil, err
	}
	return &dto.OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}, nil
}
