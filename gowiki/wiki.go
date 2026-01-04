package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Password string `json:"password"`
}

type ReginsterResponse struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Scuccess bool   `json:"success"`
	Token    string `json:"token"`
}

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

func debugHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	fmt.Fprintf(w, "Method: %s\n", r.Method)
	fmt.Fprintf(w, "URL: %s\n", r.URL.String())
	fmt.Fprintf(w, "Path: %s\n", r.URL.Path)

	fmt.Fprintf(w, "Headers:\n")
	for key, values := range r.Header {
		fmt.Fprintf(w, "%s: %v\n", key, values)
	}

	fmt.Fprintf(w, "Body : \n")
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error reading body: %v\n", err)
		return
	}
	fmt.Fprintf(w, "%s\n", bodyBytes)

}

func handleHelth(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Status : OK")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	term := queryParams.Get("q")

	if term == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	tags := queryParams["tag"]
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Search Term: %s\n", term)
	fmt.Fprintf(w, "Tags: %v\n", tags)

}

func registerUserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if newUser.Username == "" || newUser.Email == "" || newUser.Password == "" {
		http.Error(w, "Username, Email and Password are required", http.StatusBadRequest)
		return
	}

	response := ReginsterResponse{
		Id:      12345,
		Message: fmt.Sprintf("User %s registered successfully!", newUser.Username),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(response)

}

func loginUserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginrequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginrequest)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
	}

	response := LoginResponse{
		Token:    "abcdef",
		Scuccess: true,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)

}

func main() {

	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/debug/", debugHandler)
	http.HandleFunc("/health/", handleHelth)
	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/register/", registerUserHandler)
	http.HandleFunc("/login/", loginUserHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
