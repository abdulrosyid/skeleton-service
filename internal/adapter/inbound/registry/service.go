package registry

import (
	cfg "skeleton-service/configs"
	repoReg "skeleton-service/internal/adapter/outbound/repository/mysqldb"
	csvc "skeleton-service/internal/core/service"
	iservice "skeleton-service/internal/port/inbound/service"
)

type ServiceRegistry struct {
	repositoryRegistry repoReg.RepoSQL
	authService        iservice.AuthService
	jwtSecret          string
}

func NewServiceRegistry(r repoReg.RepoSQL, c cfg.Config) *ServiceRegistry {
	return &ServiceRegistry{
		repositoryRegistry: r,
		authService:        csvc.NewAuthService(r, c.JWTSecret, c.JWTExpireMinutes),
		jwtSecret:          c.JWTSecret,
	}
}

func (r *ServiceRegistry) GetAuthService() iservice.AuthService {
	return r.authService
}

func (r *ServiceRegistry) GetJWTSecret() string {
	return r.jwtSecret
}
