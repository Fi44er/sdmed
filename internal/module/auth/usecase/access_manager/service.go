package accessmanager_service

import (
	"context"
	"strings"

	auth_adapters "github.com/Fi44er/sdmed/internal/module/auth/infrastucture/adapters"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type Manager struct {
	Enforcer           *casbin.Enforcer
	db                 *gorm.DB
	userUsecaseAdapter *auth_adapters.UserUsecaseAdapter
}

const rbacModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`

func NewManager(db *gorm.DB, userUsecaseAdapter *auth_adapters.UserUsecaseAdapter) (*Manager, error) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	m, err := model.NewModelFromString(rbacModel)
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	return &Manager{Enforcer: enforcer, db: db, userUsecaseAdapter: userUsecaseAdapter}, nil
}

func (m *Manager) SyncRolePermissions(ctx context.Context) error {
	roles, err := m.userUsecaseAdapter.GetAllRoles(ctx)
	if err != nil {
		return err
	}

	m.Enforcer.ClearPolicy()

	for _, role := range roles {
		for _, perm := range role.Permissions {
			parts := strings.Split(perm.Name, ":")
			if len(parts) == 2 {
				_, _ = m.Enforcer.AddPolicy(role.Name, parts[0], parts[1])
			}
		}
	}
	return m.Enforcer.SavePolicy()
}
