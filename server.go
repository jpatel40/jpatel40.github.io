package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const baseURL = "https://groupietrackers.herokuapp.com/api"

type MyArtistFull struct {
	ID             int                 `json:"id"`
	Image          string              `json:"image"`
	Name           string              `json:"name"`
	Members        []string            `json:"members"`
	CreationDate   int                 `json:"creationDate"`
	FirstAlbum     string              `json:"firstAlbum"`
	Locations      []string            `json:"locations"`
	ConcertDates   []string            `json:"concertDates"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type MyArtist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type MyLocation struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type MyRelation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type MyDate struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type MyDates struct {
	Index []MyDate `json:"index"`
}

type MyLocations struct {
	Index []MyLocation `json:"index"`
}

type MyRelations struct {
	Index []MyRelation `json:"index"`
}

var (
	ArtistsFull []MyArtistFull
	Artists     []MyArtist
	Dates       MyDates
	Locations   MyLocations
	Relations   MyRelations
	data        []MyArtistFull
	err100      error
	port        string
)

var tpl *template.Template

func GetArtistsData() error {
	resp, err := http.Get(baseURL + "/artists")
	if err != nil {
		return errors.New("Error by get")
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error by ReadAll")
	}
	json.Unmarshal(bytes, &Artists)
	return nil
}

func GetDatesData() error {
	resp, err := http.Get(baseURL + "/dates")
	if err != nil {
		return errors.New("Error by get")
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error by ReadAll")
	}
	json.Unmarshal(bytes, &Dates)
	return nil
}

func GetLocationsData() error {
	resp, err := http.Get(baseURL + "/locations")
	if err != nil {
		return errors.New("Error by get")
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error by ReadAll")
	}
	json.Unmarshal(bytes, &Locations)
	return nil
}

func GetRelationsData() error {
	resp, err := http.Get(baseURL + "/relation")
	if err != nil {
		return errors.New("Error by get")
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error by ReadAll")
	}
	json.Unmarshal(bytes, &Relations)
	return nil
}

func GetData() error {
	if len(ArtistsFull) != 0 {
		return nil
	}
	err1 := GetArtistsData()
	err2 := GetLocationsData()
	err3 := GetDatesData()
	err4 := GetRelationsData()
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return errors.New("Error by get data artists, locations, dates")
	}
	for i := range Artists {
		var tmpl MyArtistFull
		tmpl.ID = i + 1
		tmpl.Image = Artists[i].Image
		tmpl.Name = Artists[i].Name
		tmpl.Members = Artists[i].Members
		tmpl.CreationDate = Artists[i].CreationDate
		tmpl.FirstAlbum = Artists[i].FirstAlbum
		tmpl.Locations = Locations.Index[i].Locations
		tmpl.ConcertDates = Dates.Index[i].Dates
		tmpl.DatesLocations = Relations.Index[i].DatesLocations
		ArtistsFull = append(ArtistsFull, tmpl)
	}
	return nil
}

func GetArtistByID(id int) (MyArtist, error) {
	for _, artist := range Artists {
		if artist.ID == id {
			return artist, nil
		}
	}
	return MyArtist{}, errors.New("Not found")
}

func GetDateByID(id int) (MyDate, error) {
	for _, date := range Dates.Index {
		if date.ID == id {
			return date, nil
		}
	}
	return MyDate{}, errors.New("Not found")
}

func GetLocationByID(id int) (MyLocation, error) {
	for _, location := range Locations.Index {
		if location.ID == id {
			return location, nil
		}
	}
	return MyLocation{}, errors.New("Not found")
}

func GetRelationByID(id int) (MyRelation, error) {
	for _, relation := range Relations.Index {
		if relation.ID == id {
			return relation, nil
		}
	}
	return MyRelation{}, errors.New("Not found")
}

func GetFullDataById(id int) (MyArtistFull, error) {
	for _, artist := range ArtistsFull {
		if artist.ID == id {
			fmt.Printf("GetFullDataById|%v-%v\n", artist.ID, MyArtistFull{})
			return artist, nil
		}
	}

	return MyArtistFull{}, errors.New("Not found")
}

func ConverterStructToString() ([]string, error) {
	var data []string
	for i := 1; i <= len(Artists); i++ {
		artist, err1 := GetArtistByID(i)
		locations, err2 := GetLocationByID(i)
		date, err3 := GetDateByID(i)
		if err1 != nil || err2 != nil || err3 != nil {
			return data, errors.New("Error by converter")
		}

		str := artist.Name + " "
		for _, member := range artist.Members {
			str += member + " "
		}
		str += strconv.Itoa(artist.CreationDate) + " "
		str += artist.FirstAlbum + " "
		for _, location := range locations.Locations {
			str += location + " "
		}
		for _, d := range date.Dates {
			str += d + " "
		}

		data = append(data, str)
	}
	println("Convert to str Done!")
	return data, nil
}

////////////////////////////THIS FUNCTION ONLY TO SUPPORT FIX WHERE SEARCH FOR JUST QUEEN ARTIST////////////////////////
func ConverterStructToStringQ() ([]string, error) {
	var data []string

	for i := 1; i <= 1; i++ {
		artist, err1 := GetArtistByID(i)
		locations, err2 := GetLocationByID(i)
		date, err3 := GetDateByID(i)
		if err1 != nil || err2 != nil || err3 != nil {
			return data, errors.New("Error by converter")
		}

		str := artist.Name + " "
		for _, member := range artist.Members {
			str += member + " "
		}
		str += strconv.Itoa(artist.CreationDate) + " "
		str += artist.FirstAlbum + " "
		for _, location := range locations.Locations {
			str += location + " "
		}
		for _, d := range date.Dates {
			str += d + " "
		}
		data = append(data, str)
	}
	println("Convert to str Done!")
	return data, nil
}

////////////////////////////END OF FUNCTION  TO SUPPORT FIX WHERE SEARCH FOR JUST  QUEEN ARTIST////////////////////////

func Search(search string) []MyArtistFull {
	if search == "" {
		return ArtistsFull
	}
	//////////////////////THIS SECTION DUE TO FIX ISSUE WITH SEARCH FOR QUEEN ARTIST WHERE APPEAR TOGETHER WITH SCORPIONS (queensland problem)
	if (search == "Queen") || search == "queen" {
		art, err := ConverterStructToStringQ()
		if err != nil {
			errors.New("Error by converter")
		}
		var search_artist []MyArtistFull

		for i, artist := range art {
			lower_band := strings.ToLower(artist)
			for i_name, l_name := range []byte(lower_band) {
				lower_search := strings.ToLower(search)
				if lower_search[0] == l_name {
					lenght_name := 0
					indx := i_name
					for _, l := range []byte(lower_search) {
						if l == lower_band[indx] {
							if indx+1 == len(lower_band) {
								break
							}
							indx++
							lenght_name++
						} else {
							break
						}
					}
					if len(search) == lenght_name {
						band, _ := GetFullDataById(i + 1)
						search_artist = append(search_artist, band)
						break
					}
				}
			}

		}
		println("Search str Done!")
		return search_artist
	}
	///////////////////////////////////////////END OF FIX SEARCH SECTION///////////////////////////////////////////////////
	art, err := ConverterStructToString()
	if err != nil {
		errors.New("Error by converter")
	}
	var search_artist []MyArtistFull

	for i, artist := range art {
		lower_band := strings.ToLower(artist)
		for i_name, l_name := range []byte(lower_band) {
			lower_search := strings.ToLower(search)
			if lower_search[0] == l_name {
				lenght_name := 0
				indx := i_name
				for _, l := range []byte(lower_search) {
					if l == lower_band[indx] {
						if indx+1 == len(lower_band) {
							break
						}
						indx++
						lenght_name++
					} else {
						break
					}
				}
				if len(search) == lenght_name {
					band, _ := GetFullDataById(i + 1)
					search_artist = append(search_artist, band)
					break
				}
			}
		}

	}
	println("Search str Done!")
	return search_artist
}

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

	err := GetData()
	if err != nil {
		errors.New("Error by get data")
	}
	main := r.FormValue("main")
	search := r.FormValue("search")
	if main == "Main Page" {
		data = Search("a")
		data = ArtistsFull
	}
	if !(search == "" && len(data) != 0) {
		data = Search(search)
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
	artist, _ := GetFullDataById(id)

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

	//port = ":8080"
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
