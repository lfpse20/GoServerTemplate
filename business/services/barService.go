package services

import "GoServerTemplate/business/entities"

type BarServicer interface {
	AddBar(user entities.Bar) (entities.Bar, error)
}

type BarService struct{}

func NewBarService() BarServicer {
	return BarService{}
}

func (barService BarService) AddBar(bar entities.Bar) (entities.Bar, error) {
	return bar, nil
}
