package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// ----------------- Структури -----------------

// Тверде паливо
type FuelComposition struct {
	H float64 // Водень, %
	C float64 // Вуглець, %
	S float64 // Сірка, %
	N float64 // Азот, %
	O float64 // Кисень, %
	W float64 // Волога, %
	A float64 // Зола, %
}

// Результати обчислень для твердого палива
type FuelResults struct {
	SDry         FuelComposition // Склад сухої маси
	SCombustible FuelComposition // Склад горючої маси
	QpH          float64         // Нижча теплота згоряння (робоча маса)
	QdH          float64         // Нижча теплота згоряння (суха маса)
	QdafH        float64         // Нижча теплота згоряння (горюча маса)
}

// Мазут
type MazutComposition struct {
	C    float64 // Вуглець, %
	H    float64 // Водень, %
	O    float64 // Кисень, %
	S    float64 // Сірка, %
	V    float64 // Ванадій, %
	W    float64 // Волога, %
	A    float64 // Зола, %
	Qdaf float64 // Нижча теплота згоряння (МДж/кг)
}

// Результати обчислень для мазуту
type MazutResults struct {
	Combustible MazutComposition // Склад горючої маси
	Qp          float64          // Нижча теплота згоряння робочої маси
}

// ----------------- Функції обчислення -----------------

// Розрахунок для твердого палива
func calculateComposition(f FuelComposition) (FuelResults, error) {
	// Простенька перевірка
	if (f.H + f.C + f.S + f.N + f.O + f.W + f.A) > 100 {
		return FuelResults{}, fmt.Errorf("сумарний відсоток компонентів перевищує 100")
	}

	KRS := 100.0 / (100.0 - f.W)         // перерахунок до сухої маси
	KRG := 100.0 / (100.0 - f.W - f.A)     // перерахунок до горючої маси

	sDry := FuelComposition{ // склад сухої маси
		H: f.H * KRS,
		C: f.C * KRS,
		S: f.S * KRS,
		N: f.N * KRS,
		O: f.O * KRS,
		A: f.A * KRS,
		W: 0,
	}

	sCombustible := FuelComposition{ // склад горючої маси
		H: f.H * KRG,
		C: f.C * KRG,
		S: f.S * KRG,
		N: f.N * KRG,
		O: f.O * KRG,
		A: 0,
		W: 0,
	}

	// Формула нижчої теплоти згоряння (Дюлонга, спрощено)
	QpH := (339*f.C + 1030*f.H - 108.8*(f.O-f.S) - 25*f.W) / 1000
	QdH := QpH * KRS   // суха маса
	QdafH := QpH * KRG // горюча маса

	return FuelResults{
		SDry:         sDry,
		SCombustible: sCombustible,
		QpH:          QpH,
		QdH:          QdH,
		QdafH:        QdafH,
	}, nil
}

// Розрахунок для мазуту
func calculateMazutComposition(f MazutComposition) (MazutResults, error) {

	Kp := (100.0 - f.W - f.A) / 100.0
	mP := MazutComposition{
		C:    f.C * Kp,
		H:    f.H * Kp,
		O:    f.O * Kp,
		S:    f.S * Kp,
		V:    f.V * Kp,
		W:    f.W,
		A:    f.A,
		Qdaf: f.Qdaf,
	}

	Qp := f.Qdaf*(100.0-f.W-f.A)/100.0 - 0.025*f.W

	return MazutResults{
		Combustible: mP,
		Qp:          Qp,
	}, nil
}

// ----------------- Завантаження шаблонів -----------------

var fuelTpl *template.Template
var mazutTpl *template.Template

func init() {
	var err error
	fuelTpl, err = template.ParseFiles("templates/fuel.html")
	if err != nil {
		log.Fatalf("Не вдалося завантажити шаблон fuel.html: %v", err)
	}
	mazutTpl, err = template.ParseFiles("templates/mazut.html")
	if err != nil {
		log.Fatalf("Не вдалося завантажити шаблон mazut.html: %v", err)
	}
}

// ----------------- HTTP-обробники -----------------

// Обробник для калькулятора твердого палива
func handleFuel(w http.ResponseWriter, r *http.Request) {
	type PageData struct {
		Input   FuelComposition
		Results *FuelResults
		Error   string
	}

	// Дефолтні значення (наприклад, варіант 4)
	data := PageData{
		Input: FuelComposition{
			H: 3.4,
            C: 70.6,
            S: 2.7,
            N: 1.2,
            O: 1.9,
            W: 5.0,
            A: 15.2,
		},
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			data.Error = "Неможливо розпарсити форму"
			fuelTpl.Execute(w, data)
			return
		}

		parse := func(key string) float64 {
			val, _ := strconv.ParseFloat(r.FormValue(key), 64)
			return val
		}

		data.Input.H = parse("H")
		data.Input.C = parse("C")
		data.Input.S = parse("S")
		data.Input.N = parse("N")
		data.Input.O = parse("O")
		data.Input.W = parse("W")
		data.Input.A = parse("A")

		res, errCalc := calculateComposition(data.Input)
		if errCalc != nil {
			data.Error = errCalc.Error()
		} else {
			data.Results = &res
		}
	}

	fuelTpl.Execute(w, data)
}

// Обробник для калькулятора мазуту
func handleMazut(w http.ResponseWriter, r *http.Request) {
	type PageData struct {
		Input   MazutComposition
		Results *MazutResults
		Error   string
	}

	// Дефолтні значення
	data := PageData{
		Input: MazutComposition{
			H:    11.20,  // Водень, %
			C:    85.50,  // Вуглець, %
			S:    2.50,   // Сірка, %
			O:    0.80,   // Кисень, %
			V:    0.00,   // Ванадій (за потреби перерахувати mg/kg у %)
			W:    0.00,   // Волога, %
			A:    1.50,   // Зола, %
			Qdaf: 40.40,   // Нижча теплота згоряння горючої маси, МДж/кг
		},
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			data.Error = "Неможливо розпарсити форму"
			mazutTpl.Execute(w, data)
			return
		}

		parse := func(key string) float64 {
			val, _ := strconv.ParseFloat(r.FormValue(key), 64)
			return val
		}

		data.Input.C = parse("C")
		data.Input.H = parse("H")
		data.Input.O = parse("O")
		data.Input.S = parse("S")
		data.Input.V = parse("V")
		data.Input.W = parse("W")
		data.Input.A = parse("A")
		data.Input.Qdaf = parse("Qdaf")

		res, errCalc := calculateMazutComposition(data.Input)
		if errCalc != nil {
			data.Error = errCalc.Error()
		} else {
			data.Results = &res
		}
	}

	mazutTpl.Execute(w, data)
}

// ----------------- Main -----------------

func main() {
	http.HandleFunc("/fuel", handleFuel)
	http.HandleFunc("/mazut", handleMazut)
	// Редірект з кореня на /fuel
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/fuel", http.StatusFound)
	})

	fmt.Println("Сервер запущено на http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
