package auth

// Roles handled ONLY by auth service
const (
	RoleUser  = "USER"  // Default authenticated user
	RoleAdmin = "ADMIN" // Platform administrator
)

// ValidRoles allowed in auth service
var ValidRoles = map[string]bool{
	RoleUser:  true,
	RoleAdmin: true,
}
