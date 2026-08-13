// weight_tracker.go — Go версия

package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type WeightEntry struct {
	ID     int     `json:"id"`
	Date   string  `json:"date"`
	Weight float64 `json:"weight"`
	Notes  string  `json:"notes"`
}

type Tracker struct {
	Entries []WeightEntry `json:"entries"`
	file    string
}

func NewTracker(file string) *Tracker {
	t := &Tracker{file: file}
	t.load()
	return t
}

func (t *Tracker) load() {
	data, err := os.ReadFile(t.file)
	if err != nil {
		t.Entries = []WeightEntry{}
		return
	}
	json.Unmarshal(data, &t.Entries)
}

func (t *Tracker) save() {
	data, _ := json.MarshalIndent(t.Entries, "", "  ")
	os.WriteFile(t.file, data, 0644)
}

func (t *Tracker) add(date string, weight float64, notes string) int {
	id := len(t.Entries) + 1
	t.Entries = append(t.Entries, WeightEntry{ID: id, Date: date, Weight: weight, Notes: notes})
	t.save()
	return id
}

func (t *Tracker) listAll() {
	if len(t.Entries) == 0 {
		fmt.Println("\u001B[33mНет записей.\u001B[0m")
		return
	}
	fmt.Printf("\u001B[36m%-4s %-12s %-12s %-30s\u001B[0m\n", "ID", "Дата", "Вес (кг)", "Заметки")
	fmt.Println(strings.Repeat("-", 60))
	for _, e := range t.Entries {
		fmt.Printf("%-4d %-12s %-12.1f %-30s\n", e.ID, e.Date, e.Weight, e.Notes)
	}
}

func (t *Tracker) getStats() map[string]interface{} {
	if len(t.Entries) == 0 {
		return nil
	}
	var weights []float64
	for _, e := range t.Entries {
		weights = append(weights, e.Weight)
	}
	sorted := t.Entries
	// Сортируем по дате
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i].Date > sorted[j].Date {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	first := sorted[0].Weight
	last := sorted[len(sorted)-1].Weight
	change := last - first
	sum := 0.0
	min := weights[0]
	max := weights[0]
	for _, w := range weights {
		sum += w
		if w < min {
			min = w
		}
		if w > max {
			max = w
		}
	}
	avg := sum / float64(len(weights))
	var bmi *float64
	if last > 0 {
		height := 1.75
		bmiVal := last / (height * height)
		bmi = &bmiVal
	}
	return map[string]interface{}{
		"count":   len(weights),
		"current": last,
		"average": avg,
		"min":     min,
		"max":     max,
		"change":  change,
		"bmi":     bmi,
	}
}

func (t *Tracker) printStats() {
	stats := t.getStats()
	if stats == nil {
		fmt.Println("\u001B[33mНет данных для статистики.\u001B[0m")
		return
	}
	fmt.Println("\u001B[36m📊 Статистика:\u001B[0m")
	fmt.Printf("  Всего записей: %d\n", stats["count"])
	fmt.Printf("  Текущий вес: %.1f кг\n", stats["current"])
	fmt.Printf("  Средний вес: %.1f кг\n", stats["average"])
	fmt.Printf("  Минимальный вес: %.1f кг\n", stats["min"])
	fmt.Printf("  Максимальный вес: %.1f кг\n", stats["max"])
	change := stats["change"].(float64)
	changeColor := "\u001B[32m"
	if change > 0 {
		changeColor = "\u001B[31m"
	} else if change == 0 {
		changeColor = "\u001B[33m"
	}
	fmt.Printf("  Изменение: %s%+.1f кг\u001B[0m\n", changeColor, change)
	if stats["bmi"] != nil {
		bmi := stats["bmi"].(*float64)
		bmiColor := "\u001B[32m"
		if *bmi >= 25 && *bmi <= 29.9 {
			bmiColor = "\u001B[33m"
		} else if *bmi >= 30 {
			bmiColor = "\u001B[31m"
		}
		fmt.Printf("  ИМТ: %s%.1f\u001B[0m\n", bmiColor, *bmi)
	}
}

func (t *Tracker) exportCSV(filename string) {
	if len(t.Entries) == 0 {
		fmt.Println("\u001B[33mНет данных для экспорта.\u001B[0m")
		return
	}
	file, _ := os.Create(filename)
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"ID", "Дата", "Вес (кг)", "Заметки"})
	for _, e := range t.Entries {
		writer.Write([]string{strconv.Itoa(e.ID), e.Date, fmt.Sprintf("%.1f", e.Weight), e.Notes})
	}
	fmt.Printf("\u001B[32m💾 Экспорт CSV: %s\u001B[0m\n", filename)
}

func (t *Tracker) exportJSON(filename string) {
	if len(t.Entries) == 0 {
		fmt.Println("\u001B[33mНет данных для экспорта.\u001B[0m")
		return
	}
	data, _ := json.MarshalIndent(t.Entries, "", "  ")
	os.WriteFile(filename, data, 0644)
	fmt.Printf("\u001B[32m💾 Экспорт JSON: %s\u001B[0m\n", filename)
}

func (t *Tracker) delete(id int) bool {
	for i, e := range t.Entries {
		if e.ID == id {
			t.Entries = append(t.Entries[:i], t.Entries[i+1:]...)
			t.save()
			return true
		}
	}
	return false
}

func (t *Tracker) edit(id int, field, value string) bool {
	for i, e := range t.Entries {
		if e.ID == id {
			switch field {
			case "date":
				t.Entries[i].Date = value
			case "weight":
				if w, err := strconv.ParseFloat(value, 64); err == nil {
					t.Entries[i].Weight = w
				} else {
					return false
				}
			case "notes":
				t.Entries[i].Notes = value
			default:
				return false
			}
			t.save()
			return true
		}
	}
	return false
}

func (t *Tracker) filterByDate(start, end string) {
	results := []WeightEntry{}
	for _, e := range t.Entries {
		if e.Date >= start && e.Date <= end {
			results = append(results, e)
		}
	}
	if len(results) == 0 {
		fmt.Println("\u001B[33mНет записей за указанный период.\u001B[0m")
		return
	}
	for _, e := range results {
		fmt.Printf("%d: %s | %.1f кг | %s\n", e.ID, e.Date, e.Weight, e.Notes)
	}
}

func main() {
	tracker := NewTracker("weight_data.json")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n\u001B[36m⚖️ Weight Tracker Pro (Go)\u001B[0m")
		fmt.Println("1. Добавить запись")
		fmt.Println("2. Показать все записи")
		fmt.Println("3. Статистика")
		fmt.Println("4. Экспорт в CSV")
		fmt.Println("5. Экспорт в JSON")
		fmt.Println("6. Удалить запись")
		fmt.Println("7. Редактировать запись")
		fmt.Println("8. Фильтр по дате")
		fmt.Println("9. Выход")
		fmt.Print("Выберите действие: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		switch choice {
		case "1":
			fmt.Print("Дата (ГГГГ-ММ-ДД, Enter для сегодня): ")
			date, _ := reader.ReadString('\n')
			date = strings.TrimSpace(date)
			if date == "" {
				date = time.Now().Format("2006-01-02")
			}
			fmt.Print("Вес (кг): ")
			weightStr, _ := reader.ReadString('\n')
			weight, _ := strconv.ParseFloat(strings.TrimSpace(weightStr), 64)
			fmt.Print("Заметки (опционально): ")
			notes, _ := reader.ReadString('\n')
			notes = strings.TrimSpace(notes)
			id := tracker.add(date, weight, notes)
			fmt.Printf("\u001B[32m✅ Запись добавлена (ID: %d)\u001B[0m\n", id)
		case "2":
			tracker.listAll()
		case "3":
			tracker.printStats()
		case "4":
			fmt.Print("Имя CSV файла (по умолч. weight_data.csv): ")
			filename, _ := reader.ReadString('\n')
			filename = strings.TrimSpace(filename)
			if filename == "" {
				filename = "weight_data.csv"
			}
			tracker.exportCSV(filename)
		case "5":
			fmt.Print("Имя JSON файла (по умолч. weight_data_export.json): ")
			filename, _ := reader.ReadString('\n')
			filename = strings.TrimSpace(filename)
			if filename == "" {
				filename = "weight_data_export.json"
			}
			tracker.exportJSON(filename)
		case "6":
			tracker.listAll()
			fmt.Print("Введите ID для удаления: ")
			idStr, _ := reader.ReadString('\n')
			id, _ := strconv.Atoi(strings.TrimSpace(idStr))
			if tracker.delete(id) {
				fmt.Println("\u001B[32m✅ Запись удалена.\u001B[0m")
			} else {
				fmt.Println("\u001B[31m❌ Запись не найдена.\u001B[0m")
			}
		case "7":
			tracker.listAll()
			fmt.Print("Введите ID для редактирования: ")
			idStr, _ := reader.ReadString('\n')
			id, _ := strconv.Atoi(strings.TrimSpace(idStr))
			fmt.Print("Какое поле редактировать (date, weight, notes): ")
			field, _ := reader.ReadString('\n')
			field = strings.TrimSpace(strings.ToLower(field))
			fmt.Print("Новое значение: ")
			value, _ := reader.ReadString('\n')
			value = strings.TrimSpace(value)
			if tracker.edit(id, field, value) {
				fmt.Println("\u001B[32m✅ Запись обновлена.\u001B[0m")
			} else {
				fmt.Println("\u001B[31m❌ Не удалось обновить.\u001B[0m")
			}
		case "8":
			fmt.Print("Начальная дата (ГГГГ-ММ-ДД): ")
			start, _ := reader.ReadString('\n')
			start = strings.TrimSpace(start)
			fmt.Print("Конечная дата (ГГГГ-ММ-ДД): ")
			end, _ := reader.ReadString('\n')
			end = strings.TrimSpace(end)
			tracker.filterByDate(start, end)
		case "9":
			fmt.Println("До свидания!")
			return
		default:
			fmt.Println("\u001B[31mНеверный выбор.\u001B[0m")
		}
	}
}
