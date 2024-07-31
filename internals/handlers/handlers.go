package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kh3rld/biasharaid/blockchain"
	"github.com/kh3rld/biasharaid/internals/renders"
)

var data renders.FormData

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	renders.RenderTemplate(w, "home.page.html", nil)
}

func Verification(w http.ResponseWriter, r *http.Request) {
	renders.RenderTemplate(w, "verify.page.html", nil)
}

func Contact(w http.ResponseWriter, r *http.Request) {
	renders.RenderTemplate(w, "contact.page.html", nil)
}

func About(w http.ResponseWriter, r *http.Request) {
	renders.RenderTemplate(w, "about.page.html", nil)
}

func Help(w http.ResponseWriter, r *http.Request) {
	data := renders.FormData{
		CurrentYear: currentYear,
		Title:       "About Us - BiasharaID",
	}
	renders.RenderTemplate(w, "help.page.html", &data)
}

func Details(w http.ResponseWriter, r *http.Request) {

}

func DummyHandler(w http.ResponseWriter, r *http.Request) {
	resp := blockchain.BlockchainInstance.Blocks
	renders.RenderTemplate(w, "dummy.page.html", resp)
}

func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Parse the form to retrieve file
		err := r.ParseMultipartForm(10 << 20) // 10 MB limit
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusInternalServerError)
			return
		}

		// Retrieve the file from the form
		file, _, err := r.FormFile("certificate")
		if err != nil {
			http.Error(w, "Failed to retrieve file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Define the temporary file path
		tempFilePath := "./static/uploads/temp_image.jpeg"
		tempFile, err := os.Create(tempFilePath)
		if err != nil {
			http.Error(w, "Failed to create file", http.StatusInternalServerError)
			return
		}
		defer tempFile.Close()

		// Copy the uploaded file content to the temporary file
		_, err = io.Copy(tempFile, file)
		if err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}

		// Call analyzeImage to process the temporary image
		analyzeImageWithOCRSpace(tempFilePath)

		// Send a success response
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("Image analyzed successfully"))
		if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func analyzeImageWithOCRSpace(imagePath string) {
	apiKey := "K84026493788957"
	url := "https://api.ocr.space/parse/image"

	file, err := os.Open(imagePath)
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		return
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	fileName := fileInfo.Name()

	// Create a new buffer to hold the multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create a form file field and copy the file content to it
	formFile, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		log.Printf("Failed to create form file field: %v", err)
		return
	}

	if _, err := io.Copy(formFile, file); err != nil {
		log.Printf("Failed to copy file content: %v", err)
		return
	}

	// Close the writer to finalize the multipart form data
	writer.Close()

	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return
	}

	req.Header.Set("apikey", apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return
	}
	ProcessImageText(string(responseBody))
}

func ProcessImageText(resp string) string {
	var res string

	for _, field := range strings.Split(resp, "\\r\\n") {
		fmt.Println(field)
	}
	return res
}

// UploadHandler handles file upload requests
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Render the upload form
		renders.RenderTemplate(w, "test.page.html", nil)

	case "POST":
		// Parse the form to retrieve file
		err := r.ParseMultipartForm(10 << 20) // 10 MB limit
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusInternalServerError)
			return
		}

		// Retrieve the file from the form
		file, _, err := r.FormFile("certificate")
		if err != nil {
			http.Error(w, "Failed to retrieve file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Define the upload directory and file path with .jpeg extension
		uploadDir := "./static/uploads"
		err = os.MkdirAll(uploadDir, os.ModePerm) // Create directory if it doesn't exist
		if err != nil {
			http.Error(w, "Failed to create directory", http.StatusInternalServerError)
			return
		}

		// Create a file in the upload directory with a .jpeg extension
		filePath := filepath.Join(uploadDir, "uploaded_nationalID.jpg")
		outFile, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Failed to create file", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		// Copy the uploaded file content to the new file
		_, err = io.Copy(outFile, file)
		if err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}

		// Call analyzeImage to process the uploaded image
		analyzeImageWithOCRSpace(filePath)

		// Send a success response
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("File uploaded and saved as uploaded_nationalID.jpeg successfully"))
		if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		renders.RenderTemplate(w, "verify.page.html", nil)
		return
	case "POST":
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusInternalServerError)
			return
		}

		nationalID := r.FormValue("national_id")

		if nationalID == "" {
			renders.RenderTemplate(w, "verify.page.html", nil)
			return
		}

		var block *blockchain.Block
		for _, b := range blockchain.BlockchainInstance.Blocks {
			if b.Data.NationalID == nationalID {
				block = b
				break
			}
		}

		if block == nil {
			renders.RenderTemplate(w, "notverified.page.html", nil)
			return
		}

		renders.RenderTemplate(w, "verify.page.html", block)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func Add(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		renders.RenderTemplate(w, "signup.page.html", nil)
		return
	case "POST":
		var entrepreneur blockchain.Entrepreneur

		// Parse the form data, including files
		if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB limit for files
			http.Error(w, "Failed to parse form", http.StatusInternalServerError)
			return
		}
		fmt.Println(len(blockchain.BlockchainInstance.Blocks))
		first_name := r.FormValue("firstName")
		second_name := r.FormValue("secondName")
		location := r.FormValue("location")
		phone := r.FormValue("phone")
		nationalID := r.FormValue("national_id")
		businessID := r.FormValue("business_id")
		status := r.FormValue("status")
		businessValueStr := r.FormValue("businessValue")
		businessName := r.FormValue("businessName")
		businessAddress := r.FormValue("businessaddress")

		// Convert business value from string to appropriate type

		// Create a Business struct
		business := blockchain.Business{
			BusinessID:    businessID,
			Status:        status,
			BusinessValue: businessValueStr,
			Name:          businessName,
			Address:       businessAddress,
		}

		// Create an Entrepreneur struct
		entrepreneur = blockchain.Entrepreneur{
			FirstName:  first_name,
			SecondName: second_name,
			Location:   location,
			Business:   business,
			Phone:      phone,
			NationalID: nationalID,
			IsGenesis:  false,
		}

		// Handle file upload
		file, _, err := r.FormFile("certificate")
		if err != nil {
			http.Error(w, "Failed to retrieve file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Process the file if needed
		// Example: Save the file to the server, check its type, etc.

		// Add the entrepreneur to the blockchain
		if blockchain.BlockchainInstance == nil {
			http.Error(w, "Blockchain instance not initialized", http.StatusInternalServerError)
			return
		}
		blockchain.BlockchainInstance.AddBlock(entrepreneur)

		// Render the template
		renders.RenderTemplate(w, "signup.page.html", nil)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func Addpage(w http.ResponseWriter, r *http.Request) {

	renders.RenderTemplate(w, "signup.page.html", data)
}
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	renders.RenderTemplate(w, "404.page.html", nil)
}

func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	renders.RenderTemplate(w, "400.page.html", nil)
}

func ServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	renders.RenderTemplate(w, "500.page.html", nil)
}
