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
	"time"

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
		CurrentYear: r.FormValue("currentYear"),
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
		errer := analyzeImageWithOCRSpace(tempFilePath)
		if errer == "Error analyzing image" {
			fmt.Println("Image analysis failed. Please try again.")
			http.Error(w, "Your National ID cannot be verified", http.StatusBadRequest)
		}

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

func analyzeImageWithOCRSpace(imagePath string) string {
	fmt.Println("heeeey")
	apiKey := "K84026493788957"
	url := "https://api.ocr.space/parse/image"

	file, err := os.Open(imagePath)
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		return "Error analyzing image"
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
		return "Error analyzing image"
	}

	if _, err := io.Copy(formFile, file); err != nil {
		log.Printf("Failed to copy file content: %v", err)
		return "Error analyzing image"
	}

	// Close the writer to finalize the multipart form data
	writer.Close()

	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return "Error analyzing image"
	}

	req.Header.Set("apikey", apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		return "Error analyzing image"
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return "Error analyzing image"
	}
	return ProcessImageText(string(responseBody))
}

func ProcessImageText(resp string) string {
	var res string

	for _, field := range strings.Split(resp, "\\r\\n") {
		fmt.Println(field)
	}
	return res
}

// UploadHandler handles file upload requests
func uploader(w http.ResponseWriter, r *http.Request) {
	// Parse the form to retrieve file
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusInternalServerError)
		return
	}

	// Retrieve the file from the form
	file, _, err := r.FormFile("nationalID")
	if err != nil {
		fmt.Printf("Error retrieving file: %v", err)
		log.Printf("Error retrieving file: %v", err) // Log the specific error
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	// Define the upload directory
	uploadDir := "./static/uploads"
	err = os.MkdirAll(uploadDir, os.ModePerm) // Create directory if it doesn't exist
	if err != nil {
		http.Error(w, "Failed to create directory", http.StatusInternalServerError)
		return
	}

	// Generate a unique file name to avoid overwriting
	fileName := fmt.Sprintf("uploaded_%d.jpg", time.Now().UnixNano())
	filePath := filepath.Join(uploadDir, fileName)

	// Create a file in the upload directory
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
	_, err = w.Write([]byte(fmt.Sprintf("File uploaded and saved as %s successfully", fileName)))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
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
		fmt.Println("It is here!!!")
		var entrepreneur blockchain.Entrepreneur

		// Parse the form data, including files
		if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB limit for files
			http.Error(w, "Failed to parse form", http.StatusInternalServerError)
			return
		}

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

		// Retrieve the file from the form
		file, _, err := r.FormFile("nationalID")
		if err != nil {
			fmt.Printf("Error retrieving file: %v", err)
			log.Printf("Error retrieving file: %v", err) // Log the specific error
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Define the upload directory
		uploadDir := "./static/uploads"
		err = os.MkdirAll(uploadDir, os.ModePerm) // Create directory if it doesn't exist
		if err != nil {
			http.Error(w, "Failed to create directory", http.StatusInternalServerError)
			return
		}

		// Generate a unique file name to avoid overwriting
		fileName := fmt.Sprintf("uploaded_%d.jpg", time.Now().UnixNano())
		filePath := filepath.Join(uploadDir, fileName)

		// Create a file in the upload directory
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

		// Process the file if needed
		// Example: Save the file to the server, check its type, etc.
		analyzeImageWithOCRSpace(filePath)

		// Add the entrepreneur to the blockchain
		if blockchain.BlockchainInstance == nil {
			http.Error(w, "Blockchain instance not initialized", http.StatusInternalServerError)
			return
		}
		blockchain.BlockchainInstance.AddBlock(entrepreneur)
		clientID := r.FormValue("national_id")

		// Render the template
		var block *blockchain.Block
		for _, b := range blockchain.BlockchainInstance.Blocks {
			if b.Data.NationalID == clientID {
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
