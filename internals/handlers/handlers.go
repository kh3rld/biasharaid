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

// var data renders.FormData
var currentYear = time.Now().Format("2006")

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	data := renders.FormData{
		CurrentYear: currentYear,
		Title:       "Welcome to BiasharaID - Your Secure Blockchain Identity Verification",
	}
	renders.RenderTemplate(w, "home.page.html", &data)
}

func Verification(w http.ResponseWriter, r *http.Request) {
	data := renders.FormData{
		CurrentYear: currentYear,
		Title:       "Verify Your Identity - BiasharaID",
	}
	renders.RenderTemplate(w, "verify.page.html", &data)
}

func Contact(w http.ResponseWriter, r *http.Request) {
	data := renders.FormData{
		CurrentYear: currentYear,
		Title:       "Contact Us - BiasharaID",
	}
	renders.RenderTemplate(w, "contact.page.html", &data)
}

func About(w http.ResponseWriter, r *http.Request) {
	data := renders.FormData{
		CurrentYear: currentYear,
		Title:       "About Us - BiasharaID",
	}
	renders.RenderTemplate(w, "about.page.html", &data)
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
		data := renders.FormData{
			CurrentYear: currentYear,
			Title:       "Upload Form  - BiasharaID",
		}
		renders.RenderTemplate(w, "test.page.html", &data)

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
		data := renders.FormData{
			CurrentYear: currentYear,
			Title:       "Verify - BiasharaID",
		}
		renders.RenderTemplate(w, "verify.page.html", &data)
		return
	case "POST":
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusInternalServerError)
			return
		}

		nationalID := r.FormValue("national_id")

		if nationalID == "" {
			data := renders.FormData{
				CurrentYear: currentYear,
				Title:       "Verify - BiasharaID",
			}
			renders.RenderTemplate(w, "verify.page.html", &data)
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
			data := renders.FormData{
				CurrentYear: currentYear,
				Title:       "Page Not Found - BiasharaID",
			}
			renders.RenderTemplate(w, "404.page.html", &data)
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
		data := renders.FormData{
			CurrentYear: currentYear,
		}
		renders.RenderTemplate(w, "signup.page.html", &data)
		return
	case "POST":
		var entrepreneur blockchain.Entrepreneur
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusInternalServerError)
			return
		}
		fmt.Println(len(blockchain.BlockchainInstance.Blocks))
		first_name := r.FormValue("firstName")
		second_name := r.FormValue("secondName")
		location := r.FormValue("location")
		phone := r.FormValue("phone")
		national_id := r.FormValue("national_id")
		business_id := r.FormValue("business_id")
		status := r.FormValue("status")
		business_value := r.FormValue("businessValue")
		name := r.FormValue("businessName")
		address := r.FormValue("businessaddress")

		business := blockchain.Business{
			BusinessID:    business_id,
			Status:        status,
			BusinessValue: business_value,
			Name:          name,
			Address:       address,
		}

		// Create an instance of Entrepreneur
		entrepreneur = blockchain.Entrepreneur{
			FirstName:  first_name,
			SecondName: second_name,
			Location:   location,
			Business:   business,
			Phone:      phone,
			NationalID: national_id,
			IsGenesis:  false,
		}

		blockchain.BlockchainInstance.AddBlock(entrepreneur)
		fmt.Println(len(blockchain.BlockchainInstance.Blocks))

		data := renders.FormData{
			CurrentYear: currentYear,
			Title:       "SignUp - BiasharaID",
		}
		renders.RenderTemplate(w, "signup.page.html", &data)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func Addpage(w http.ResponseWriter, r *http.Request) {
	data := renders.FormData{
		CurrentYear: currentYear,
		Title:       "SignUp - BiasharaID",
	}
	renders.RenderTemplate(w, "signup.page.html", data)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	data := renders.FormData{
		CurrentYear: currentYear,
		Title:       "Not Found - BiasharaID",
	}
	renders.RenderTemplate(w, "404.page.html", &data)
}

func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	data := renders.FormData{
		CurrentYear: currentYear,
		Title:       "Not Found - BiasharaID",
	}
	renders.RenderTemplate(w, "400.page.html", data)
}

func ServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	data := renders.FormData{
		CurrentYear: currentYear,
		Title:       "Internal Server Error - BiasharaID",
	}
	renders.RenderTemplate(w, "500.page.html", &data)
}
