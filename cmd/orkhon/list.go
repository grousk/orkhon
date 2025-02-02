// Package main
package main

import (
	"encoding/json"
	"fmt"
	"orkhon/internal/log"
	"orkhon/internal/mod"

	"github.com/SladeThe/yav"
	"github.com/SladeThe/yav/vstring"
)

type listArgs struct {
	ModsFile string `arg:"-f,--" default:"mods.json"`
}

func (a *listArgs) Validate() error {
	return yav.Chain(
		"mods_file", a.ModsFile,
		vstring.Required,
	)
}

func (a *listArgs) Run(log *log.Logger) {
	if err := a.Validate(); err != nil {
		log.Fatal().Err(err).Msg("Invalid arguments")
	}

	mods, err := mod.LoadFromFile(a.ModsFile)
	if err != nil {
		log.Error().Err(err).Msg("Unable parse modlist")
	}
	log.Info().Msgf("loaded %d mods", len(mods))

	for _, mod := range mods {
		formatted, err := formatMod(mod)
		if err != nil {
			log.Warn().Err(err).Msg("Unable to indent mod")
		}

		fmt.Println(formatted)
	}
}

func formatMod(mod *mod.Mod) (string, error) {
	// Marshal to json with identiation (with tab)
	formatted, err := json.MarshalIndent(mod, "", "\t")
	if err != nil {
		return "", err
	}

	return string(formatted), err
}
