package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	gt "./data"
)

const baseURL = "https://groupietrackers.herokuapp.com/api"

var (
	err100 error
	data   []gt.MyArtistFull
	port   string
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		aaa, err := template.ParseFiles("404.html")
		if err != nil {
			http.Error(w, err.Error(), 400)
			http.Error(w, "Resources NotFound-400", 400)
			return
		}
		if err := aaa.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		return
	}

	_, err_url := http.Get(baseURL)

	if err_url != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["API_BadRequest 404"] = baseURL
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
	}

	err := gt.GetData()
	if err != nil {
		errors.New("Error by get data")
	}
	main := r.FormValue("main")
	search := r.FormValue("search")
	if main == "Main Page" {
		data = gt.Search("a")
		data = gt.ArtistsFull
	}
	if !(search == "" && len(data) != 0) {
		data = gt.Search(search)
	}

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		http.Error(w, "Resources NotFound-400", 400)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func concertPage(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("concert")
	id, _ := strconv.Atoi(idStr)
	artist, _ := gt.GetFullDataById(id)

	for key, value := range artist.DatesLocations {
		fmt.Print(key + "  - ")
		for _, e := range value {
			println(e)
		}
	}

	tmpl, err := template.ParseFiles("concert.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		http.Error(w, "Resources NotFound-400", 400)
		return
	}
	if err := tmpl.Execute(w, artist); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func ServerStatus(w http.ResponseWriter, r *http.Request) { // validate by client on http://localhost:8080/status  - 200 OK / 500 Error
	_, err100 = http.Get("http://127.0.0.1" + port + "/")
	if err100 == nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["StatusOK"] = "200"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["StatusInternalServerError"] = "500"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
	}
}

func main() {
	// port = ":8080"
	port = os.Getenv("PORT")

	if port == "" {
		port = ":8080"
	} else {
		port = ":" + port
	}

	http.HandleFunc("/", mainPage)
	// static folder
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/concert", concertPage)
	http.HandleFunc("/status", ServerStatus)

	println("Server listen on port:", port)
	count := 0
	for count < 2 {
		if err100 != nil {
			fmt.Println("StatusInternalServerError", http.StatusInternalServerError)
			break
		} else {
			fmt.Println("StatusOK", http.StatusOK)
		}
		err100 := http.ListenAndServe(port, nil)
		if err100 != nil {
			log.Fatal("Listen and Server", err100)
		}
		count++
	}
}
