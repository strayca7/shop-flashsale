package admindomain

type Permission string

type Role struct {
    name string
    permissions map[Permission]bool
}

func (r *Role) HasPermission(p Permission) bool {
	return false
}