// Package main
package main

import (
	"fmt"
	"path/filepath"
	"time"

	"orkhon/internal/log"
	"orkhon/internal/mod"
	"orkhon/internal/provider"
	"orkhon/pkg/file"

	"github.com/SladeThe/yav"
	"github.com/SladeThe/yav/vstring"
	"golang.org/x/sync/errgroup"
	"resty.dev/v3"
)

// separateArgs is a arguments struct for seperate subcommand
type separateArgs struct {
	ModsFile      string `arg:"-f,--" default:"mods.json"`
	ClientModsDir string `arg:"-c,--client-dir" placeholder:"DIR"`
	ServerModsDir string `arg:"-s,--server-dir" placeholder:"DIR"`
}

// Validate validates seperateArgs
func (a *separateArgs) Validate() error {
	return yav.Join(
		yav.Chain(
			"mods_file", a.ModsFile,
			vstring.Required,
		),
		yav.Chain(
			"client_mods_dir", a.ClientModsDir,
			vstring.Required,
		),
		yav.Chain(
			"server_mods_dir", a.ServerModsDir,
			vstring.Required,
		),
	)
}

func (a *separateArgs) Run(log *log.Logger) {
	if err := a.Validate(); err != nil {
		log.Fatal().Err(err).Msg("Invalid arguments")
	}

	log.Info().Str("file_path", a.ModsFile).Msg("Loading mods from file")
	mods, err := mod.LoadFromFile(a.ModsFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to parse modlist")
	}
	log.Info().Msgf("loaded %d mods", len(mods))

	var eg errgroup.Group

	// TODO: add a arg to let user change this value
	eg.SetLimit(50)
	log.Debug().Int("limit", 50).Msg("Created a ErrGroup")

	// Fetching mod data
	client := resty.New().
		SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(15 * time.Second)
	defer client.Close()

	provider := provider.NewModrinthProvider(client)

	log.Info().Msg("Fetching mod data")
	for _, mod := range mods {
		fetchMod(mod, provider, &eg)
	}

	// Seperate mods
	if err := eg.Wait(); err != nil {
		log.Warn().Err(err).Msg("Unable to fetch data")
	}

	log.Info().Msg("Seperating mods")
	if err := separateMods(mods, a.ClientModsDir, a.ServerModsDir); err != nil {
		log.Error().Err(err).Msg("Unable to seperate mods")
	}

	log.Info().Msg("Done!")
}

func fetchMod(mod *mod.Mod, p provider.Provider, eg *errgroup.Group) {
	modCopy := mod

	eg.Go(func() error {
		if err := modCopy.WithProvider(p).UpdatePlatform(); err != nil {
			return err
		}
		return nil
	})
}

func separateMods(mods []*mod.Mod, clientDir, serverDir string) error {
	for _, mod := range mods {
		modPath := filepath.Base(mod.Path)
		cdir, sdir := filepath.Join(clientDir, modPath), filepath.Join(serverDir, modPath)

		switch mod.Platform {
		case "both":
			// NOTE: might cause slower performance
			if err := file.Copy(mod.Path, cdir); err != nil {
				return err
			}
			if err := file.Copy(mod.Path, sdir); err != nil {
				return err
			}
		case "client":
			if err := file.Copy(mod.Path, cdir); err != nil {
				return err
			}
		case "server":
			if err := file.Copy(mod.Path, sdir); err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid platform for mod: %s", mod.Name)
		}
	}

	return nil
}
