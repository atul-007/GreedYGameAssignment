package services

type Services struct {
	Data DbServicesInterface
}

func Init() Services {

	dataService := DbServices{}

	services := Services{
		Data: &dataService,
	}
	return services
}
