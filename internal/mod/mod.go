// Package mod provides Mod struct for mod data related stuff.
package mod

import (
	"encoding/json"
	"fmt"
	"orkhon/internal/provider"
	"os"

	"github.com/SladeThe/yav"
	"github.com/SladeThe/yav/vstring"
)

// Mod provides whole mod data.
type Mod struct {
	ID       string `json:"id"`
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Platform string `json:"platform"`
	// Provider is used for determining which data provider is used.
	Provider provider.Provider
}

// Validate validates the Mod struct and returns an error.
func (m *Mod) Validate() error {
	return yav.Join(
		yav.Chain(
			"name", m.Name,
			vstring.Required,
		),
		yav.Chain(
			"path", m.Path,
			vstring.Required,
		),
	)
}

// LoadFromFile parses the modlist file
func LoadFromFile(filep string) ([]*Mod, error) {
	file, err := os.Open(filep)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var mods []*Mod
	if err := json.NewDecoder(file).Decode(&mods); err != nil {
		return nil, err
	}

	return mods, nil
}

// UnmarshalJSON custom Unmarshaler for encoding/json
func (m *Mod) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	var temp map[string]interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	// Mapping "Identifier" field
	if id, ok := temp["id"].(string); ok {
		m.ID = id
	} else {
		return fmt.Errorf("invalid or missing 'id' field")
	}
	if slug, ok := temp["slug"].(string); ok {
		m.Slug = slug
	} else {
		return fmt.Errorf("invalid or missing 'slug' field")
	}
	if name, ok := temp["name"].(string); ok {
		m.Name = name
	}
	if path, ok := temp["path"].(string); ok {
		m.Path = path
	}
	return nil
}

// WithProvider sets Provider field of Mod
func (m *Mod) WithProvider(p provider.Provider) *Mod {
	m.Provider = p
	return m
}

// UpdatePlatform fetchs platform data from provider and changes its Platform value with it
func (m *Mod) UpdatePlatform() error {
	dto, err := m.Provider.FetchMod(m.ID)
	if err != nil {
		return err
	}

	m.Platform = dto.Platform

	return nil
}
