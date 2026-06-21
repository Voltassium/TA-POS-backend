package constants

type (
	UserRole    string
	UserStatus  string
	UserRoleMap map[UserRole]bool
)

const (
	UserRoleSuperadmin UserRole = "superadmin"
	UserRoleOwner      UserRole = "owner"
	UserRoleChef       UserRole = "chef"
	UserRoleStaff      UserRole = "staff"
)

const (
	UserStatusActive = "Active"
)

func (receiver UserRole) IsValidEnum() bool {
	switch receiver {
	case UserRoleSuperadmin, UserRoleOwner, UserRoleChef, UserRoleStaff:
		return true
	default:
		return false
	}
}

func (receiver UserStatus) IsValidEnum() bool {
	switch receiver {
	case UserStatusActive:
		return true
	default:
		return false
	}
}
