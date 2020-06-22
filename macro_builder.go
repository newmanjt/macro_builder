package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Gender string

const (
	Male   Gender = "Male"
	Female        = "Female"
)

func calculateBMR(gender Gender, weight float64, height float64, age float64) float64 {
	var BMR float64
	if gender == Male {
		BMR = 66 + (6.23 * weight) + (12.7 * height) - (6.8 * age)
	} else if gender == Female {
		BMR = 655 + (4.35 * weight) + (4.7 * height) - (4.7 * age)
	}
	return BMR
}

func loadPage(title string) string {
	filename := title + ".html"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return ""
	}
	return string(body)
}

func handler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, loadPage("index"))
}

func bmrHandler(w http.ResponseWriter, req *http.Request) {
	//calculating bmr
	gender, ok := req.URL.Query()["gender"]
	if !ok || len(gender[0]) < 1 {
		log.Println("Url param gender is mising")
		return
	}
	log.Println("%s", gender)
	weight, ok := req.URL.Query()["weight"]
	if !ok || len(weight[0]) < 1 {
		log.Println("Url param weight is mising")
		return
	}
	log.Println("%s", weight)
	height, ok := req.URL.Query()["height"]
	if !ok || len(height[0]) < 1 {
		log.Println("Url param height is mising")
		return
	}
	log.Println("%s", height)
	age, ok := req.URL.Query()["age"]
	if !ok || len(age[0]) < 1 {
		log.Println("Url param age is mising")
		return
	}
	log.Println("%s", age)
}

func main() {
	fmt.Println("vim-go")
	http.HandleFunc("/", handler)
	http.HandleFunc("/bmr", bmrHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
