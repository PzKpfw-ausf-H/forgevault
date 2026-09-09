package domain

type Role string

const (
	RoleStudent    Role = "student"
	RoleEmployee   Role = "employee"
	RoleAdmin      Role = "admin"
	RoleSuperadmin Role = "superadmin"
)

func (r Role) Valid() bool {
	switch r {
	case RoleAdmin, RoleEmployee, RoleStudent, RoleSuperadmin:
		return true
	default:
		return false
	}
}
