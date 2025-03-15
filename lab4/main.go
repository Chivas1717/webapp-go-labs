package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

var templates *template.Template

func init() {
	var err error
	templates, err = template.ParseFiles(
		"templates/index.html",
		"templates/three_phase.html",
		"templates/single_phase.html",
		"templates/stability.html",
	)
	if err != nil {
		log.Fatalf("Помилка завантаження шаблонів: %v", err)
	}
}

// ThreePhaseHandler – розрахунок струму трифазного короткого замикання
func ThreePhaseHandler(w http.ResponseWriter, r *http.Request) {
	var result string
	if r.Method == http.MethodPost {
		voltage, _ := strconv.ParseFloat(r.FormValue("voltage"), 64)
		impedance, _ := strconv.ParseFloat(r.FormValue("impedance"), 64)
		if impedance != 0 {
			// Для трифазного КЗ струм = U / (Z * √3)
			result = fmt.Sprintf("Струм трифазного КЗ: %.2f A", voltage/(impedance*math.Sqrt(3)))
		} else {
			result = "Помилка: Імпеданс не може бути нулем."
		}
	}
	templates.ExecuteTemplate(w, "three_phase.html", result)
}

// SinglePhaseHandler – розрахунок струму однофазного короткого замикання
func SinglePhaseHandler(w http.ResponseWriter, r *http.Request) {
	var result string
	if r.Method == http.MethodPost {
		voltage, _ := strconv.ParseFloat(r.FormValue("voltage"), 64)
		impedance, _ := strconv.ParseFloat(r.FormValue("impedance"), 64)
		if impedance != 0 {
			// Для однофазного КЗ струм = U / Z
			result = fmt.Sprintf("Струм однофазного КЗ: %.2f A", voltage/impedance)
		} else {
			result = "Помилка: Імпеданс не може бути нулем."
		}
	}
	templates.ExecuteTemplate(w, "single_phase.html", result)
}

// StabilityHandler – перевірка термічної стійкості (розрахунок A²·с)
func StabilityHandler(w http.ResponseWriter, r *http.Request) {
	var result string
	if r.Method == http.MethodPost {
		current, _ := strconv.ParseFloat(r.FormValue("current"), 64)
		duration, _ := strconv.ParseFloat(r.FormValue("duration"), 64)
		result = fmt.Sprintf("Термічна стійкість: %.2f A²·с", current*current*duration)
	}
	templates.ExecuteTemplate(w, "stability.html", result)
}

// IndexHandler – головне меню з посиланнями на всі калькулятори
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/threephase", ThreePhaseHandler)
	http.HandleFunc("/singlephase", SinglePhaseHandler)
	http.HandleFunc("/stability", StabilityHandler)

	fmt.Println("Сервер запущено на http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
