package application

import (
	"context"

	"app/internal/core/model"
	"app/internal/core/services"
)

type Application struct {
	services *services.Services
}

func NewApplication(store services.Store) *Application {
	return &Application{
		services: services.NewServices(store),
	}
}

func (app *Application) GetEmployeesByName(
	ctx context.Context,
	req model.GetEmployeesRequest,
) ([]model.Employee, error) {
	return app.services.GetEmployeesByName(ctx, req)
}
