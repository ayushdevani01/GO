package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type WikiPage struct {
	Title string
	Body  []byte
}

func (p *WikiPage) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func LoadPage(title string) (*WikiPage, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}
	return &WikiPage{Title: title, Body: body}, err
}

func (p *WikiPage) String() string {
	return fmt.Sprintf("Title: %s\nBody: %s", p.Title, string(p.Body))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := LoadPage(title)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "<h1>%s</h1> %s", p.Title, string(p.Body))
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := LoadPage(title)
	if err != nil {
		p = &WikiPage{Title: title}
	}
	fmt.Fprintf(w, "<h1> Editing %s </h1>"+"<form action = \"/save/%s\" method=\"post\"><textarea name=\"body\">%s</textarea><br> <input type=\"submit\" value=\"save\"> </form>", p.Title, p.Title, p.Body)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &WikiPage{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
