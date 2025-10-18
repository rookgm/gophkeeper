package main

import "github.com/spf13/cobra"

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync data",
	RunE:  runSync,
}

func runSync(cmd *cobra.Command, args []string) error {
	return nil
}
