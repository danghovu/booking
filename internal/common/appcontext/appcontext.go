package appcontext

import "booking-event/config"

type AppContext interface {
	InfraRegistry() InfraRegistry
	RepositoryRegistry() RepositoryRegistry
	ServiceRegistry() ServiceRegistry
}

type appContext struct {
	infraRegistry      InfraRegistry
	repositoryRegistry RepositoryRegistry
	serviceRegistry    ServiceRegistry
}

func NewAppContext(
	config config.Config,
) AppContext {
	infraRegistry := NewInfraRegistry(config)
	repositoryRegistry := NewRepositoryRegistry(config, infraRegistry)
	serviceRegistry := NewServiceRegistry(config, infraRegistry, repositoryRegistry)
	return &appContext{
		infraRegistry:      infraRegistry,
		repositoryRegistry: repositoryRegistry,
		serviceRegistry:    serviceRegistry,
	}
}

func (a *appContext) InfraRegistry() InfraRegistry {
	return a.infraRegistry
}

func (a *appContext) RepositoryRegistry() RepositoryRegistry {
	return a.repositoryRegistry
}

func (a *appContext) ServiceRegistry() ServiceRegistry {
	return a.serviceRegistry
}
