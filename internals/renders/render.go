package renders

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type FormData struct {
	Body string
}

// Data is a global variable to hold the form data
var Data FormData

// functions is a map of template functions
var functions = template.FuncMap{}

// RenderTemplate is a helper function to render HTML templates
func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, _ := getTemplateCache()
	ts, ok := t[tmpl]
	if !ok {
		renderServerErrorTemplate(w, tmpl+" is missing, contact the Network Admin.")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := ts.Execute(w, data)
	if err != nil {
		return
	}
}

// getTemplateCache is a helper function to cache all HTML templates as a map
func getTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	baseDir := GetProjectRoot("views", "templates")

	templatesDir := filepath.Join(baseDir, "*.page.html")
	pages, err := filepath.Glob(templatesDir)
	if err != nil {
		return myCache, fmt.Errorf("error globbing templates: %v", err)
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, fmt.Errorf("error parsing page %s: %v", name, err)
		}

		layoutsPath := filepath.Join(baseDir, "*.layout.html")
		matches, err := filepath.Glob(layoutsPath)
		if err != nil {
			return myCache, fmt.Errorf("error finding layout files: %v", err)
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(layoutsPath)
			if err != nil {
				return myCache, fmt.Errorf("error parsing layout files: %v", err)
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}

// GetProjectRoot dynamically finds the project root directory
func GetProjectRoot(first, second string) string {
	cwd, _ := os.Getwd()
	baseDir := cwd
	if strings.HasSuffix(baseDir, "web") {
		baseDir = filepath.Join(cwd, "../")
	}
	return filepath.Join(baseDir, first, second)
}

// renderServerErrorTemplate renders a simple error template directly
func renderServerErrorTemplate(w http.ResponseWriter, errMsg string) {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Server Error</title>
	<style>
		body {color: bisque; background-color: #333; font-family: Arial, sans-serif; }
		.container { text-align: center; margin-top: 50px; }
		.btn { display: inline-block; margin-top: 20px; padding: 10px 20px; background-color: #007bff; color: white; text-decoration: none; border-radius: 5px; }
	</style>
</head>
<body>
	<div class="container">
		<h1>500 Oops We can't find what you are looking for! 🙁</h1>
		<h2>Something went wrong.</h2>
		<h3>{{.Error}}</h3>
		<a href="/" title="Go back to the home page" class="btn">
			<h1>Home</h1>
		</a>
	</div>
</body>
</html>`

	t, err := template.New("error").Parse(tmpl)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
	}

	data := struct {
		Error string
	}{
		Error: errMsg,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	t.Execute(w, data)
}
