package business

import "GoServerTemplate/business/services"

func InitServices() services.ServiceList {
	serviceList := services.ServiceList{}

	serviceList.FooService = services.NewFooService()
	serviceList.BarService = services.NewBarService()

	return serviceList
}
