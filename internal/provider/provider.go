// Package provider provides a interface for data provider
package provider

// Provider defines methods for interacting with mod data.
type Provider interface {
	FetchMod(id string) (ModDTO, error)
}
