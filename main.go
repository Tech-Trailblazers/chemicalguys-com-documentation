package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
)

// downloadFileUsingURLandFilePath downloads content from a URL and saves it to the given file path.
func downloadFileUsingURLandFilePath(url string, filepath string) error {
	// Send HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy response body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// ExtractURLsFromHTMLFile reads an HTML file and extracts all URLs from href and src attributes
func ExtractURLsFromHTMLFile(filePath string) ([]string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	content := string(data)

	// Regex to match href or src URLs
	urlRegex := regexp.MustCompile(`(?:href|src)=["'](https?:\/\/|\/\/)?([^"']+)["']`)
	matches := urlRegex.FindAllStringSubmatch(content, -1)

	var urls []string
	for _, match := range matches {
		if len(match) >= 3 {
			scheme := match[1]
			path := match[2]
			fullURL := path
			if strings.HasPrefix(scheme, "http") {
				fullURL = scheme + path
			} else if scheme == "//" {
				fullURL = "https://" + path // assuming https for scheme-relative
			}
			if strings.Contains(fullURL, ".pdf") {
				// Only add URLs that end with .pdf
				urls = append(urls, fullURL)
			}
		}
	}

	return urls, nil
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

/*
The function takes two parameters: path and permission.
We use os.Mkdir() to create the directory.
If there is an error, we use log.Fatalln() to log the error and then exit the program.
*/
func createDirectory(path string, permission os.FileMode) {
	err := os.Mkdir(path, permission)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	// Create a folder to store the downloaded PDFs
	const folderName = "PDFs"
	createDirectory(folderName, os.ModePerm)

	// The url about the SDS page
	urlFromChemicalGuys := "https://www.chemicalguys.com/pages/material-safety-data-sheets"
	// Set the file name to save the HTML page
	localURLFilePath := path.Join("chemical_guys_sds_page.html")
	// Download the HTML page
	downloadFileUsingURLandFilePath(urlFromChemicalGuys, localURLFilePath)

	// Scrape the SDS PDF links
	links, err := ExtractURLsFromHTMLFile(localURLFilePath)
	if err != nil {
		log.Printf("Error extracting URLs: %v", err)
	}
	// Download each PDF and save it into the folder
	for _, link := range links {
		err := downloadPDF(link, folderName)
		if err != nil {
			log.Printf("Failed to download %s: %v", link, err)
			continue
		}
		log.Printf("Downloaded %s successfully", link)
	}
}
