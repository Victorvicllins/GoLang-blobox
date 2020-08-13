package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from blogBox"))
}

func showBlog(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	// w.Write([]byte("Display a blog instance..."))
	fmt.Fprintf(w, "Display a snippet with ID %d..", id)
}

func createBlog(w http.ResponseWriter, r *http.Request) {
	// Checking http method could either be  http.MethodPost OR "POST"
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		//w.WriteHeader(405)
		//w.Write([]byte("Method Not Allowed"))
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	w.Write([]byte("Create a new blog..."))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/blog", showBlog)
	mux.HandleFunc("/blog/create", createBlog)

	log.Println("Server started on :5000")
	err := http.ListenAndServe(":5000", mux)
	log.Fatal(err)
}
