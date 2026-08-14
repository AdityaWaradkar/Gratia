package user

type Role string

const (
	RoleDonor Role = "DONOR"
	RoleNGO   Role = "NGO"
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "USER"
)