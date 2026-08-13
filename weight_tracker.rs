// weight_tracker.rs — Rust версия

use serde::{Deserialize, Serialize};
use std::fs;
use std::io::{self, Write};

#[derive(Serialize, Deserialize, Clone)]
struct WeightEntry {
    id: usize,
    date: String,
    weight: f64,
    notes: String,
}

struct Tracker {
    entries: Vec<WeightEntry>,
    file: String,
}

impl Tracker {
    fn new(file: &str) -> Self {
        let mut t = Tracker { entries: Vec::new(), file: file.to_string() };
        t.load();
        t
    }

    fn load(&mut self) {
        if let Ok(data) = fs::read_to_string(&self.file) {
            if let Ok(entries) = serde_json::from_str(&data) {
                self.entries = entries;
                return;
            }
        }
        self.entries = Vec::new();
    }

    fn save(&self) {
        let data = serde_json::to_string_pretty(&self.entries).unwrap();
        fs::write(&self.file, data).unwrap();
    }

    fn add(&mut self, date: String, weight: f64, notes: String) -> usize {
        let id = self.entries.len() + 1;
        self.entries.push(WeightEntry { id, date, weight, notes });
        self.save();
        id
    }

    fn list_all(&self) {
        if self.entries.is_empty() {
            println!("\x1b[33mНет записей.\x1b[0m");
            return;
        }
        println!("\x1b[36m{:<4} {:<12} {:<12} {:<30}\x1b[0m", "ID", "Дата", "Вес (кг)", "Заметки");
        println!("{}", "-".repeat(60));
        for e in &self.entries {
            println!("{:<4} {:<12} {:<12.1} {:<30}", e.id, e.date, e.weight, e.notes);
        }
    }

    fn get_stats(&self) -> Option<Stats> {
        if self.entries.is_empty() {
            return None;
        }
        let weights: Vec<f64> = self.entries.iter().map(|e| e.weight).collect();
        let mut sorted = self.entries.clone();
        sorted.sort_by(|a, b| a.date.cmp(&b.date));
        let first = sorted.first().unwrap().weight;
        let last = sorted.last().unwrap().weight;
        let change = last - first;
        let sum: f64 = weights.iter().sum();
        let avg = sum / weights.len() as f64;
        let min = *weights.iter().min_by(|a, b| a.partial_cmp(b).unwrap()).unwrap();
        let max = *weights.iter().max_by(|a, b| a.partial_cmp(b).unwrap()).unwrap();
        let bmi = if last > 0.0 {
            let height = 1.75;
            Some(last / (height * height))
        } else {
            None
        };
        Some(Stats { count: weights.len(), current: last, average: avg, min, max, change, bmi })
    }

    fn print_stats(&self) {
        if let Some(stats) = self.get_stats() {
            println!("\x1b[36m📊 Статистика:\x1b[0m");
            println!("  Всего записей: {}", stats.count);
            println!("  Текущий вес: {:.1} кг", stats.current);
            println!("  Средний вес: {:.1} кг", stats.average);
            println!("  Минимальный вес: {:.1} кг", stats.min);
            println!("  Максимальный вес: {:.1} кг", stats.max);
            let change_color = if stats.change < 0.0 { "\x1b[32m" } else if stats.change > 0.0 { "\x1b[31m" } else { "\x1b[33m" };
            println!("  Изменение: {}{:+.1} кг\x1b[0m", change_color, stats.change);
            if let Some(bmi) = stats.bmi {
                let bmi_color = if bmi >= 18.5 && bmi <= 24.9 { "\x1b[32m" } else if bmi >= 25.0 && bmi <= 29.9 { "\x1b[33m" } else { "\x1b[31m" };
                println!("  ИМТ: {}{:.1}\x1b[0m", bmi_color, bmi);
            }
        } else {
            println!("\x1b[33mНет данных для статистики.\x1b[0m");
        }
    }

    fn export_csv(&self, filename: &str) {
        if self.entries.is_empty() {
            println!("\x1b[33mНет данных для экспорта.\x1b[0m");
            return;
        }
        let mut csv = String::from("ID,Дата,Вес (кг),Заметки\n");
        for e in &self.entries {
            csv.push_str(&format!("{},{},{:.1},{}\n", e.id, e.date, e.weight, e.notes));
        }
        fs::write(filename, csv).unwrap();
        println!("\x1b[32m💾 Экспорт CSV: {}\x1b[0m", filename);
    }

    fn export_json(&self, filename: &str) {
        if self.entries.is_empty() {
            println!("\x1b[33mНет данных для экспорта.\x1b[0m");
            return;
        }
        let data = serde_json::to_string_pretty(&self.entries).unwrap();
        fs::write(filename, data).unwrap();
        println!("\x1b[32m💾 Экспорт JSON: {}\x1b[0m", filename);
    }

    fn delete(&mut self, id: usize) -> bool {
        let pos = self.entries.iter().position(|e| e.id == id);
        if let Some(idx) = pos {
            self.entries.remove(idx);
            self.save();
            true
        } else {
            false
        }
    }

    fn edit(&mut self, id: usize, field: &str, value: &str) -> bool {
        for e in &mut self.entries {
            if e.id == id {
                match field {
                    "date" => e.date = value.to_string(),
                    "weight" => {
                        if let Ok(w) = value.parse() {
                            e.weight = w;
                        } else {
                            return false;
                        }
                    }
                    "notes" => e.notes = value.to_string(),
                    _ => return false,
                }
                self.save();
                return true;
            }
        }
        false
    }

    fn filter_by_date(&self, start: &str, end: &str) {
        let results: Vec<&WeightEntry> = self.entries.iter()
            .filter(|e| e.date >= start && e.date <= end)
            .collect();
        if results.is_empty() {
            println!("\x1b[33mНет записей за указанный период.\x1b[0m");
            return;
        }
        for e in results {
            println!("{}: {} | {:.1} кг | {}", e.id, e.date, e.weight, e.notes);
        }
    }
}

struct Stats {
    count: usize,
    current: f64,
    average: f64,
    min: f64,
    max: f64,
    change: f64,
    bmi: Option<f64>,
}

fn main() {
    let mut tracker = Tracker::new("weight_data.json");
    loop {
        println!("\n\x1b[36m⚖️ Weight Tracker Pro (Rust)\x1b[0m");
        println!("1. Добавить запись");
        println!("2. Показать все записи");
        println!("3. Статистика");
        println!("4. Экспорт в CSV");
        println!("5. Экспорт в JSON");
        println!("6. Удалить запись");
        println!("7. Редактировать запись");
        println!("8. Фильтр по дате");
        println!("9. Выход");
        print!("Выберите действие: ");
        io::stdout().flush().unwrap();
        let mut choice = String::new();
        io::stdin().read_line(&mut choice).unwrap();
        match choice.trim() {
            "1" => {
                print!("Дата (ГГГГ-ММ-ДД, Enter для сегодня): ");
                io::stdout().flush().unwrap();
                let mut date = String::new();
                io::stdin().read_line(&mut date).unwrap();
                let date = if date.trim().is_empty() {
                    chrono::Local::now().format("%Y-%m-%d").to_string()
                } else {
                    date.trim().to_string()
                };
                print!("Вес (кг): ");
                io::stdout().flush().unwrap();
                let mut weight_str = String::new();
                io::stdin().read_line(&mut weight_str).unwrap();
                let weight: f64 = weight_str.trim().parse().unwrap();
                print!("Заметки (опционально): ");
                io::stdout().flush().unwrap();
                let mut notes = String::new();
                io::stdin().read_line(&mut notes).unwrap();
                let notes = notes.trim().to_string();
                let id = tracker.add(date, weight, notes);
                println!("\x1b[32m✅ Запись добавлена (ID: {})\x1b[0m", id);
            }
            "2" => tracker.list_all(),
            "3" => tracker.print_stats(),
            "4" => {
                print!("Имя CSV файла (по умолч. weight_data.csv): ");
                io::stdout().flush().unwrap();
                let mut filename = String::new();
                io::stdin().read_line(&mut filename).unwrap();
                let filename = if filename.trim().is_empty() { "weight_data.csv".to_string() } else { filename.trim().to_string() };
                tracker.export_csv(&filename);
            }
            "5" => {
                print!("Имя JSON файла (по умолч. weight_data_export.json): ");
                io::stdout().flush().unwrap();
                let mut filename = String::new();
                io::stdin().read_line(&mut filename).unwrap();
                let filename = if filename.trim().is_empty() { "weight_data_export.json".to_string() } else { filename.trim().to_string() };
                tracker.export_json(&filename);
            }
            "6" => {
                tracker.list_all();
                print!("Введите ID для удаления: ");
                io::stdout().flush().unwrap();
                let mut id_str = String::new();
                io::stdin().read_line(&mut id_str).unwrap();
                let id: usize = id_str.trim().parse().unwrap();
                if tracker.delete(id) {
                    println!("\x1b[32m✅ Запись удалена.\x1b[0m");
                } else {
                    println!("\x1b[31m❌ Запись не найдена.\x1b[0m");
                }
            }
            "7" => {
                tracker.list_all();
                print!("Введите ID для редактирования: ");
                io::stdout().flush().unwrap();
                let mut id_str = String::new();
                io::stdin().read_line(&mut id_str).unwrap();
                let id: usize = id_str.trim().parse().unwrap();
                print!("Какое поле редактировать (date, weight, notes): ");
                io::stdout().flush().unwrap();
                let mut field = String::new();
                io::stdin().read_line(&mut field).unwrap();
                let field = field.trim().to_lowercase();
                print!("Новое значение: ");
                io::stdout().flush().unwrap();
                let mut value = String::new();
                io::stdin().read_line(&mut value).unwrap();
                if tracker.edit(id, &field, value.trim()) {
                    println!("\x1b[32m✅ Запись обновлена.\x1b[0m");
                } else {
                    println!("\x1b[31m❌ Не удалось обновить.\x1b[0m");
                }
            }
            "8" => {
                print!("Начальная дата (ГГГГ-ММ-ДД): ");
                io::stdout().flush().unwrap();
                let mut start = String::new();
                io::stdin().read_line(&mut start).unwrap();
                print!("Конечная дата (ГГГГ-ММ-ДД): ");
                io::stdout().flush().unwrap();
                let mut end = String::new();
                io::stdin().read_line(&mut end).unwrap();
                tracker.filter_by_date(start.trim(), end.trim());
            }
            "9" => {
                println!("До свидания!");
                break;
            }
            _ => println!("\x1b[31mНеверный выбор.\x1b[0m"),
        }
    }
}
