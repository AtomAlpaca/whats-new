package main

import (
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "whats-new",
		Short: "A tool to fetch and analyze web content",
	}

	rootCmd.AddCommand(linksCmd)
	rootCmd.AddCommand(extractCmd)
	rootCmd.AddCommand(fullTextCmd)
	rootCmd.AddCommand(metadataCmd)
	rootCmd.AddCommand(memoryCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
