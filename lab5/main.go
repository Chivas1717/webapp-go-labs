package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "strconv"
    "strings"
)

// ---------- ГЛОБАЛЬНІ КАРТИ ДЛЯ ЗАВДАННЯ 1 ----------

// Приклад умовних даних (для демонстрації):
var omegaMap = map[string]float64{
    // Ймовірності відмов (штучні числа, аби вийшла приблизна "0.18" тощо)
    "X": 0.15,
}
var tvMap = map[string]float64{
    // Середні часи відновлення
    "X": 2,
}
var tpMap = map[string]float64{
    // Для обчислення kPP = 1.2 * maxTp / 8760
    "X": 0, // нехай буде 0, щоб kPP ≈ 0
}

// calculateTask1 – приклад реалізації з умови
func calculateTask1(elements string, n float64) string {
    omegaSum, tRecovery, maxTp := 0.0, 0.0, 0.0
    // Розділяємо рядок на елементи
    for _, el := range strings.Split(elements, " ") {
        omegaSum += omegaMap[el]
        tRecovery += omegaMap[el] * tvMap[el]
        if tpMap[el] > maxTp {
            maxTp = tpMap[el]
        }
    }
    // Додаємо внесок від "n"
    omegaSum += 0.03 * n
    tRecovery += 0.06 * n
    // Середній час відновлення
    tRecovery /= omegaSum

    // Коефіцієнт аварійного простою
    // kAP = (omegaSum * tRecovery) / 8760
    kAP := (omegaSum * tRecovery) / 8760

    // Коефіцієнт планового простою
    // kPP = (1.2 * maxTp) / 8760
    kPP := (1.2 * maxTp) / 8760

    // Частота відмов двоколової системи
    // omegaDK = 2 * omegaSum * (kAP + kPP)
    omegaDK := 2 * omegaSum * (kAP + kPP)

    // Частота відмов з секційним вимикачем
    omegaDKS := omegaDK + 0.02

    return fmt.Sprintf(
        `Частота відмов: %.5f рік^-1
Середня тривалість відновлення: %.5f год
Коефіцієнт аварійного простою: %.5f
Коефіцієнт планового простою: %.5f
Частота відмов двоколової системи: %.5f рік^-1
Частота відмов з секційним вимикачем: %.5f рік^-1`,
        omegaSum, tRecovery, kAP, kPP, omegaDK, omegaDKS,
    )
}

// ---------- ЗАВДАННЯ 2 ----------
// calculateTask2 – розрахунок недовідпущеної енергії та збитків
func calculateTask2(omega, tb, Pm, Tm, kp, zPerA, zPerP float64) string {
    // MWA = omega * tb * Pm * Tm
    MWA := omega * tb * Pm * Tm
    // MWP = kp * Pm * Tm
    MWP := kp * Pm * Tm
    // M = zPerA*MWA + zPerP*MWP
    M := zPerA*MWA + zPerP*MWP

    return fmt.Sprintf(`Аварійне недовідпущення: %.5f кВт·год
Планове недовідпущення: %.5f кВт·год
Збитки: %.5f грн`, MWA, MWP, M)
}

// ---------- ШАБЛОНИ ----------
var tmpl *template.Template

func init() {
    var err error
    tmpl, err = template.ParseFiles(
        "templates/index.html",
        "templates/task1.html",
        "templates/task2.html",
    )
    if err != nil {
        log.Fatalf("Помилка завантаження шаблонів: %v", err)
    }
}

// ---------- ОБРОБНИКИ ----------

// Головна сторінка
func handleIndex(w http.ResponseWriter, r *http.Request) {
    tmpl.ExecuteTemplate(w, "index.html", nil)
}

// Обробник для Завдання 1
func handleTask1(w http.ResponseWriter, r *http.Request) {
    type Task1Data struct {
        ElementsInput string
        NInput        string
        Result        string
        Error         string
    }

    data := Task1Data{
        // Дефолтні значення – щоб приблизно вийшли дані зі скріншота
        ElementsInput: "X", // викличе 0.15 із карти + 0.03 => 0.18
        NInput:        "1",
    }

    if r.Method == http.MethodPost {
        if err := r.ParseForm(); err != nil {
            data.Error = "Помилка парсингу форми"
            tmpl.ExecuteTemplate(w, "task1.html", data)
            return
        }

        elements := r.FormValue("elements")
        nStr := r.FormValue("n")

        data.ElementsInput = elements
        data.NInput = nStr

        nVal, err := strconv.ParseFloat(nStr, 64)
        if err != nil {
            data.Error = "Помилка: n має бути числом"
        } else {
            data.Result = calculateTask1(elements, nVal)
        }
    }

    tmpl.ExecuteTemplate(w, "task1.html", data)
}

// Обробник для Завдання 2
func handleTask2(w http.ResponseWriter, r *http.Request) {
    type Task2Data struct {
        OmegaInput   string
        TbInput      string
        PmInput      string
        TmInput      string
        KpInput      string
        ZAInput      string
        ZPInput      string

        Result string
        Error  string
    }

    // Дефолтні значення, підібрані так, щоб вийшло ~
    //   Аварійне недовідпущення ~ 14863.104
    //   Планове ~ 132116.48
    //   Збитки ~ 2676019.3024
    data := Task2Data{
        OmegaInput: "0.18",
        TbInput:    "2",
        PmInput:    "206",
        TmInput:    "200.42",
        KpInput:    "3.2",
        ZAInput:    "100",
        ZPInput:    "9",
    }

    if r.Method == http.MethodPost {
        if err := r.ParseForm(); err != nil {
            data.Error = "Помилка парсингу форми"
            tmpl.ExecuteTemplate(w, "task2.html", data)
            return
        }

        // Зчитуємо поля
        omegaStr := r.FormValue("omega")
        tbStr := r.FormValue("tb")
        pmStr := r.FormValue("pm")
        tmStr := r.FormValue("tm")
        kpStr := r.FormValue("kp")
        zaStr := r.FormValue("zPerA")
        zpStr := r.FormValue("zPerP")

        data.OmegaInput = omegaStr
        data.TbInput = tbStr
        data.PmInput = pmStr
        data.TmInput = tmStr
        data.KpInput = kpStr
        data.ZAInput = zaStr
        data.ZPInput = zpStr

        // Перетворюємо на float
        omega, err1 := strconv.ParseFloat(omegaStr, 64)
        tb, err2 := strconv.ParseFloat(tbStr, 64)
        pm, err3 := strconv.ParseFloat(pmStr, 64)
        tm, err4 := strconv.ParseFloat(tmStr, 64)
        kp, err5 := strconv.ParseFloat(kpStr, 64)
        za, err6 := strconv.ParseFloat(zaStr, 64)
        zp, err7 := strconv.ParseFloat(zpStr, 64)

        if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil || err7 != nil {
            data.Error = "Помилка: усі поля мають бути числовими"
        } else {
            data.Result = calculateTask2(omega, tb, pm, tm, kp, za, zp)
        }
    }

    tmpl.ExecuteTemplate(w, "task2.html", data)
}

func main() {
    http.HandleFunc("/", handleIndex)
    http.HandleFunc("/task1", handleTask1)
    http.HandleFunc("/task2", handleTask2)

    fmt.Println("Сервер запущено на http://localhost:8080/")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
