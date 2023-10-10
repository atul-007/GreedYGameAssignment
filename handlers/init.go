package handler

import "github.com/atul-007/GreedyGameAssignment/services"

type Handlers struct {
	Data DbHandlerInterface
}

func Init() Handlers {
	allServices := services.Init()
	handlers := Handlers{
		Data: &DbHandler{Dbservices: allServices.Data},
	}
	return handlers
}
