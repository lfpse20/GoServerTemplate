package services

import "GoServerTemplate/business/entities"

type FooServicer interface {
	AddFoo(user entities.Foo) (entities.Foo, error)
}

type fooService struct{}

func NewFooService() FooServicer {
	return fooService{}
}

func (fooService fooService) AddFoo(foo entities.Foo) (entities.Foo, error) {
	return foo, nil
}
