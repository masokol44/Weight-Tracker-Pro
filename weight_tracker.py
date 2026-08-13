

### 1. `weight_tracker.py` (Python)

```python
# weight_tracker.py — Python версия

import json
import os
import csv
from datetime import datetime, timedelta
from colorama import init, Fore, Style

init(autoreset=True)
DATA_FILE = "weight_data.json"

class WeightEntry:
    def __init__(self, id, date, weight, notes=""):
        self.id = id
        self.date = date
        self.weight = weight
        self.notes = notes

    def to_dict(self):
        return {"id": self.id, "date": self.date, "weight": self.weight, "notes": self.notes}

    @classmethod
    def from_dict(cls, data):
        return cls(data["id"], data["date"], data["weight"], data.get("notes", ""))

class WeightTracker:
    def __init__(self):
        self.entries = []
        self.load()

    def load(self):
        if os.path.exists(DATA_FILE):
            try:
                with open(DATA_FILE, 'r', encoding='utf-8') as f:
                    data = json.load(f)
                    self.entries = [WeightEntry.from_dict(e) for e in data]
            except:
                self.entries = []

    def save(self):
        with open(DATA_FILE, 'w', encoding='utf-8') as f:
            json.dump([e.to_dict() for e in self.entries], f, indent=2, ensure_ascii=False)

    def add(self, date, weight, notes=""):
        id = len(self.entries) + 1
        entry = WeightEntry(id, date, weight, notes)
        self.entries.append(entry)
        self.save()
        return id

    def list_all(self):
        if not self.entries:
            print(Fore.YELLOW + "Нет записей.")
            return
        print(Fore.CYAN + f"{'ID':<4} {'Дата':<12} {'Вес (кг)':<12} {'Заметки':<30}")
        print("-" * 60)
        for e in self.entries:
            print(f"{e.id:<4} {e.date:<12} {e.weight:<12.1f} {e.notes:<30}")

    def get_stats(self):
        if not self.entries:
            return None
        weights = [e.weight for e in self.entries]
        dates = [e.date for e in self.entries]
        sorted_entries = sorted(self.entries, key=lambda x: x.date)
        first_weight = sorted_entries[0].weight
        last_weight = sorted_entries[-1].weight
        change = last_weight - first_weight
        # Вычисляем ИМТ (примерно)
        bmi = None
        # Предполагаем рост 175 см для демонстрации
        height_m = 1.75
        if last_weight > 0:
            bmi = last_weight / (height_m * height_m)
        return {
            'count': len(weights),
            'current': last_weight,
            'average': sum(weights) / len(weights),
            'min': min(weights),
            'max': max(weights),
            'change': change,
            'bmi': bmi
        }

    def print_stats(self):
        stats = self.get_stats()
        if not stats:
            print(Fore.YELLOW + "Нет данных для статистики.")
            return
        print(Fore.CYAN + "📊 Статистика:")
        print(f"  Всего записей: {stats['count']}")
        print(f"  Текущий вес: {stats['current']:.1f} кг")
        print(f"  Средний вес: {stats['average']:.1f} кг")
        print(f"  Минимальный вес: {stats['min']:.1f} кг")
        print(f"  Максимальный вес: {stats['max']:.1f} кг")
        change_color = Fore.GREEN if stats['change'] < 0 else Fore.RED if stats['change'] > 0 else Fore.YELLOW
        print(f"  Изменение: {change_color}{stats['change']:+.1f} кг{Style.RESET_ALL}")
        if stats['bmi']:
            bmi_color = Fore.GREEN if 18.5 <= stats['bmi'] <= 24.9 else Fore.YELLOW if 25 <= stats['bmi'] <= 29.9 else Fore.RED
            print(f"  ИМТ: {bmi_color}{stats['bmi']:.1f}{Style.RESET_ALL}")

    def export_csv(self, filename="weight_data.csv"):
        if not self.entries:
            print(Fore.YELLOW + "Нет данных для экспорта.")
            return
        with open(filename, 'w', newline='', encoding='utf-8') as f:
            writer = csv.writer(f)
            writer.writerow(["ID", "Дата", "Вес (кг)", "Заметки"])
            for e in self.entries:
                writer.writerow([e.id, e.date, e.weight, e.notes])
        print(Fore.GREEN + f"💾 Экспорт CSV: {filename}")

    def export_json(self, filename="weight_data_export.json"):
        if not self.entries:
            print(Fore.YELLOW + "Нет данных для экспорта.")
            return
        data = [e.to_dict() for e in self.entries]
        with open(filename, 'w', encoding='utf-8') as f:
            json.dump(data, f, indent=2, ensure_ascii=False)
        print(Fore.GREEN + f"💾 Экспорт JSON: {filename}")

    def delete(self, id):
        for i, e in enumerate(self.entries):
            if e.id == id:
                del self.entries[i]
                self.save()
                return True
        return False

    def edit(self, id, field, value):
        for e in self.entries:
            if e.id == id:
                if field == "date":
                    e.date = value
                elif field == "weight":
                    e.weight = float(value)
                elif field == "notes":
                    e.notes = value
                else:
                    return False
                self.save()
                return True
        return False

    def filter_by_date(self, start_date, end_date):
        results = [e for e in self.entries if start_date <= e.date <= end_date]
        if not results:
            print(Fore.YELLOW + "Нет записей за указанный период.")
            return
        for e in results:
            print(f"{e.id}: {e.date} | {e.weight:.1f} кг | {e.notes}")

def main():
    tracker = WeightTracker()
    while True:
        print(Fore.CYAN + "\n⚖️ Weight Tracker Pro (Python)")
        print("1. Добавить запись")
        print("2. Показать все записи")
        print("3. Статистика")
        print("4. Экспорт в CSV")
        print("5. Экспорт в JSON")
        print("6. Удалить запись")
        print("7. Редактировать запись")
        print("8. Фильтр по дате")
        print("9. Выход")
        choice = input("Выберите действие: ").strip()
        if choice == "1":
            date = input("Дата (ГГГГ-ММ-ДД, Enter для сегодня): ").strip()
            if not date:
                date = datetime.now().strftime("%Y-%m-%d")
            weight = float(input("Вес (кг): "))
            notes = input("Заметки (опционально): ")
            id = tracker.add(date, weight, notes)
            print(Fore.GREEN + f"✅ Запись добавлена (ID: {id})")
        elif choice == "2":
            tracker.list_all()
        elif choice == "3":
            tracker.print_stats()
        elif choice == "4":
            filename = input("Имя CSV файла (по умолч. weight_data.csv): ").strip()
            if not filename:
                filename = "weight_data.csv"
            tracker.export_csv(filename)
        elif choice == "5":
            filename = input("Имя JSON файла (по умолч. weight_data_export.json): ").strip()
            if not filename:
                filename = "weight_data_export.json"
            tracker.export_json(filename)
        elif choice == "6":
            tracker.list_all()
            id = int(input("Введите ID для удаления: "))
            if tracker.delete(id):
                print(Fore.GREEN + "✅ Запись удалена.")
            else:
                print(Fore.RED + "❌ Запись не найдена.")
        elif choice == "7":
            tracker.list_all()
            id = int(input("Введите ID для редактирования: "))
            field = input("Какое поле редактировать (date, weight, notes): ").lower()
            value = input("Новое значение: ")
            if tracker.edit(id, field, value):
                print(Fore.GREEN + "✅ Запись обновлена.")
            else:
                print(Fore.RED + "❌ Не удалось обновить.")
        elif choice == "8":
            start = input("Начальная дата (ГГГГ-ММ-ДД): ")
            end = input("Конечная дата (ГГГГ-ММ-ДД): ")
            tracker.filter_by_date(start, end)
        elif choice == "9":
            print("До свидания!")
            break
        else:
            print(Fore.RED + "Неверный выбор.")

if __name__ == "__main__":
    main()
