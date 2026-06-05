package module

type ApplicationModule struct{}

func NewApplicationModule(infra *InfraModule) *ApplicationModule {
	return &ApplicationModule{}
}
