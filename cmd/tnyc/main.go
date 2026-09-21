// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/mdhender/tnyc"
	"github.com/mdhender/tnyc/internal/dotenv"
	"github.com/spf13/cobra"
)

func main() {
	// TNYC_ENV selects which dotenv files load, and scopes the credential file.
	// It is read before flag parsing because those files populate the
	// environment ff then reads.
	env, ok := os.LookupEnv("TNYC_ENV")
	if !ok {
		env = "development"
	}
	// Load also rejects an unknown environment, which is what lets env be used
	// as a path segment in the credential file without further checking.
	if err := dotenv.Load(env); err != nil {
		fmt.Fprintf(os.Stderr, "tnyc: %v\n", err)
		os.Exit(1)
	}

	if err := run(context.Background(), os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "tnyc: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	cmd := newRootCommand()
	cmd.SetArgs(args)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	return cmd.ExecuteContext(ctx)
}

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "tnyc",
		Short:         "An unfaithful remake of a lost classic",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Fprintln(cmd.OutOrStdout(), tnyc.Version())
		},
	})
	cmd.AddCommand(newWorldCommand())
	cmd.AddCommand(newDatabaseCommand())
	cmd.AddCommand(&cobra.Command{
		Use:   "game",
		Short: "Run the game",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "api",
		Short: "Run the API server",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	})

	return cmd
}

func newDatabaseCommand() *cobra.Command {
	const (
		defaultDatabasePath = "var/"
		defaultWorldMap     = "var/tnyc-world.json"
	)
	cmd := &cobra.Command{
		Use:   "database",
		Short: "Manage the database",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	var dbPath, worldMap string
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a database",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := tnyc.CreateDatastore(dbPath, worldMap); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "created database at %s\n", dbPath)
			return nil
		},
	}
	createCmd.Flags().StringVar(&dbPath, "db-path", defaultDatabasePath, "existing directory for tnyc.json")
	createCmd.Flags().StringVar(&worldMap, "world-map", defaultWorldMap, "T'Nyc world JSON path")
	cmd.AddCommand(createCmd)
	return cmd
}

func newWorldCommand() *cobra.Command {
	const (
		defaultInputPath  = "var/wgvc-export.json"
		defaultOutputPath = "var/tnyc-world.json"
	)
	var inputPath, outputPath string
	cmd := &cobra.Command{
		Use:   "world",
		Short: "Import and validate a WGVC world",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			world, err := tnyc.LoadWorld(inputPath)
			if err != nil {
				return err
			}
			if err := tnyc.SaveWorld(outputPath, world); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "saved world to %s\n", outputPath)
			return nil
		},
	}
	cmd.Flags().StringVarP(&inputPath, "input", "i", defaultInputPath, "WGVC schema-v1 JSON input path")
	cmd.Flags().StringVarP(&outputPath, "output", "o", defaultOutputPath, "T'Nyc world JSON output path")
	return cmd
}
