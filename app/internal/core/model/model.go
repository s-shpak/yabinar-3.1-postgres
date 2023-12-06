package model

type Employee struct {
	ID         int    `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Salary     int    `json:"salary"`
	PositionID int    `json:"position_id"`
	Email      string `json:"email"`
}

type OffsetRequest struct {
	Limit  int
	LastID int
}

type GetEmployeesRequest struct {
	OffsetRequest
	Name string
}
