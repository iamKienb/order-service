package module

import (
	"inventory-query-module/internal/application/services/inventory"
)

type ApplicationModule struct {
}

func NewApplicationModule(infra *InfraModule) *ApplicationModule {
	inventoryQueryService := inventory.NewQueryService(infra.ESService)

}
