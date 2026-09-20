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

	return cmd
}
