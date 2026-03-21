package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

var fullTextCmd = &cobra.Command{
	Use:   "full-text [url]",
	Short: "Get full HTML content from a webpage",
	Args:  cobra.ExactArgs(1),
	Run:   runFullText,
}

func runFullText(cmd *cobra.Command, args []string) {
	targetURL := args[0]
	content, err := getFullHTML(targetURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	result := map[string]interface{}{
		"url":     targetURL,
		"content": content,
	}

	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}

func getFullHTML(targetURL string) (string, error) {
	_, err := url.Parse(targetURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	doc, err := fetchDocument(targetURL)
	if err != nil {
		return "", err
	}

	html, err := doc.Html()
	if err != nil {
		return "", fmt.Errorf("failed to get HTML: %w", err)
	}

	return html, nil
}
