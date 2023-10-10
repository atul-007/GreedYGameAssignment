package handler

import "github.com/atul-007/GreedyGameAssignment/services"

type DbHandlerInterface interface {
}
type DbHandler struct {
	Dbservices services.DbServicesInterface
}
