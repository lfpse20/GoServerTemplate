package services

import "GoServerTemplate/business/entities"

type BarServicer interface {
	AddBar(user entities.Bar) (entities.Bar, error)
}

type barService struct{}

func NewBarService() BarServicer {
	return barService{}
}

func (barService barService) AddBar(bar entities.Bar) (entities.Bar, error) {
	return bar, nil
}
