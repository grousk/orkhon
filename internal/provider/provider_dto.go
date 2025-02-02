// Package provider provides a interface for data provider
package provider

import "encoding/json"

// ModDTO hold data that will be used later in processes.
type ModDTO struct {
	Platform string `json:"platform"`
}

// UnmarshalJSON custom Unmarshaler for encoding/json
func (m *ModDTO) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Mapping "Platform" field
	if clientSide, ok := temp["client_side"].(string); ok && (clientSide == "required" || clientSide == "optional") {
		m.Platform = "client"
	}

	if serverSide, ok := temp["server_side"].(string); ok && (serverSide == "required" || serverSide == "optional") {
		if m.Platform == "client" {
			m.Platform = "both" // If both are required, set to "BOTH"
		} else {
			m.Platform = "server"
		}
	}

	return nil
}
