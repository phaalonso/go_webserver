package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

// Template cache
var templates = template.Must(template.ParseFiles("templates/edit.html", "templates/view.html"))

// Path validation
// Operação vai sempre começar com edit, save ou edit
// após ela, deve existir o nome da página o qual é composto por letras maiusculas, minúsculas ou números
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := "data/" + p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600) // 0600 -> read and write for current user
}

func loadPage(title string) (*Page, error) {
	log.Printf("Loading the page [page/%v.txt]", title)
	filename := "data/" + title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	log.Printf("Rendering the template [%v]", tmpl)
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	log.Printf("Opening the view [%v]", title)
	p, err := loadPage(title)

	// Caso não encontrar a página, irá o redirecionar para /edit/{{title}}
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	log.Printf("Editing content from the page [%v]", title)

	p, err := loadPage(title)

	if err != nil {
		p = &Page{Title: title}
	}

	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	log.Printf("Saving in [%v] the content -> \"%v\"", title, body)
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	log.Println("Loading the handlers")
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Println("Starting the webserver")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
