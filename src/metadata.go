package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

var metadataCmd = &cobra.Command{
	Use:   "metadata [url]",
	Short: "Extract metadata from a webpage",
	Args:  cobra.ExactArgs(1),
	Run:   runMetadata,
}

func runMetadata(cmd *cobra.Command, args []string) {
	targetURL := args[0]
	metadata, err := extractMetadata(targetURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	output, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}

func extractMetadata(targetURL string) (map[string]string, error) {
	baseURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	doc, err := fetchDocument(targetURL)
	if err != nil {
		return nil, err
	}

	metadata := make(map[string]string)

	metadata["url"] = targetURL
	metadata["title"] = doc.Find("title").Text()

	metadata["description"] = getMetaContent(doc, "description")
	metadata["keywords"] = getMetaContent(doc, "keywords")
	metadata["author"] = getMetaContent(doc, "author")

	metadata["og:title"] = getOpenGraph(doc, "og:title")
	metadata["og:description"] = getOpenGraph(doc, "og:description")
	metadata["og:image"] = getOpenGraph(doc, "og:image")
	metadata["og:url"] = getOpenGraph(doc, "og:url")

	favicon := getFavicon(doc, baseURL)
	if favicon != "" {
		metadata["favicon"] = favicon
	}

	metadata["site_name"] = getOpenGraph(doc, "og:site_name")

	cleanMetadata(metadata)
	return metadata, nil
}

func getMetaContent(doc *goquery.Document, name string) string {
	var content string
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if attr, _ := s.Attr("name"); attr == name {
			if c, ok := s.Attr("content"); ok {
				content = c
			}
		}
	})
	return content
}

func getOpenGraph(doc *goquery.Document, property string) string {
	var content string
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if attr, _ := s.Attr("property"); attr == property {
			if c, ok := s.Attr("content"); ok {
				content = c
			}
		}
	})
	return content
}

func getFavicon(doc *goquery.Document, base *url.URL) string {
	var favicon string
	doc.Find("link[rel='icon'], link[rel='shortcut icon']").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok && favicon == "" {
			faviconURL, err := url.Parse(href)
			if err == nil {
				if faviconURL.IsAbs() {
					favicon = href
				} else {
					favicon = base.ResolveReference(faviconURL).String()
				}
			}
		}
	})
	return favicon
}

func cleanMetadata(m map[string]string) {
	for k, v := range m {
		if v == "" {
			delete(m, k)
		}
	}
}
