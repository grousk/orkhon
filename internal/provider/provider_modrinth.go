// Package provider provides a interface for data provider
package provider

import (
	"fmt"
	"net/http"

	"resty.dev/v3"
)

// ModrinthProvider interacts with the Modrinth API
type ModrinthProvider struct {
	client *resty.Client
}

// NewModrinthProvider creates a new ModrinthProvider
func NewModrinthProvider(c *resty.Client) *ModrinthProvider {
	return &ModrinthProvider{
		client: c,
	}
}

// FetchMod retrieves mod data from Modrinth API
func (p *ModrinthProvider) FetchMod(id string) (ModDTO, error) {
	url := fmt.Sprintf("https://api.modrinth.com/v2/project/%s", id)

	var dto ModDTO

	res, err := p.client.R().SetResult(&dto).Get(url)
	if err != nil {
		return ModDTO{}, err
	}

	if res.StatusCode() != http.StatusOK {
		return ModDTO{}, fmt.Errorf("failed to fetch mod: %d", res.StatusCode())
	}

	return dto, nil
}
