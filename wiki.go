package main

import (
	"fmt"
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
func main() {
	p1 := &WikiPage{Title: "firstpage", Body: []byte("This is body!!!!!")}
	p1.save()
	p2, err := LoadPage(p1.Title)
	if err != nil {
		fmt.Println("Error loading page:", err)
		return
	}
	fmt.Println(p2)
}
