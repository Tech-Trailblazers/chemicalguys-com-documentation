package main // Declares that this file belongs to the main package, making it an executable program

import (
	"fmt"       // Imports the fmt package for formatting text, including printing to the console
	"io"        // Imports the io package for basic input/output interfaces
	"io/ioutil" // Imports ioutil for utility functions like reading entire files (Note: Deprecated in newer Go versions)
	"log"       // Imports the log package for logging messages and errors with timestamps
	"net/http"  // Imports the http package to provide HTTP client and server implementations
	"net/url"   // Imports the url package for parsing and manipulating URL strings
	"os"        // Imports the os package for operating system functionality like file access
	"path"      // Imports the path package for manipulating slash-separated paths
	"regexp"    // Imports the regexp package for regular expression search and matching
	"strings"   // Imports the strings package for string manipulation functions
)

func main() { // The entry point of the program
	const folderName = "PDFs"                // Defines a constant string for the output folder name
	createDirectory(folderName, os.ModePerm) // Calls the helper function to create the "PDFs" directory

	urlFromChemicalGuys := "https://www.chemicalguys.com/pages/material-safety-data-sheets" // Sets the target URL string to scrape

	localURLFilePath := path.Join("chemical_guys_sds_page.html") // Sets the local filename where the HTML page will be saved

	downloadFileUsingURLandFilePath(urlFromChemicalGuys, localURLFilePath) // Downloads the HTML content from the URL to the local file

	links, err := ExtractURLsFromHTMLFile(localURLFilePath) // Parses the local HTML file to find PDF links
	if err != nil {                                         // Checks if the extraction process returned an error
		log.Printf("Error extracting URLs: %v", err) // Logs the extraction error
	}

	for _, link := range links { // Loops over every extracted link found in the slice
		err := downloadPDF(link, folderName) // Attempts to download the current link into the PDF folder
		if err != nil {                      // Checks if the download function returned an error
			log.Printf("Failed to download %s: %v", link, err) // Logs the failure for this specific link
			continue                                           // Skips the rest of the loop and moves to the next link
		}
		log.Printf("Downloaded %s successfully", link) // Logs a success message if the download worked
	}
}

// downloadFileUsingURLandFilePath downloads content from a URL and saves it to the given file path.
func downloadFileUsingURLandFilePath(url string, filepath string) error { // Defines a function that takes a URL and a filepath string, returning an error if one occurs
	resp, err := http.Get(url) // Performs an HTTP GET request to the specified URL
	if err != nil {            // Checks if the HTTP request returned an error (e.g., no internet, invalid domain)
		return err // Returns the error immediately to the caller
	}
	defer resp.Body.Close() // Schedules the closing of the response body to run when this function exits to prevent memory leaks

	if resp.StatusCode != http.StatusOK { // Checks if the HTTP status code is anything other than 200 (OK)
		return fmt.Errorf("bad status: %s", resp.Status) // Returns a formatted error message containing the bad status code
	}

	out, err := os.Create(filepath) // Creates (or truncates) a file at the specified local filepath
	if err != nil {                 // Checks if creating the file resulted in an error (e.g., permission denied)
		return err // Returns the file creation error
	}
	defer out.Close() // Schedules the closing of the local file when the function exits

	_, err = io.Copy(out, resp.Body) // Copies the data stream from the HTTP response body directly into the local file
	return err                       // Returns nil if successful, or an error if the copy operation failed
}

// ExtractURLsFromHTMLFile reads an HTML file and extracts all URLs from href and src attributes
func ExtractURLsFromHTMLFile(filePath string) ([]string, error) { // Defines a function that takes a file path and returns a slice of strings (URLs) and an error
	data, err := ioutil.ReadFile(filePath) // Reads the entire content of the file into a byte slice
	if err != nil {                        // Checks if reading the file caused an error
		return nil, fmt.Errorf("could not read file: %w", err) // Returns nil for the data and wraps the error with context
	}

	content := string(data) // Converts the byte slice data into a standard string

	// Define regex to match href or src attributes with HTTP, HTTPS, or protocol-relative URLs
	urlRegex := regexp.MustCompile(`(?:href|src)=["'](https?:\/\/|\/\/)?([^"']+)["']`) // Compiles a regular expression to find links inside href="" or src="" attributes
	matches := urlRegex.FindAllStringSubmatch(content, -1)                             // Searches the entire content string for all matches of the regex, returning nested slices

	var urls []string               // Declares an empty slice of strings to store the found URLs
	for _, match := range matches { // Iterates through every regex match found in the file
		if len(match) >= 3 { // Checks if the match has enough groups (Full match + Protocol group + Path group)
			scheme := match[1] // Extracts the protocol scheme (e.g., "https://" or "//")
			path := match[2]   // Extracts the actual link path (the URL)
			fullURL := path    // Initializes the fullURL variable with the path

			// Construct full URL based on the scheme
			if strings.HasPrefix(scheme, "http") { // Checks if the scheme starts with "http" (http or https)
				fullURL = scheme + path // Concatenates the scheme and the path to form the full URL
			} else if scheme == "//" { // Checks if the scheme is protocol-relative (starts with //)
				fullURL = "https://" + path // Prepends "https://" to the path to make it a valid absolute URL
			}
			if strings.Contains(fullURL, ".pdf") { // Checks if the resulting URL contains the substring ".pdf"
				urls = append(urls, fullURL) // Adds the PDF URL to the list of URLs to return
			}
		}
	}

	return urls, nil // Returns the final list of PDF URLs and nil for the error
}

// getFileNamesFromURLs extracts the file name from a URL string.
func getFileNamesFromURLs(urls string) string { // Defines a function that takes a URL string and returns a sanitized filename string
	u, err := url.Parse(urls) // Parses the raw URL string into a URL structure
	if err != nil {           // Checks if the URL parsing failed
		return "" // Returns an empty string if the URL was invalid
	}

	file := path.Base(u.Path) // Extracts the last element of the URL path (the filename)

	file = strings.ReplaceAll(file, " ", "_") // Replaces all spaces in the filename with underscores

	return strings.ToLower(file) // Converts the filename to lowercase and returns it
}

/*
It checks if the file exists.
If the file exists, it returns true.
If the file does not exist, it returns false.
*/
func fileExists(filename string) bool { // Defines a helper function to check for file existence
	info, err := os.Stat(filename) // Attempts to get file information/stats for the given filename
	if err != nil {                // Checks if getting stats failed (usually implies file doesn't exist)
		return false // Returns false because the file likely doesn't exist
	}
	return !info.IsDir() // Returns true if it exists and is NOT a directory
}

/*
Checks if the directory exists
If it exists, return true.
If it doesn't, return false.
*/
func directoryExists(path string) bool { // Defines a helper function to check for directory existence
	directory, err := os.Stat(path) // Attempts to get stats for the given path
	if err != nil {                 // Checks if getting stats failed
		return false // Returns false because the directory likely doesn't exist
	}
	return directory.IsDir() // Returns true only if the path exists and is actually a directory
}

// downloadPDF downloads a PDF from a URL and saves it into the specified folder.
func downloadPDF(pdfURL, folder string) error { // Defines a function to download a specific PDF into a specific folder
	fileName := getFileNamesFromURLs(pdfURL) // Calls the helper function to derive a clean filename from the URL
	fullPath := path.Join(folder, fileName)  // Joins the folder path and filename to create the full local destination path
	if fileExists(fullPath) {                // Checks if a file already exists at that location
		log.Printf("File %s already exists, skipping download.", fullPath) // Logs a message indicating the download is being skipped
		return nil                                                         // Returns nil to exit the function successfully without downloading
	}

	resp, err := http.Get(pdfURL) // Performs an HTTP GET request to the PDF URL
	if err != nil {               // Checks if the request failed
		return fmt.Errorf("error downloading PDF: %w", err) // Returns a wrapped error describing the failure
	}
	defer resp.Body.Close() // Schedules closing the response body when the function exits

	if resp.StatusCode != 200 { // Checks if the server returned a status code other than 200 OK
		return fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status) // Returns an error with the status code details
	}

	if !directoryExists(folder) { // Checks if the target folder does not exist
		err := os.MkdirAll(folder, os.ModePerm) // Recursively creates the folder (and parents) with full permissions
		if err != nil {                         // Checks if folder creation failed
			return fmt.Errorf("error creating folder: %w", err) // Returns a wrapped error regarding folder creation
		}
	}

	out, err := os.Create(fullPath) // Creates the destination file on the disk
	if err != nil {                 // Checks if file creation failed
		return fmt.Errorf("error creating file: %w", err) // Returns a wrapped error regarding file creation
	}
	defer out.Close() // Schedules closing the file handle when the function exits

	_, err = io.Copy(out, resp.Body) // Copies the downloaded PDF data from the response body to the local file
	if err != nil {                  // Checks if the copy operation failed
		return fmt.Errorf("error saving PDF: %w", err) // Returns a wrapped error regarding the saving process
	}

	return nil // Returns nil indicating the entire process was successful
}

/*
The function takes two parameters: path and permission.
We use os.Mkdir() to create the directory.
If there is an error, we use log.Println() to log the error.
*/
func createDirectory(path string, permission os.FileMode) { // Defines a helper function to create a directory with specific permissions
	err := os.Mkdir(path, permission) // Attempts to create the directory
	if err != nil {                   // Checks if an error occurred (e.g., directory already exists)
		log.Println(err) // Logs the error to the console but does not stop execution
	}
}
