package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// performCalculation реалізує формули:
// wPred = pc * 24 * (1 - delta)
// wUnpred = pc * 24 * delta
// revenue = wPred * price
// fine = wUnpred * penalty
// profit = revenue - fine
func performCalculation(pc, delta, price, penalty float64) float64 {
	wPred := pc * 24 * (1 - delta)
	wUnpred := pc * 24 * delta
	revenue := wPred * price
	fine := wUnpred * penalty
	profit := revenue - fine
	return profit
}

// PageData містить дані форми та результат розрахунку
type PageData struct {
	PC      string
	Delta   string
	Price   string
	Penalty string

	Profit string
	Error  string
}

var tmpl *template.Template

func init() {
	var err error
	tmpl, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("Помилка завантаження шаблону: %v", err)
	}
}

func handleCalc(w http.ResponseWriter, r *http.Request) {
	// Встановлюємо дефолтні значення
	data := PageData{
		PC:      "5",
		Delta:   "0.32",
		Price:   "7",
		Penalty: "7",
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			data.Error = "Помилка читання форми"
			tmpl.Execute(w, data)
			return
		}

		pcStr := r.FormValue("pc")
		deltaStr := r.FormValue("delta")
		priceStr := r.FormValue("price")
		penaltyStr := r.FormValue("penalty")

		data.PC = pcStr
		data.Delta = deltaStr
		data.Price = priceStr
		data.Penalty = penaltyStr

		pc, err1 := strconv.ParseFloat(pcStr, 64)
		delta, err2 := strconv.ParseFloat(deltaStr, 64)
		price, err3 := strconv.ParseFloat(priceStr, 64)
		penalty, err4 := strconv.ParseFloat(penaltyStr, 64)

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			data.Error = "Некоректні вхідні дані!"
		} else {
			profit := performCalculation(pc, delta, price, penalty)
			data.Profit = fmt.Sprintf("%.2f", profit)
		}
	}

	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/calc", handleCalc)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/calc", http.StatusFound)
	})

	fmt.Println("Сервер запущено на http://localhost:8080/calc")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
