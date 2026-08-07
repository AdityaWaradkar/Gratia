package auth

// Role defines a strict custom type for access control levels
type Role string

// Core roles required by the system architecture
const (
    RoleDonor Role = "DONOR"
    RoleNGO   Role = "NGO"
    RoleAdmin Role = "ADMIN"
    RoleUser  Role = "USER"
)

// ValidRoles provides an O(1) lookup map for input sanitization
var ValidRoles = map[Role]bool{
    RoleDonor: true,
    RoleNGO:   true,
    RoleAdmin: true,
    RoleUser:  true,
}