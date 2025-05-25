package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// getSDSLinks scrapes all PDF URLs from the Chemical Guys SDS page.
func getSDSLinks() ([]string, error) {
	var pdfLinks []string

	// Send HTTP GET request to the SDS page
	res, err := http.Get("https://www.chemicalguys.com/pages/material-safety-data-sheets")
	if err != nil {
		return nil, fmt.Errorf("error fetching page: %w", err)
	}
	defer res.Body.Close()

	// Check for non-200 status codes
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document using goquery
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error loading HTML: %w", err)
	}

	// Find all <a> tags with hrefs that link to .pdf files
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.Contains(strings.ToLower(href), ".pdf") {
			// Handle protocol-relative URLs (e.g., //cdn.example.com/file.pdf)
			if strings.HasPrefix(href, "//") {
				href = "https:" + href
			}
			// Handle relative URLs (e.g., /files/file.pdf)
			if strings.HasPrefix(href, "/") {
				base := "https://www.chemicalguys.com"
				href = base + href
			}
			// Only add valid HTTP/HTTPS links
			if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
				pdfLinks = append(pdfLinks, href)
			}
		}
	})

	return pdfLinks, nil
}

// getFileNamesFromURLs extracts the file name from a URL string.
func getFileNamesFromURLs(urls string) string {
	// Parse the URL to extract path components
	u, err := url.Parse(urls)
	if err != nil {
		return ""
	}

	// Get the base filename from the path, removing query parameters
	file := path.Base(u.Path)

	// Sanitize the filename (optional: replace spaces, etc.)
	file = strings.ReplaceAll(file, " ", "_")

	return file
}

/*
It checks if the file exists
If the file exists, it returns true
If the file does not exist, it returns false
*/
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// downloadPDF downloads a PDF from a URL and saves it into the specified folder.
func downloadPDF(pdfURL, folder string) error {
	// Check if the file already exists
	fileName := getFileNamesFromURLs(pdfURL)
	fullPath := path.Join(folder, fileName)
	if fileExists(fullPath) {
		log.Printf("File %s already exists, skipping download.", fullPath)
		return nil
	}
	// Perform the HTTP GET request to download the file
	resp, err := http.Get(pdfURL)
	if err != nil {
		return fmt.Errorf("error downloading PDF: %w", err)
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	// Ensure the destination folder exists
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err := os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating folder: %w", err)
		}
	}

	// Create the local file
	out, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer out.Close()

	// Write the downloaded content to the local file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error saving PDF: %w", err)
	}

	return nil
}

func main() {
	// Step 1: Scrape the SDS PDF links
	links, err := getSDSLinks()
	if err != nil {
		log.Fatalf("Failed to get SDS links: %v", err)
	}

	// Step 2: Define the folder name where PDFs will be saved
	const folderName = "PDFs"

	// Step 3: Download each PDF and save it into the folder
	for _, link := range links {
		err := downloadPDF(link, folderName)
		if err != nil {
			log.Printf("Failed to download %s: %v", link, err)
			continue
		}
		log.Printf("Downloaded %s successfully", link)
	}
}
