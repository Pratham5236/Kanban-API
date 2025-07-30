
package auth

import (
	"kanban-app/api/database"

	"github.com/casbin/casbin/v2"
)

type Service struct {
	enforcer *casbin.Enforcer
}

func NewAuthorizationService() *Service {
	return &Service{
		enforcer: database.Enforcer,
	}
}

func (s *Service) Enforce(sub, obj, act string) (bool, error) {
	return s.enforcer.Enforce(sub, obj, act)
}

func (s *Service) AddPolicy(sub, obj, act string) (bool, error) {
	return s.enforcer.AddPolicy(sub, obj, act)
}

func (s *Service) RemovePolicy(sub, obj, act string) (bool, error) {
	return s.enforcer.RemovePolicy(sub, obj, act)
}
