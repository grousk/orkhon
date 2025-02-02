// Package main
package main

import (
	"fmt"

	"orkhon/internal/log"

	"github.com/alexflint/go-arg"
	"github.com/rs/zerolog"
)

// Version hold the version of the program
// Should set at build time with ldflags
var Version string

type cli struct {
	// Subcommands
	Separate *separateArgs `arg:"subcommand:seperate|s"`
	List     *listArgs     `arg:"subcommand:list|l"`

	// "Global" (not really) flags
	Debug bool `arg:"-d,--debug" help:"show debug logs"`
}

func (cli) Version() string {
	return Version
}

// Printed in help text
func (cli) Description() string {
	return "orkhon is a CLI tool to manage Minecraft mods."
}

// Printed in help text
func (cli) Epilogue() string {
	return fmt.Sprintf("MIT License, Copyright © 2025 ÖMER FERHAT ŞENEL")
}

var args cli

func main() {
	arg.MustParse(&args)

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log := log.New()

	switch {
	case args.Separate != nil:
		args.Separate.Run(log)
	case args.List != nil:
		args.List.Run(log)
	}
}
