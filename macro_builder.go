package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Gender string
type ActivityLevel float64
type ProteinMacro float64
type FatMacro float64

const (
	Male        Gender        = "Male"
	Female                    = "Female"
	Rarely      ActivityLevel = 0.25
	Moderately                = 1.375
	Average                   = 1.55
	Often                     = 1.725
	Religiously               = 1.9
	Protein     ProteinMacro  = 1.3
	Fat         FatMacro      = 0.4
)

type Nutrients struct {
	Enerc_kcal float64
	Procnt     float64
	Fat        float64
	Chocdf     float64
	Fibtg      float64
}
type FoodInfo struct {
	FoodID    string
	Label     string
	Nutrients Nutrients
}
type Food struct {
	Food FoodInfo
}
type ParsedFood struct {
	Text   string
	Parsed []Food
}

func getActivity(level string) ActivityLevel {
	switch level {
	case "rarely":
		return Rarely
	case "moderately":
		return Moderately
	case "average":
		return Average
	case "often":
		return Often
	case "religiously":
		return Religiously
	default:
		//eh well let's say average
		return Average
	}
}

func calculateProtein(weight float64) float64 {
	return weight * float64(Protein)
}

func calculateFat(weight float64) float64 {
	return weight * float64(Fat)
}

func getProteinCalories(protein_grams float64) float64 {
	return protein_grams * 4
}

func getFatCalories(fat_grams float64) float64 {
	return fat_grams * 9
}

func getCarbGrams(carb_calories float64) float64 {
	return carb_calories / 4
}

func calculateCarb(total float64, protein float64, fat float64) float64 {
	return total - (getProteinCalories(protein) + getFatCalories(fat))
}

func calculateTDCR(bmr float64, activity_level ActivityLevel) float64 {
	return bmr * float64(activity_level)
}

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

func getFood(food string) ParsedFood {
	//make request to edamam food database
	base_url := "https://api.edamam.com/api/food-database/v2/parser?"
	v := url.Values{}
	v.Set("app_id", "c11655f5")
	v.Add("app_key", "e39f152d3776f490bf3831821ee639af")
	v.Add("ingr", food)
	log.Println(base_url + v.Encode())
	resp, err := http.Get(base_url + v.Encode())
	if err != nil {
		log.Println("error getting food")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var js ParsedFood
	json.Unmarshal(body, &js)
	// log.Println("Calories: ", js.Parsed[0].Food.Nutrients.Enerc_kcal)
	return js
}

func handler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, loadPage("index"))
}

func bmrHandler(w http.ResponseWriter, req *http.Request) {
	//calculating bmr
	genders, ok := req.URL.Query()["gender"]
	if !ok || len(genders[0]) < 1 {
		log.Println("Url param gender is mising")
		return
	}
	var gender Gender
	if genders[0] == "Male" {
		gender = Male
	} else if genders[0] == "Female" {
		gender = Female
	}
	log.Println(gender)
	weights, ok := req.URL.Query()["weight"]
	if !ok || len(weights[0]) < 1 {
		log.Println("Url param weight is mising")
		return
	}
	weight, err := strconv.ParseFloat(weights[0], 64)
	if err != nil {
		log.Println("error converting weight")
		return
	}
	log.Println(weight)
	heights, ok := req.URL.Query()["height"]
	if !ok || len(heights[0]) < 1 {
		log.Println("Url param height is mising")
		return
	}
	height, err := strconv.ParseFloat(heights[0], 64)
	if err != nil {
		log.Println("error converting height")
		return
	}
	log.Println(height)
	ages, ok := req.URL.Query()["age"]
	if !ok || len(ages[0]) < 1 {
		log.Println("Url param age is mising")
		return
	}
	age, err := strconv.ParseFloat(ages[0], 64)
	if err != nil {
		log.Println("error converting age")
		return
	}
	log.Println(age)
	activitys, ok := req.URL.Query()["activity"]
	if !ok || len(activitys[0]) < 1 {
		log.Println("Url param activity is missing")
		return
	}
	activity := getActivity(activitys[0])
	log.Println(activity)

	//string to construct the output
	var responsePage string

	responsePage = "<html><head><title>BMR</title></head><body><h1>BMR Results</h1><br>"

	//calculate BMR
	bmr := calculateBMR(gender, weight, height, age)
	log.Println(bmr)
	responsePage = responsePage + "<label>BMR: " + strconv.FormatFloat(bmr, 'f', -1, 64) + "</label><br>"
	//calculate total daily calorie requirement
	totalCalories := calculateTDCR(bmr, activity)
	log.Println(totalCalories)
	responsePage = responsePage + "<label>Total Daily Calories: " + strconv.FormatFloat(totalCalories, 'f', -1, 64) + "</label><br>"
	//calculate grams of protein
	protein := calculateProtein(weight)
	proteinCals := getProteinCalories(protein)
	log.Println("Protein: ", protein, proteinCals)
	responsePage = responsePage + "<label>Protein Grams: " + strconv.FormatFloat(protein, 'f', -1, 64) + "</label><br>"
	//calculate grams of fat
	fat := calculateFat(weight)
	fatCals := getFatCalories(fat)
	log.Println("Fat: ", fat, fatCals)
	responsePage = responsePage + "<label>Fat Grams: " + strconv.FormatFloat(fat, 'f', -1, 64) + "</label><br>"
	//calculate grams of carbs
	carbCalories := calculateCarb(totalCalories, protein, fat)
	carb := getCarbGrams(carbCalories)
	log.Println("Carb: ", carb, carbCalories)
	responsePage = responsePage + "<label>Carb Grams: " + strconv.FormatFloat(carb, 'f', -1, 64) + "</label><br>"

	responsePage = responsePage + "<form action='/build_macros'><input type=\"submit\" value=\"Build Macros\"></form></body></html>"
	w.Write([]byte(responsePage))
}

func buildMacroHandler(w http.ResponseWriter, req *http.Request) {
	macros, ok := req.URL.Query()["macros"]
	log.Println(macros)
	if !ok || len(macros[0]) < 1 {
		enter_output := "<html><head><title>BMR</title></head><body><h1>Sorry, Please Enter Your Macros First</h1><br>"
		enter_output = enter_output + "<form action='/build_macros?'>"
		inputs := [4]string{"protein", "fat", "carbs", "calories"}
		for input_index := range inputs {
			input := inputs[input_index]
			enter_output = enter_output + "<label for=\"" + input + "\">" + input + ":</label><br>"
			enter_output = enter_output + "<input type=\"number\" id=\"" + input + "\" name=\"" + input + "\" value=\"" + input + "\"><br>"
		}
		enter_output = enter_output + "<input type=\"submit\" value=\"Submit\">"
		enter_output = enter_output + "</form>"
		enter_output = enter_output + "</body></html>"
		w.Write([]byte(enter_output))
	} else {
		w.Write([]byte("hello world"))
	}
}

type Person struct {
	Gender          Gender
	Age             float64
	Weight          float64
	Height          float64
	BMR             float64
	TotalCalories   float64
	ProteinGrams    float64
	ProteinCalories float64
	FatGrams        float64
	FatCalories     float64
	CarbGrams       float64
	CarbCalories    float64
}

func ruler(submit chan Person) {

}

func main() {
	fmt.Println("vim-go")
	// initialize chans
	// var submit chan Person
	// set handlers
	http.HandleFunc("/", handler)
	http.HandleFunc("/bmr", bmrHandler)
	http.HandleFunc("/build_macros", buildMacroHandler)
	// serve website
	// ingr := getFood("pretzels")
	// log.Println("Calories: ", ingr.Parsed[0].Food.Nutrients.Enerc_kcal)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
