package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

const PUN_BASE = "http://www.punoftheday.com"
const RANDOM_PUN = PUN_BASE + "/cgi-bin/randompun.pl"
const SELECT_PUN = PUN_BASE + "/pun"

type Pun struct {
	Id   int    `json:"-"`
	Url  string `json:"url"`
	Text string `json:"text"`
}

type Error struct {
	Status int    `json:"-"`
	Detail string `json:"error"`
}

func main() {
	http.HandleFunc("/puns/today", TodayPun)
	http.HandleFunc("/puns/random", RandomPun)
	http.HandleFunc("/puns/", ShowPun)
	http.HandleFunc("/", NotFound)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	log.Printf("404 Not Found: %v\n", r.URL.Path)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := &Error{404, "not_found"}
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(err)
}

func WritePun(w http.ResponseWriter, r *http.Request, pun Pun) {
	log.Printf("200 OK: %v\n", r.URL.Path)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(pun)
}

func ShowPun(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/puns/"):]

	pun := getPun(SELECT_PUN + "/" + id)

	if id != strconv.Itoa(pun.Id) {
		NotFound(w, r)
	} else {
		WritePun(w, r, pun)
	}
}

func RandomPun(w http.ResponseWriter, r *http.Request) {
	pun := getPun(RANDOM_PUN)

	WritePun(w, r, pun)
}

func TodayPun(w http.ResponseWriter, r *http.Request) {
	pun := getPun(PUN_BASE)

	pun.Text = stripQuotes(pun.Text)

	WritePun(w, r, pun)
}

func stripQuotes(punText string) string {
	punText = strings.Replace(punText, "“", "", 1)
	punText = strings.Replace(punText, "”", "", 1)

	return punText
}

func getPun(url string) Pun {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("ERROR: Unable to access " + url)
	}

	b := resp.Body
	defer b.Close()

	pun := Pun{}

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return pun
		case tt == html.StartTagToken:
			t := z.Token()

			isParagraph := t.Data == "p"
			isInput := t.Data == "input"

			if isParagraph && pun.Text == "" {
				z.Next()
				t := z.Token()
				pun.Text = t.Data
			}
			if isInput {
				if getAttr("name", t) == "PunID" {
					pun.Id, _ = strconv.Atoi(getAttr("value", t))
					pun.Url = SELECT_PUN + "/" + strconv.Itoa(pun.Id)
				}
			}
		}
	}
}

func getAttr(at string, t html.Token) string {
	var val string
	for _, a := range t.Attr {
		if a.Key == at {
			val = a.Val
		}
	}

	return val
}
