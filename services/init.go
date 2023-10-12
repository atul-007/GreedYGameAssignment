package services

import "github.com/atul-007/GreedyGameAssignment/models"

type Services struct {
	Data DbServicesInterface
}

func Init() Services {

	dataService := DbServices{store: make(map[string]models.KeyValue)}

	services := Services{
		Data: &dataService,
	}
	return services
}
