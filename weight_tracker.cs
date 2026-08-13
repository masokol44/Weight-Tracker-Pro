// weight_tracker.cs — C# версия

using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text.Json;

class WeightEntry {
    public int Id { get; set; }
    public string Date { get; set; }
    public double Weight { get; set; }
    public string Notes { get; set; }
}

class WeightTracker {
    private List<WeightEntry> entries = new List<WeightEntry>();
    private const string DataFile = "weight_data.json";

    public WeightTracker() {
        Load();
    }

    private void Load() {
        if (File.Exists(DataFile)) {
            try {
                string json = File.ReadAllText(DataFile);
                entries = JsonSerializer.Deserialize<List<WeightEntry>>(json) ?? new List<WeightEntry>();
            } catch {
                entries = new List<WeightEntry>();
            }
        }
    }

    private void Save() {
        string json = JsonSerializer.Serialize(entries, new JsonSerializerOptions { WriteIndented = true });
        File.WriteAllText(DataFile, json);
    }

    public int Add(string date, double weight, string notes = "") {
        int id = entries.Count + 1;
        entries.Add(new WeightEntry { Id = id, Date = date, Weight = weight, Notes = notes });
        Save();
        return id;
    }

    public void ListAll() {
        if (entries.Count == 0) {
            Console.WriteLine("\u001B[33mНет записей.\u001B[0m");
            return;
        }
        Console.WriteLine($"\u001B[36m{"ID",-4} {"Дата",-12} {"Вес (кг)",-12} {"Заметки",-30}\u001B[0m");
        Console.WriteLine(new string('-', 60));
        foreach (var e in entries) {
            Console.WriteLine($"{e.Id,-4} {e.Date,-12} {e.Weight,-12:F1} {e.Notes,-30}");
        }
    }

    public void PrintStats() {
        if (entries.Count == 0) {
            Console.WriteLine("\u001B[33mНет данных для статистики.\u001B[0m");
            return;
        }
        var sorted = entries.OrderBy(e => e.Date).ToList();
        double first = sorted.First().Weight;
        double last = sorted.Last().Weight;
        double change = last - first;
        double avg = entries.Average(e => e.Weight);
        double min = entries.Min(e => e.Weight);
        double max = entries.Max(e => e.Weight);
        double? bmi = null;
        if (last > 0) {
            double height = 1.75;
            bmi = last / (height * height);
        }
        Console.WriteLine("\u001B[36m📊 Статистика:\u001B[0m");
        Console.WriteLine($"  Всего записей: {entries.Count}");
        Console.WriteLine($"  Текущий вес: {last:F1} кг");
        Console.WriteLine($"  Средний вес: {avg:F1} кг");
        Console.WriteLine($"  Минимальный вес: {min:F1} кг");
        Console.WriteLine($"  Максимальный вес: {max:F1} кг");
        string changeColor = change < 0 ? "\u001B[32m" : change > 0 ? "\u001B[31m" : "\u001B[33m";
        Console.WriteLine($"  Изменение: {changeColor}{change:+F1} кг\u001B[0m");
        if (bmi.HasValue) {
            string bmiColor = bmi >= 18.5 && bmi <= 24.9 ? "\u001B[32m" : bmi >= 25 && bmi <= 29.9 ? "\u001B[33m" : "\u001B[31m";
            Console.WriteLine($"  ИМТ: {bmiColor}{bmi:F1}\u001B[0m");
        }
    }

    public void ExportCSV(string filename = "weight_data.csv") {
        if (entries.Count == 0) {
            Console.WriteLine("\u001B[33mНет данных для экспорта.\u001B[0m");
            return;
        }
        using var writer = new StreamWriter(filename);
        writer.WriteLine("ID,Дата,Вес (кг),Заметки");
        foreach (var e in entries) {
            writer.WriteLine($"{e.Id},{e.Date},{e.Weight:F1},{e.Notes}");
        }
        Console.WriteLine($"\u001B[32m💾 Экспорт CSV: {filename}\u001B[0m");
    }

    public void ExportJSON(string filename = "weight_data_export.json") {
        if (entries.Count == 0) {
            Console.WriteLine("\u001B[33mНет данных для экспорта.\u001B[0m");
            return;
        }
        string json = JsonSerializer.Serialize(entries, new JsonSerializerOptions { WriteIndented = true });
        File.WriteAllText(filename, json);
        Console.WriteLine($"\u001B[32m💾 Экспорт JSON: {filename}\u001B[0m");
    }

    public bool Delete(int id) {
        var entry = entries.FirstOrDefault(e => e.Id == id);
        if (entry != null) {
            entries.Remove(entry);
            Save();
            return true;
        }
        return false;
    }

    public bool Edit(int id, string field, string value) {
        var entry = entries.FirstOrDefault(e => e.Id == id);
        if (entry == null) return false;
        switch (field.ToLower()) {
            case "date": entry.Date = value; break;
            case "weight": if (double.TryParse(value, out double w)) entry.Weight = w; else return false; break;
            case "notes": entry.Notes = value; break;
            default: return false;
        }
        Save();
        return true;
    }

    public void FilterByDate(string start, string end) {
        var results = entries.Where(e => e.Date.CompareTo(start) >= 0 && e.Date.CompareTo(end) <= 0).ToList();
        if (results.Count == 0) {
            Console.WriteLine("\u001B[33mНет записей за указанный период.\u001B[0m");
            return;
        }
        foreach (var e in results) {
            Console.WriteLine($"{e.Id}: {e.Date} | {e.Weight:F1} кг | {e.Notes}");
        }
    }

    public static void Main() {
        var tracker = new WeightTracker();
        while (true) {
            Console.WriteLine("\n\u001B[36m⚖️ Weight Tracker Pro (C#)\u001B[0m");
            Console.WriteLine("1. Добавить запись");
            Console.WriteLine("2. Показать все записи");
            Console.WriteLine("3. Статистика");
            Console.WriteLine("4. Экспорт в CSV");
            Console.WriteLine("5. Экспорт в JSON");
            Console.WriteLine("6. Удалить запись");
            Console.WriteLine("7. Редактировать запись");
            Console.WriteLine("8. Фильтр по дате");
            Console.WriteLine("9. Выход");
            Console.Write("Выберите действие: ");
            string choice = Console.ReadLine();
            switch (choice) {
                case "1":
                    Console.Write("Дата (ГГГГ-ММ-ДД, Enter для сегодня): ");
                    string date = Console.ReadLine();
                    if (string.IsNullOrEmpty(date)) date = DateTime.Now.ToString("yyyy-MM-dd");
                    Console.Write("Вес (кг): ");
                    double weight = double.Parse(Console.ReadLine());
                    Console.Write("Заметки (опционально): ");
                    string notes = Console.ReadLine();
                    int id = tracker.Add(date, weight, notes);
                    Console.WriteLine($"\u001B[32m✅ Запись добавлена (ID: {id})\u001B[0m");
                    break;
                case "2": tracker.ListAll(); break;
                case "3": tracker.PrintStats(); break;
                case "4": tracker.ExportCSV(); break;
                case "5": tracker.ExportJSON(); break;
                case "6":
                    tracker.ListAll();
                    Console.Write("Введите ID для удаления: ");
                    int delId = int.Parse(Console.ReadLine());
                    if (tracker.Delete(delId)) {
                        Console.WriteLine("\u001B[32m✅ Запись удалена.\u001B[0m");
                    } else {
                        Console.WriteLine("\u001B[31m❌ Запись не найдена.\u001B[0m");
                    }
                    break;
                case "7":
                    tracker.ListAll();
                    Console.Write("Введите ID для редактирования: ");
                    int editId = int.Parse(Console.ReadLine());
                    Console.Write("Какое поле редактировать (date, weight, notes): ");
                    string field = Console.ReadLine().ToLower();
                    Console.Write("Новое значение: ");
                    string value = Console.ReadLine();
                    if (tracker.Edit(editId, field, value)) {
                        Console.WriteLine("\u001B[32m✅ Запись обновлена.\u001B[0m");
                    } else {
                        Console.WriteLine("\u001B[31m❌ Не удалось обновить.\u001B[0m");
                    }
                    break;
                case "8":
                    Console.Write("Начальная дата (ГГГГ-ММ-ДД): ");
                    string start = Console.ReadLine();
                    Console.Write("Конечная дата (ГГГГ-ММ-ДД): ");
                    string end = Console.ReadLine();
                    tracker.FilterByDate(start, end);
                    break;
                case "9":
                    Console.WriteLine("До свидания!");
                    return;
                default:
                    Console.WriteLine("\u001B[31mНеверный выбор.\u001B[0m");
                    break;
            }
        }
    }
}
