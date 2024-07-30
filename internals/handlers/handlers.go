package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/kh3rld/biasharaid/blockchain"
	"github.com/kh3rld/biasharaid/internals/renders"
	"google.golang.org/api/option"
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
		analyzeImage(tempFilePath)

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

func analyzeImage(imagePath string) {
	ctx := context.Background()
	client, err := vision.NewImageAnnotatorClient(ctx, option.WithCredentialsFile("path/to/your-service-account-key.json"))
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}
	defer client.Close()

	file, err := os.Open(imagePath)
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		return
	}
	defer file.Close()

	image, err := vision.NewImageFromReader(file)
	if err != nil {
		log.Printf("Failed to create image: %v", err)
		return
	}

	response, err := client.DetectLabels(ctx, image, nil, 10)
	if err != nil {
		log.Printf("Failed to detect labels: %v", err)
		return
	}

	fmt.Println("Labels:")
	for _, label := range response {
		fmt.Printf("%s (confidence: %f)\n", label.Description, label.Score)
	}
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
		filePath := filepath.Join(uploadDir, "uploaded_nationalID.jpeg")
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
		analyzeImage(filePath)

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
			renders.RenderTemplate(w, "404.page.html", nil)
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
