package services

import (
	"context"
	"fmt"
	"strings"

	"app/internal/core/model"
)

type Store interface {
	GetEmployeesByName(
		ctx context.Context,
		name string,
		offsetOpts model.OffsetRequest,
	) ([]model.Employee, error)
}

type Services struct {
	s Store
}

func NewServices(s Store) *Services {
	return &Services{s}
}

func (s *Services) GetEmployeesByName(
	ctx context.Context,
	req model.GetEmployeesRequest,
) ([]model.Employee, error) {
	name := strings.ToLower(req.Name)
	emps, err := s.s.GetEmployeesByName(ctx, name, req.OffsetRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the employees list from store: %w", err)
	}
	return emps, nil
}
