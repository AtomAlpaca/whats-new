package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

var extractCmd = &cobra.Command{
	Use:   "extract [url]",
	Short: "Extract main content from a webpage as plain text",
	Args:  cobra.ExactArgs(1),
	Run:   runExtract,
}

func runExtract(cmd *cobra.Command, args []string) {
	targetURL := args[0]
	content, err := extractContent(targetURL)
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

func extractContent(targetURL string) (string, error) {
	baseURL, err := url.Parse(targetURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	doc, err := fetchDocument(targetURL)
	if err != nil {
		return "", err
	}

	removeUnwanted(doc)
	selection := findMainContent(doc)
	text := extractText(selection)

	log.Printf("Extracted %d characters from %s", len(text), targetURL)
	_ = baseURL
	return text, nil
}

func removeUnwanted(doc *goquery.Document) {
	doc.Find("script, style, nav, header, footer, aside, .advertisement, .ad, .comments, .social, .share, .menu, .nav, .breadcrumb, .related").Remove()
}

func findMainContent(doc *goquery.Document) *goquery.Selection {
	selectors := []string{
		"article",
		"[role='main']",
		"main",
		".content",
		".post-content",
		".article-content",
		".entry-content",
		"#content",
		".text-content",
		".main-content",
	}

	for _, sel := range selectors {
		if doc.Find(sel).Length() > 0 {
			return doc.Find(sel).First()
		}
	}

	return doc.Find("body")
}

func extractText(sel *goquery.Selection) string {
	var paragraphs []string

	sel.Find("p, h1, h2, h3, h4, h5, h6, li, blockquote, pre, div").Each(func(i int, s *goquery.Selection) {
		text := cleanTextReturns(s.Text())
		if isMeaningful(text) {
			tagName := getTagName(s)
			if tagName == "p" || tagName == "li" || tagName == "blockquote" {
				paragraphs = append(paragraphs, text)
			} else if isHeading(tagName) {
				paragraphs = append(paragraphs, "\n## "+text+"\n")
			} else {
				paragraphs = append(paragraphs, text)
			}
		}
	})

	if len(paragraphs) == 0 {
		return cleanTextReturns(sel.Text())
	}

	result := strings.Join(paragraphs, "\n\n")
	result = regexp.MustCompile(`\n{3,}`).ReplaceAllString(result, "\n\n")
	return strings.TrimSpace(result)
}

func getTagName(s *goquery.Selection) string {
	return s.Nodes[0].Data
}

func isHeading(tag string) bool {
	return tag == "h1" || tag == "h2" || tag == "h3" || tag == "h4" || tag == "h5" || tag == "h6"
}

func cleanTextReturns(text string) string {
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)
	return text
}

func isMeaningful(text string) bool {
	text = strings.TrimSpace(text)
	if len(text) < 20 {
		return false
	}
	noisePatterns := []string{
		"^follow us",
		"^share on",
		"^read more",
		"^learn more",
		"^click here",
		"^subscribe",
		"^copyright",
		"^all rights reserved",
		"^privacy policy",
		"^terms of",
		"^[a-z]+@[a-z]+\\.[a-z]+",
		"^\\d+$",
	}
	lower := strings.ToLower(text)
	for _, pattern := range noisePatterns {
		if matched, _ := regexp.MatchString(pattern, lower); matched {
			return false
		}
	}
	return true
}
