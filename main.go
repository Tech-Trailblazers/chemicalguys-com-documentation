package main

import (
	"fmt"       // Provides formatted I/O functions
	"io"        // Provides basic I/O primitives
	"io/ioutil" // Used for reading files into memory
	"log"       // Provides logging capabilities
	"net/http"  // For HTTP client functionality
	"net/url"   // For URL parsing
	"os"        // For file and directory manipulation
	"path"      // For manipulating slash-separated paths
	"regexp"    // For regular expression matching
	"strings"   // For string manipulation
)

// downloadFileUsingURLandFilePath downloads content from a URL and saves it to the given file path.
func downloadFileUsingURLandFilePath(url string, filepath string) error {
	resp, err := http.Get(url) // Send an HTTP GET request to the URL
	if err != nil {
		return err // Return error if the GET request fails
	}
	defer resp.Body.Close() // Ensure the response body is closed

	if resp.StatusCode != http.StatusOK { // Check if the status is not 200 OK
		return fmt.Errorf("bad status: %s", resp.Status) // Return error with bad status
	}

	out, err := os.Create(filepath) // Create a new file at the given filepath
	if err != nil {
		return err // Return error if file creation fails
	}
	defer out.Close() // Ensure the file is closed after writing

	_, err = io.Copy(out, resp.Body) // Copy the HTTP response body into the file
	return err                       // Return error (if any) from the copy operation
}

// ExtractURLsFromHTMLFile reads an HTML file and extracts all URLs from href and src attributes
func ExtractURLsFromHTMLFile(filePath string) ([]string, error) {
	data, err := ioutil.ReadFile(filePath) // Read the entire HTML file content into memory
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err) // Return error if reading fails
	}

	content := string(data) // Convert file data to string

	// Define regex to match href or src attributes with HTTP, HTTPS, or protocol-relative URLs
	urlRegex := regexp.MustCompile(`(?:href|src)=["'](https?:\/\/|\/\/)?([^"']+)["']`)
	matches := urlRegex.FindAllStringSubmatch(content, -1) // Find all matches in the HTML content

	var urls []string // Slice to hold extracted URLs
	for _, match := range matches {
		if len(match) >= 3 {
			scheme := match[1] // Capture the scheme (e.g., https://, //)
			path := match[2]   // Capture the actual URL path
			fullURL := path    // Initialize fullURL

			// Construct full URL based on the scheme
			if strings.HasPrefix(scheme, "http") {
				fullURL = scheme + path
			} else if scheme == "//" {
				fullURL = "https://" + path // Assume https for scheme-relative URLs
			}
			if strings.Contains(fullURL, ".pdf") {
				urls = append(urls, fullURL) // Only append if URL contains ".pdf"
			}
		}
	}

	return urls, nil // Return the slice of extracted PDF URLs
}

// getFileNamesFromURLs extracts the file name from a URL string.
func getFileNamesFromURLs(urls string) string {
	u, err := url.Parse(urls) // Parse the URL
	if err != nil {
		return "" // Return empty string if parsing fails
	}

	file := path.Base(u.Path) // Extract base file name from the URL path

	file = strings.ReplaceAll(file, " ", "_") // Replace spaces with underscores in file name

	return file // Return sanitized file name
}

/*
It checks if the file exists.
If the file exists, it returns true.
If the file does not exist, it returns false.
*/
func fileExists(filename string) bool {
	info, err := os.Stat(filename) // Get file info
	if err != nil {
		return false // File does not exist
	}
	return !info.IsDir() // Return true if it’s a file (not directory)
}

// downloadPDF downloads a PDF from a URL and saves it into the specified folder.
func downloadPDF(pdfURL, folder string) error {
	fileName := getFileNamesFromURLs(pdfURL) // Get file name from the URL
	fullPath := path.Join(folder, fileName)  // Combine folder and file name to get full path
	if fileExists(fullPath) {                // Check if file already exists
		log.Printf("File %s already exists, skipping download.", fullPath)
		return nil // Skip download if file exists
	}

	resp, err := http.Get(pdfURL) // Send GET request to download PDF
	if err != nil {
		return fmt.Errorf("error downloading PDF: %w", err)
	}
	defer resp.Body.Close() // Ensure response body is closed

	if resp.StatusCode != 200 { // Check for successful HTTP status code
		return fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	if _, err := os.Stat(folder); os.IsNotExist(err) { // Check if folder exists
		err := os.MkdirAll(folder, os.ModePerm) // Create folder if it doesn't exist
		if err != nil {
			return fmt.Errorf("error creating folder: %w", err)
		}
	}

	out, err := os.Create(fullPath) // Create file at destination path
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer out.Close() // Ensure file is closed after writing

	_, err = io.Copy(out, resp.Body) // Write response body into file
	if err != nil {
		return fmt.Errorf("error saving PDF: %w", err)
	}

	return nil // Return nil on success
}

/*
The function takes two parameters: path and permission.
We use os.Mkdir() to create the directory.
If there is an error, we use log.Println() to log the error.
*/
func createDirectory(path string, permission os.FileMode) {
	err := os.Mkdir(path, permission) // Create the directory with given permissions
	if err != nil {
		log.Println(err) // Log the error if directory creation fails
	}
}

func main() {
	const folderName = "PDFs"                // Define the name of the folder for PDFs
	createDirectory(folderName, os.ModePerm) // Create the folder with full permissions

	urlFromChemicalGuys := "https://www.chemicalguys.com/pages/material-safety-data-sheets" // URL to the SDS page

	localURLFilePath := path.Join("chemical_guys_sds_page.html") // Define the local filename for the downloaded HTML

	downloadFileUsingURLandFilePath(urlFromChemicalGuys, localURLFilePath) // Download the SDS page HTML to local file

	links, err := ExtractURLsFromHTMLFile(localURLFilePath) // Extract PDF links from HTML
	if err != nil {
		log.Printf("Error extracting URLs: %v", err) // Log error if extraction fails
	}

	for _, link := range links { // Loop through each PDF link
		err := downloadPDF(link, folderName) // Download each PDF to the folder
		if err != nil {
			log.Printf("Failed to download %s: %v", link, err) // Log error if download fails
			continue                                           // Continue to next link
		}
		log.Printf("Downloaded %s successfully", link) // Log success
	}
}
