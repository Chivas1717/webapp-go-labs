package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "strconv"
)

// ------------------ ФУНКЦІЯ РОЗРАХУНКУ ------------------
func calculateEmissions(coalStr, fuelOilStr, naturalGasStr string) string {
    coalAmount, _ := strconv.ParseFloat(coalStr, 64)
    fuelOilAmount, _ := strconv.ParseFloat(fuelOilStr, 64)
    gasAmount, _ := strconv.ParseFloat(naturalGasStr, 64)

    emissionFactorCoal := 150.0
    emissionFactorFuelOil := 0.57
    emissionFactorGas := 0.0

    heatValueCoal := 20.47
    heatValueFuelOil := 40.40
    heatValueGas := 33.08

    totalCoalEmissions := (emissionFactorCoal * heatValueCoal * coalAmount) / 1_000_000
    totalFuelOilEmissions := (emissionFactorFuelOil * heatValueFuelOil * fuelOilAmount) / 1_000_000
    totalGasEmissions := (emissionFactorGas * heatValueGas * gasAmount) / 1_000_000

    totalEmissions := totalCoalEmissions + totalFuelOilEmissions + totalGasEmissions

    result := fmt.Sprintf(`Валові викиди при спалюванні палива:
• Вугілля: %.4f т
• Мазут: %.4f т
• Природний газ: %.4f т
--------------------------------
Загальна кількість викидів: %.4f т`,
        totalCoalEmissions, totalFuelOilEmissions, totalGasEmissions, totalEmissions)

    return result
}

// ------------------ СТРУКТУРА ДЛЯ ШАБЛОНУ ------------------
type PageData struct {
    CoalInput         string
    FuelOilInput      string
    GasInput          string
    CalculationResult string
    Error             string
}

// ------------------ ЗАВАНТАЖЕННЯ ШАБЛОНУ ------------------
var tmpl *template.Template

func init() {
    // Підключаємо шаблон з файлу
    var err error
    tmpl, err = template.ParseFiles("templates/emissions.html")
    if err != nil {
        log.Fatalf("Помилка при завантаженні шаблону: %v", err)
    }
}

// ------------------ ОБРОБНИК ------------------
func handleEmissions(w http.ResponseWriter, r *http.Request) {
    data := PageData{
        // Дефолтні значення (за бажанням можна змінити)
        CoalInput:    "100",
        FuelOilInput: "50",
        GasInput:     "200",
    }

    if r.Method == http.MethodPost {
        if err := r.ParseForm(); err != nil {
            data.Error = "Помилка парсингу форми"
            tmpl.Execute(w, data)
            return
        }

        coalStr := r.FormValue("coal")
        fuelOilStr := r.FormValue("fuelOil")
        gasStr := r.FormValue("gas")

        data.CoalInput = coalStr
        data.FuelOilInput = fuelOilStr
        data.GasInput = gasStr

        data.CalculationResult = calculateEmissions(coalStr, fuelOilStr, gasStr)
    }

    tmpl.Execute(w, data)
}

// ------------------ MAIN ------------------
func main() {
    // Маршрут для калькулятора викидів
    http.HandleFunc("/emissions", handleEmissions)

    // Якщо звертаються на "/", перенаправляємо на /emissions
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, "/emissions", http.StatusFound)
    })

    fmt.Println("Сервер запущено на http://localhost:8080/emissions")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
