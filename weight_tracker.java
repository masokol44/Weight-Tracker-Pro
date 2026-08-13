// weight_tracker.java — Java версия

import java.io.*;
import java.nio.file.*;
import java.time.*;
import java.util.*;

class WeightEntry {
    int id;
    String date;
    double weight;
    String notes;

    WeightEntry(int id, String date, double weight, String notes) {
        this.id = id;
        this.date = date;
        this.weight = weight;
        this.notes = notes;
    }

    String toJson() {
        return String.format("{\"id\":%d,\"date\":\"%s\",\"weight\":%.1f,\"notes\":\"%s\"}",
                id, date, weight, notes != null ? notes : "");
    }
}

public class weight_tracker {
    private static List<WeightEntry> entries = new ArrayList<>();
    private static final String DATA_FILE = "weight_data.json";
    private static Scanner scanner = new Scanner(System.in);

    public static void main(String[] args) {
        load();
        while (true) {
            System.out.println("\n\u001B[36m⚖️ Weight Tracker Pro (Java)\u001B[0m");
            System.out.println("1. Добавить запись");
            System.out.println("2. Показать все записи");
            System.out.println("3. Статистика");
            System.out.println("4. Экспорт в CSV");
            System.out.println("5. Экспорт в JSON");
            System.out.println("6. Удалить запись");
            System.out.println("7. Редактировать запись");
            System.out.println("8. Фильтр по дате");
            System.out.println("9. Выход");
            System.out.print("Выберите действие: ");
            String choice = scanner.nextLine();
            switch (choice) {
                case "1": addEntry(); break;
                case "2": listAll(); break;
                case "3": printStats(); break;
                case "4": exportCSV(); break;
                case "5": exportJSON(); break;
                case "6": deleteEntry(); break;
                case "7": editEntry(); break;
                case "8": filterByDate(); break;
                case "9": System.out.println("До свидания!"); return;
                default: System.out.println("\u001B[31mНеверный выбор.\u001B[0m");
            }
        }
    }

    private static void load() {
        try {
            String content = new String(Files.readAllBytes(Paths.get(DATA_FILE)));
            // Упрощённо, для демонстрации используем пустой список
            entries = new ArrayList<>();
        } catch (IOException e) {
            entries = new ArrayList<>();
        }
    }

    private static void save() {
        try {
            StringBuilder sb = new StringBuilder("[");
            for (int i = 0; i < entries.size(); i++) {
                sb.append(entries.get(i).toJson());
                if (i < entries.size() - 1) sb.append(",");
            }
            sb.append("]");
            Files.write(Paths.get(DATA_FILE), sb.toString().getBytes());
        } catch (IOException e) {
            System.out.println("Ошибка сохранения.");
        }
    }

    private static void addEntry() {
        System.out.print("Дата (ГГГГ-ММ-ДД, Enter для сегодня): ");
        String date = scanner.nextLine();
        if (date.isEmpty()) {
            date = LocalDate.now().toString();
        }
        System.out.print("Вес (кг): ");
        double weight = Double.parseDouble(scanner.nextLine());
        System.out.print("Заметки (опционально): ");
        String notes = scanner.nextLine();
        int id = entries.size() + 1;
        entries.add(new WeightEntry(id, date, weight, notes));
        save();
        System.out.println("\u001B[32m✅ Запись добавлена (ID: " + id + ")\u001B[0m");
    }

    private static void listAll() {
        if (entries.isEmpty()) {
            System.out.println("\u001B[33mНет записей.\u001B[0m");
            return;
        }
        System.out.printf("\u001B[36m%-4s %-12s %-12s %-30s\u001B[0m\n", "ID", "Дата", "Вес (кг)", "Заметки");
        System.out.println("-".repeat(60));
        for (WeightEntry e : entries) {
            System.out.printf("%-4d %-12s %-12.1f %-30s\n", e.id, e.date, e.weight, e.notes);
        }
    }

    private static void printStats() {
        if (entries.isEmpty()) {
            System.out.println("\u001B[33mНет данных для статистики.\u001B[0m");
            return;
        }
        double[] weights = entries.stream().mapToDouble(e -> e.weight).toArray();
        entries.sort(Comparator.comparing(e -> e.date));
        double first = entries.get(0).weight;
        double last = entries.get(entries.size()-1).weight;
        double change = last - first;
        double sum = 0;
        double min = weights[0];
        double max = weights[0];
        for (double w : weights) {
            sum += w;
            if (w < min) min = w;
            if (w > max) max = w;
        }
        double avg = sum / weights.length;
        Double bmi = null;
        if (last > 0) {
            double height = 1.75;
            bmi = last / (height * height);
        }
        System.out.println("\u001B[36m📊 Статистика:\u001B[0m");
        System.out.printf("  Всего записей: %d\n", weights.length);
        System.out.printf("  Текущий вес: %.1f кг\n", last);
        System.out.printf("  Средний вес: %.1f кг\n", avg);
        System.out.printf("  Минимальный вес: %.1f кг\n", min);
        System.out.printf("  Максимальный вес: %.1f кг\n", max);
        String changeColor = change < 0 ? "\u001B[32m" : change > 0 ? "\u001B[31m" : "\u001B[33m";
        System.out.printf("  Изменение: %s%+.1f кг\u001B[0m\n", changeColor, change);
        if (bmi != null) {
            String bmiColor = bmi >= 18.5 && bmi <= 24.9 ? "\u001B[32m" : bmi >= 25 && bmi <= 29.9 ? "\u001B[33m" : "\u001B[31m";
            System.out.printf("  ИМТ: %s%.1f\u001B[0m\n", bmiColor, bmi);
        }
    }

    private static void exportCSV() {
        if (entries.isEmpty()) {
            System.out.println("\u001B[33mНет данных для экспорта.\u001B[0m");
            return;
        }
        try (FileWriter fw = new FileWriter("weight_data.csv")) {
            fw.write("ID,Дата,Вес (кг),Заметки\n");
            for (WeightEntry e : entries) {
                fw.write(e.id + "," + e.date + "," + String.format("%.1f", e.weight) + "," + e.notes + "\n");
            }
            System.out.println("\u001B[32m💾 Экспорт CSV: weight_data.csv\u001B[0m");
        } catch (IOException e) {
            System.out.println("Ошибка экспорта.");
        }
    }

    private static void exportJSON() {
        if (entries.isEmpty()) {
            System.out.println("\u001B[33mНет данных для экспорта.\u001B[0m");
            return;
        }
        try {
            StringBuilder sb = new StringBuilder("[");
            for (int i = 0; i < entries.size(); i++) {
                sb.append(entries.get(i).toJson());
                if (i < entries.size() - 1) sb.append(",");
            }
            sb.append("]");
            Files.write(Paths.get("weight_data_export.json"), sb.toString().getBytes());
            System.out.println("\u001B[32m💾 Экспорт JSON: weight_data_export.json\u001B[0m");
        } catch (IOException e) {
            System.out.println("Ошибка экспорта.");
        }
    }

    private static void deleteEntry() {
        listAll();
        System.out.print("Введите ID для удаления: ");
        int id = Integer.parseInt(scanner.nextLine());
        Iterator<WeightEntry> it = entries.iterator();
        while (it.hasNext()) {
            if (it.next().id == id) {
                it.remove();
                save();
                System.out.println("\u001B[32m✅ Запись удалена.\u001B[0m");
                return;
            }
        }
        System.out.println("\u001B[31m❌ Запись не найдена.\u001B[0m");
    }

    private static void editEntry() {
        listAll();
        System.out.print("Введите ID для редактирования: ");
        int id = Integer.parseInt(scanner.nextLine());
        WeightEntry target = null;
        for (WeightEntry e : entries) {
            if (e.id == id) {
                target = e;
                break;
            }
        }
        if (target == null) {
            System.out.println("\u001B[31m❌ Запись не найдена.\u001B[0m");
            return;
        }
        System.out.print("Какое поле редактировать (date, weight, notes): ");
        String field = scanner.nextLine().toLowerCase();
        System.out.print("Новое значение: ");
        String value = scanner.nextLine();
        boolean ok = true;
        switch (field) {
            case "date": target.date = value; break;
            case "weight": target.weight = Double.parseDouble(value); break;
            case "notes": target.notes = value; break;
            default: ok = false;
        }
        if (ok) {
            save();
            System.out.println("\u001B[32m✅ Запись обновлена.\u001B[0m");
        } else {
            System.out.println("\u001B[31m❌ Не удалось обновить.\u001B[0m");
        }
    }

    private static void filterByDate() {
        System.out.print("Начальная дата (ГГГГ-ММ-ДД): ");
        String start = scanner.nextLine();
        System.out.print("Конечная дата (ГГГГ-ММ-ДД): ");
        String end = scanner.nextLine();
        List<WeightEntry> results = new ArrayList<>();
        for (WeightEntry e : entries) {
            if (e.date.compareTo(start) >= 0 && e.date.compareTo(end) <= 0) {
                results.add(e);
            }
        }
        if (results.isEmpty()) {
            System.out.println("\u001B[33mНет записей за указанный период.\u001B[0m");
            return;
        }
        for (WeightEntry e : results) {
            System.out.printf("%d: %s | %.1f кг | %s\n", e.id, e.date, e.weight, e.notes);
        }
    }
}
