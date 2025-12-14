package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"text/template"
)

const (
	maxRequestBytes = 11 << 20 // upto 10mb per-request upload limit
	maxParseMemory  = 50 << 20 // memory used while parsing multipart form
	maxRAMBytes     = 50 << 20 // 50 MB total RAM allowed
)

var (
	currentRAMBytes int64
	ramMutex        sync.Mutex
)

type HtmlPageData struct {
	Link  string
	Error string
}

func renderPage(w http.ResponseWriter, data HtmlPageData) {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

func main() {

	// Serving static files like CSS /static/ directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// when user hits / route index.html will be served
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderPage(w, HtmlPageData{})
	})

	// when user clicks upload button redirects to /upload
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {

		// Limiting the file send / request body to 10MB
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBytes)

		// This method seprates the values and files from the form recived
		// If the request body was truncated by MaxBytesReader,
		// parsing will fail and return an error here.
		err := r.ParseMultipartForm(maxParseMemory)
		if err != nil {
			renderPage(w, HtmlPageData{
				Error: "Encryption Failed: File is too large supports upto 10MB",
			})
			return
		}
		// grab the file from the provided form key
		file, header, err := r.FormFile("userUploadedFile")
		if err != nil {
			renderPage(w, HtmlPageData{
				Error: "Failed to read uploaded file",
			})
			return
		}
		defer file.Close()
		// read the file
		data, _ := io.ReadAll(file)

		//encrpyt the file
		encryptedFileData, err := encryptFile(data, header)
		if err != nil {
			http.Error(w, "Encryption Failed", http.StatusInternalServerError)
			return
		}

		// The length of encrypted file will be more than the size of file provided by user
		// since we add nonce and additional tags
		// The user can slowly upload files less than maxRequestBytes one after another
		// and exhaust the server RAM
		// To avoid that we keep track of current RAM used and reject requests
		// if the new file would exceed the maxRAMBytes limit
		encryptedFileSize := int64(len(encryptedFileData.EncryptedBytes))
		ramMutex.Lock()
		if encryptedFileSize+currentRAMBytes > maxRAMBytes {
			ramMutex.Unlock()
			renderPage(w, HtmlPageData{
				Error: "Server memory limit reached. Try again later.",
			})
			return
		}
		currentRAMBytes += encryptedFileSize
		ramMutex.Unlock()

		// generate a random token for the file
		token := generateToken()

		// store the encrypted file data in the map with the token as key
		StoreMutex.Lock()
		FileStore[token] = encryptedFileData
		StoreMutex.Unlock()

		// Send response to user with download link
		renderPage(w, HtmlPageData{
			Link: fmt.Sprintf("https://localhost:8080/download?token=%s", token),
		})

	})

	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {

		// Get the token from the URL query parameters
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "Missing Token", http.StatusBadRequest)
			return
		}

		// Retrieve the file data from the store
		StoreMutex.Lock()
		fileData, found := FileStore[token]
		if !found {
			StoreMutex.Unlock()
			http.Error(w, "Link is invalid or expired", http.StatusNotFound)
			return
		}
		// Remove the file data from the store after retrieval
		delete(FileStore, token)
		// free up the RAM used by the file
		ramMutex.Lock()
		currentRAMBytes -= int64(len(fileData.EncryptedBytes))
		ramMutex.Unlock()
		StoreMutex.Unlock()

		// Decrypt the file
		decryptedFile, err := decryptFile(fileData.EncryptedBytes, fileData.Key)
		if err != nil {
			http.Error(w, "Decryption Failed", http.StatusInternalServerError)
			return
		}

		//Serve the decrypted fille as a download
		w.Header().Set("Content-Disposition", "attachment; filename="+fileData.Filename)
		w.Header().Set("Content-type", "application/octet-stream")
		w.Write(decryptedFile)

	})

	// Just to check is server is up and running correctly (dev log)
	fmt.Print("Server is running and serving index.html\n")

	// Start the HTTPS server
	err := http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil)
	if err != nil {
		fmt.Println("Server failed to start: ", err)
	}
}
