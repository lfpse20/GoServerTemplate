package services

import "GoServerTemplate/business/entities"

type FooServicer interface {
	AddFoo(user entities.Foo) (entities.Foo, error)
}

type FooService struct{}

func NewFooService() FooServicer {
	return FooService{}
}

func (fooService FooService) AddFoo(foo entities.Foo) (entities.Foo, error) {
	return foo, nil
}
