package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	log.Println(read())
	http.HandleFunc("/create", create)
	http.HandleFunc("/complete", complete)
	http.HandleFunc("/", homePage)
	http.HandleFunc("/clear", clear)
	log.Println("Starting web server on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", List)
}

func renderTemplate(w http.ResponseWriter, page string, data any) {
	t, err := template.ParseFiles(page)
	if err != nil {
		log.Println(err)
		return
	}

	err = t.Execute(w, data)

	if err != nil {
		log.Println(err)
		return
	}
}

var List []Task

type Task struct {
	Description string
	Completed   bool
}

func create(w http.ResponseWriter, r *http.Request) {
	desc := r.URL.Query().Get("description")

	List = append(List, Task{desc, false})
	log.Println(save())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func complete(w http.ResponseWriter, r *http.Request) {
	for i, _ := range List {
		if List[i].Completed == false && r.URL.Query().Get(strconv.Itoa(i)) == "true" {
			List[i].Completed = true
		} else if List[i].Completed == true && r.URL.Query().Get(strconv.Itoa(i)) == "" {
			List[i].Completed = false
		}
	}

	//for i, _ := range r.URL.Query() {
	//	n, _ := strconv.Atoi(i)
	//	List[n].Completed = true
	//}
	log.Println(save())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func save() error {

	b := []byte{}

	for i, _ := range List {
		b = append(b, []byte(List[i].Description)...)
		b = append(b, ':')
		b = strconv.AppendBool(b, List[i].Completed)
		b = append(b, '\n')
	}
	err := os.WriteFile("tasks.txt", b, 0666)

	return err
}

func read() error {
	i, err := os.ReadFile("tasks.txt")
	s := string(i)
	lines := strings.Split(s, "\n")

	words := []string{}

	for i := 0; i < len(lines); i++ {
		words = append(words, strings.Split(lines[i], ":")...)
	}

	for i := 0; i < len(words)-1; i = i + 2 {
		x, _ := strconv.ParseBool(words[i+1])
		List = append(List, Task{Description: words[i], Completed: x})
	}
	return err
}

func clear(w http.ResponseWriter, r *http.Request) {
	clearList()
	log.Println(save())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func clearList() {
	List = nil
}
