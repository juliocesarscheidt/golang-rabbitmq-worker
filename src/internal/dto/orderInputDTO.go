package dto

type OrderInputDTO struct {
	ID         string	`json:"id"`
	Price      float64	`json:"price"`
	Tax        float64	`json:"tax"`
}
