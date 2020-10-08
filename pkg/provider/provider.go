package provider

import "github.com/ryarnyah/dblock/pkg/model"

// RegistredProviders provider available for sourcing
var RegistredProviders = make(map[string]Provider)

// RegisterProvider permit to add new provider
func RegisterProvider(name string, provider Provider) {
	RegistredProviders[name] = provider
}

// Provider schema provider
type Provider interface {
	GetCurrentModel() (*model.DatabaseSchema, error)
}
