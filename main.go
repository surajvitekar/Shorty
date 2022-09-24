package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"text/template"
)

type URLData struct {
	MainURL  string
	ShortURL string
}

var ReturnData []URLData
var totalRec int

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
)

func RandStringBytesMask(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; {
		if idx := int(rand.Int63() & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i++
		}
	}
	return string(b)
}

func renderTemplate(w http.ResponseWriter, r *http.Request) {
	visited := r.URL.String()
	fmt.Println("visited URL is", visited)
	if visited == "/" {
		parsedTemplate, _ := template.ParseFiles("./index.html")
		err := parsedTemplate.Execute(w, nil)
		if err != nil {
			log.Println("Error executing template :", err)
			return
		}
	} else if visited == "/favicon.ico" {
		log.Println("Do nothing")
	} else {
		newUrl := getMainURL(visited)
		newUrl = "https://" + newUrl
		http.Redirect(w, r, newUrl, http.StatusSeeOther)
		return
	}

}

func getMainURL(visited string) string {
	//visited = strings.Replace(visited, "/", "", -1)
	log.Println("Changed to ", visited)
	for _, i := range ReturnData {
		i.ShortURL = strings.Replace(i.ShortURL, "localhost:8090", "", -1)
		log.Println("comparing", visited, "with ", i.ShortURL)
		if visited == i.ShortURL {
			println("URL already present")
			return (i.MainURL)
		}
	}
	return "nil"
}

func getURL(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("URL : ", r.Form.Get("url"))
	urlReturned := r.Form.Get("url")
	urlReturned = strings.Replace(urlReturned, "https://", "", -1)
	urlReturned = strings.Replace(urlReturned, "http://", "", -1)
	urlReturned = strings.Replace(urlReturned, "Https://", "", -1)
	urlReturned = strings.Replace(urlReturned, "Http://", "", -1)
	ShortURL := saveAndReturnURL(urlReturned)

	w.Write([]byte(ShortURL))
}

func saveAndReturnURL(rec string) string {
	// check if the array is empty, if yes.. initialise the array
	hashedURL := ""

	fmt.Println("Checking if the URL presents already")
	for _, i := range ReturnData {
		if rec == i.MainURL {
			println("URL already present")
			return (i.ShortURL)
		}
	}
	println("URL not present")
	hashedURL = RandStringBytesMask(8)
	hashedURL = "localhost:8090/" + hashedURL
	fmt.Println(hashedURL)
	dtToAdd := URLData{MainURL: rec, ShortURL: hashedURL}
	ReturnData = append(ReturnData, dtToAdd)
	return hashedURL
}

func main() {
	http.HandleFunc("/", renderTemplate)
	http.HandleFunc("/url", getURL)

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println("in main printing global", urlReturned)
}
