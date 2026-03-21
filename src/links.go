package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

var linksCmd = &cobra.Command{
	Use:   "links [url]",
	Short: "Extract all links from a webpage",
	Args:  cobra.ExactArgs(1),
	Run:   runLinks,
}

func runLinks(cmd *cobra.Command, args []string) {
	targetURL := args[0]
	links, err := extractAllLinks(targetURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	result := map[string]interface{}{
		"url":   targetURL,
		"count": len(links),
		"links": links,
	}

	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}

func extractAllLinks(targetURL string) ([]string, error) {
	baseURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	doc, err := fetchDocument(targetURL)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var allLinks []string

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		link := resolveLink(href, baseURL)
		if link == "" || seen[link] {
			return
		}

		seen[link] = true
		allLinks = append(allLinks, link)
	})

	log.Printf("Found %d links from %s", len(allLinks), targetURL)
	return allLinks, nil
}

func resolveLink(href string, base *url.URL) string {
	href = strings.TrimSpace(href)
	if href == "" || strings.HasPrefix(href, "#") ||
		strings.HasPrefix(href, "javascript:") ||
		strings.HasPrefix(href, "mailto:") {
		return ""
	}

	link, err := url.Parse(href)
	if err != nil {
		return ""
	}

	return base.ResolveReference(link).String()
}
