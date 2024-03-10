package account

import uuid "github.com/satori/go.uuid"

type Role string

const (
	Employee   Role = "employee"
	Manager    Role = "manager"
	Accountant Role = "accountant"
	Admin      Role = "admin"
)

type Account struct {
	PublicID uuid.UUID `json:"public_id"`

	Name  string `json:"name"`
	Email string `json:"email"`

	Role Role `json:"role"`
}
