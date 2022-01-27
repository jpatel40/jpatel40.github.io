package gt

import (
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
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
