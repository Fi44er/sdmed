package user_entity

type Role struct {
	ID          string
	Name        string
	Premissions []string
}

func (r *Role) AddPermission(permission string) error {
	r.Premissions = append(r.Premissions, permission)
	return nil
}

func (r *Role) CanAccess(permission string) bool {
	for _, p := range r.Premissions {
		if p == permission {
			return true
		}
	}
	return false
}
