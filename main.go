package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
)

const portNumber = ":8080"

type Input struct {
	text, banner string
}

// Go application entrypoint
func main() {
	http.HandleFunc("/", Mainpage)
	http.Handle("/ascii-art", http.FileServer(http.Dir("./templates")))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))

	fmt.Println("Starting application on port", portNumber)
	err := http.ListenAndServe(portNumber, nil)
	if err != nil {
		log.Fatal("500 Internal server error", http.StatusInternalServerError) // internal server error
		return
	}
}

func Mainpage(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, err := template.ParseFiles("./templates/welcome-template.html")
	if err != nil {
		http.Error(w, "404 template NOT FOUND", http.StatusNotFound) // if doesn't find templates
		return
	}
	if r.URL.Path != "/" {
		http.Error(w, "404 address NOT FOUND", http.StatusNotFound) // wrong adress put in
		return
	}

	switch r.Method {
	case "GET":
		// send the template page and output
		parsedTemplate.Execute(w, nil)
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		name := r.FormValue("name")        // getting user's input data from the web
		InputData := r.FormValue("banner") // getting chosen banner from the web

		bnr := "banners/" + InputData + ".txt"

		font, err := ioutil.ReadFile(bnr) // reading chosen banner file from the folder
		if err != nil {
			http.Error(w, "404 banner NOT FOUND", http.StatusNotFound)
			return
		}
		dataStr := string(font)
		ascii := strings.Split(dataStr, "\n")
		newlinesplit := strings.Split(name, "\n")
		output := ""
		for _, line := range newlinesplit {
			for i := 1; i < 9; i++ {
				arrline := []rune(line)
				for j, letter := range arrline {
					if letter >= 32 && letter <= 126 { // avoiding unwanted characters
						output = output + (ascii[(letter-32)*9+rune(i)])
					} else if arrline[j] == 13 {
						output = output + string(arrline[j])
					} else {
						http.Error(w, "400 Bad request", http.StatusBadRequest)
						return

					}
				}
				output = output + "\n"
			}
		}
		error := parsedTemplate.Execute(w, output)
		if error != nil {
			log.Fatalf("400 BAD REQUEST: %s", err)
		}
	}
}
