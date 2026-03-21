package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var memoryCmd = &cobra.Command{
	Use:   "memory",
	Short: "Manage website memory",
}

var memoryGetCmd = &cobra.Command{
	Use:   "get [url]",
	Short: "Get website from database",
	Args:  cobra.ExactArgs(1),
	Run:   runMemoryGet,
}

func init() {
	memoryCmd.AddCommand(memoryGetCmd)
	memoryCmd.AddCommand(saveMemoryCmd)
}

var saveMemoryCmd = &cobra.Command{
	Use:   "save [json]",
	Short: "Save website to database",
	Args:  cobra.ExactArgs(1),
	Run:   runMemory,
}

type WebsiteInput struct {
	URL    string `json:"url"`
	Title  string `json:"title"`
	Report string `json:"report"`
}

func runMemoryGet(cmd *cobra.Command, args []string) {
	ensureDB()

	record, err := getWebsite(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if record == nil {
		fmt.Fprintf(os.Stderr, "Error: website not found\n")
		os.Exit(1)
	}

	output, _ := json.MarshalIndent(record, "", "  ")
	fmt.Println(string(output))
}

func runMemory(cmd *cobra.Command, args []string) {
	var input WebsiteInput

	err := json.Unmarshal([]byte(args[0]), &input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid JSON: %v\n", err)
		os.Exit(1)
	}

	if input.URL == "" {
		fmt.Fprintln(os.Stderr, "Error: url is required")
		os.Exit(1)
	}

	if input.Title == "" {
		metadata, err := extractMetadata(input.URL)
		if err != nil {
			input.Title = ""
		} else {
			if t, ok := metadata["title"]; ok {
				input.Title = t
			}
		}
	}

	if input.Report == "" {
		htmlContent, err := getFullHTML(input.URL)
		if err != nil {
			input.Report = ""
		} else {
			lines := strings.Split(htmlContent, "\n")
			var summary []string
			for i, line := range lines {
				if i >= 3 {
					break
				}
				line = strings.TrimSpace(line)
				if len(line) > 100 {
					line = line[:100] + "..."
				}
				if line != "" {
					summary = append(summary, line)
				}
			}
			input.Report = strings.Join(summary, "\n")
		}
	}

	htmlContent, err := getFullHTML(input.URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to fetch HTML: %v\n", err)
		os.Exit(1)
	}

	ensureDB()

	err = saveWebsite(input.URL, input.Title, htmlContent, input.Report)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	record, err := getWebsite(input.URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to get record: %v\n", err)
	}

	result := map[string]interface{}{
		"url":   input.URL,
		"title": input.Title,
		"saved": true,
	}

	if record != nil {
		result["content_hash"] = record["content_hash"]
		result["last_recorded"] = record["last_recorded"]
	}

	output, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(output))
}
