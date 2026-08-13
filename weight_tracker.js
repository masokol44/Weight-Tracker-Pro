// weight_tracker.js — JavaScript версия

const fs = require('fs');
const readline = require('readline');

const DATA_FILE = 'weight_data.json';

class WeightEntry {
    constructor(id, date, weight, notes = '') {
        this.id = id;
        this.date = date;
        this.weight = weight;
        this.notes = notes;
    }
}

class WeightTracker {
    constructor() {
        this.entries = [];
        this.load();
    }

    load() {
        if (fs.existsSync(DATA_FILE)) {
            try {
                const data = JSON.parse(fs.readFileSync(DATA_FILE, 'utf8'));
                this.entries = data.map(e => new WeightEntry(e.id, e.date, e.weight, e.notes));
            } catch {
                this.entries = [];
            }
        }
    }

    save() {
        fs.writeFileSync(DATA_FILE, JSON.stringify(this.entries, null, 2));
    }

    add(date, weight, notes = '') {
        const id = this.entries.length + 1;
        const entry = new WeightEntry(id, date, weight, notes);
        this.entries.push(entry);
        this.save();
        return id;
    }

    listAll() {
        if (this.entries.length === 0) {
            console.log('\x1b[33mНет записей.\x1b[0m');
            return;
        }
        console.log('\x1b[36m' + 'ID'.padEnd(4) + 'Дата'.padEnd(12) + 'Вес (кг)'.padEnd(12) + 'Заметки'.padEnd(30) + '\x1b[0m');
        console.log('-'.repeat(60));
        for (const e of this.entries) {
            console.log(`${String(e.id).padEnd(4)} ${e.date.padEnd(12)} ${String(e.weight).padEnd(12)} ${e.notes.padEnd(30)}`);
        }
    }

    getStats() {
        if (this.entries.length === 0) return null;
        const weights = this.entries.map(e => e.weight);
        const sorted = [...this.entries].sort((a, b) => a.date.localeCompare(b.date));
        const first = sorted[0].weight;
        const last = sorted[sorted.length - 1].weight;
        const change = last - first;
        const sum = weights.reduce((a, b) => a + b, 0);
        const min = Math.min(...weights);
        const max = Math.max(...weights);
        const avg = sum / weights.length;
        let bmi = null;
        if (last > 0) {
            const height = 1.75;
            bmi = last / (height * height);
        }
        return { count: weights.length, current: last, average: avg, min, max, change, bmi };
    }

    printStats() {
        const stats = this.getStats();
        if (!stats) {
            console.log('\x1b[33mНет данных для статистики.\x1b[0m');
            return;
        }
        console.log('\x1b[36m📊 Статистика:\x1b[0m');
        console.log(`  Всего записей: ${stats.count}`);
        console.log(`  Текущий вес: ${stats.current.toFixed(1)} кг`);
        console.log(`  Средний вес: ${stats.average.toFixed(1)} кг`);
        console.log(`  Минимальный вес: ${stats.min.toFixed(1)} кг`);
        console.log(`  Максимальный вес: ${stats.max.toFixed(1)} кг`);
        const changeColor = stats.change < 0 ? '\x1b[32m' : stats.change > 0 ? '\x1b[31m' : '\x1b[33m';
        console.log(`  Изменение: ${changeColor}${stats.change > 0 ? '+' : ''}${stats.change.toFixed(1)} кг\x1b[0m`);
        if (stats.bmi) {
            const bmiColor = stats.bmi >= 18.5 && stats.bmi <= 24.9 ? '\x1b[32m' : stats.bmi >= 25 && stats.bmi <= 29.9 ? '\x1b[33m' : '\x1b[31m';
            console.log(`  ИМТ: ${bmiColor}${stats.bmi.toFixed(1)}\x1b[0m`);
        }
    }

    exportCSV(filename = 'weight_data.csv') {
        if (this.entries.length === 0) {
            console.log('\x1b[33mНет данных для экспорта.\x1b[0m');
            return;
        }
        let csv = 'ID,Дата,Вес (кг),Заметки\n';
        for (const e of this.entries) {
            csv += `${e.id},${e.date},${e.weight.toFixed(1)},${e.notes}\n`;
        }
        fs.writeFileSync(filename, csv);
        console.log(`\x1b[32m💾 Экспорт CSV: ${filename}\x1b[0m`);
    }

    exportJSON(filename = 'weight_data_export.json') {
        if (this.entries.length === 0) {
            console.log('\x1b[33mНет данных для экспорта.\x1b[0m');
            return;
        }
        fs.writeFileSync(filename, JSON.stringify(this.entries, null, 2));
        console.log(`\x1b[32m💾 Экспорт JSON: ${filename}\x1b[0m`);
    }

    delete(id) {
        const index = this.entries.findIndex(e => e.id === id);
        if (index !== -1) {
            this.entries.splice(index, 1);
            this.save();
            return true;
        }
        return false;
    }

    edit(id, field, value) {
        const entry = this.entries.find(e => e.id === id);
        if (!entry) return false;
        switch (field) {
            case 'date': entry.date = value; break;
            case 'weight': entry.weight = parseFloat(value); break;
            case 'notes': entry.notes = value; break;
            default: return false;
        }
        this.save();
        return true;
    }

    filterByDate(start, end) {
        const results = this.entries.filter(e => e.date >= start && e.date <= end);
        if (results.length === 0) {
            console.log('\x1b[33mНет записей за указанный период.\x1b[0m');
            return;
        }
        for (const e of results) {
            console.log(`${e.id}: ${e.date} | ${e.weight.toFixed(1)} кг | ${e.notes}`);
        }
    }
}

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

const tracker = new WeightTracker();

function ask(question) {
    return new Promise(resolve => rl.question(question, resolve));
}

async function main() {
    while (true) {
        console.log('\x1b[36m\n⚖️ Weight Tracker Pro (JavaScript)\x1b[0m');
        console.log('1. Добавить запись');
        console.log('2. Показать все записи');
        console.log('3. Статистика');
        console.log('4. Экспорт в CSV');
        console.log('5. Экспорт в JSON');
        console.log('6. Удалить запись');
        console.log('7. Редактировать запись');
        console.log('8. Фильтр по дате');
        console.log('9. Выход');
        const choice = await ask('Выберите действие: ');
        switch (choice.trim()) {
            case '1': {
                let date = await ask('Дата (ГГГГ-ММ-ДД, Enter для сегодня): ');
                if (!date) date = new Date().toISOString().split('T')[0];
                const weight = parseFloat(await ask('Вес (кг): '));
                const notes = await ask('Заметки (опционально): ');
                const id = tracker.add(date, weight, notes);
                console.log(`\x1b[32m✅ Запись добавлена (ID: ${id})\x1b[0m`);
                break;
            }
            case '2': tracker.listAll(); break;
            case '3': tracker.printStats(); break;
            case '4': {
                let filename = await ask('Имя CSV файла (по умолч. weight_data.csv): ');
                if (!filename) filename = 'weight_data.csv';
                tracker.exportCSV(filename);
                break;
            }
            case '5': {
                let filename = await ask('Имя JSON файла (по умолч. weight_data_export.json): ');
                if (!filename) filename = 'weight_data_export.json';
                tracker.exportJSON(filename);
                break;
            }
            case '6': {
                tracker.listAll();
                const id = parseInt(await ask('Введите ID для удаления: '));
                if (tracker.delete(id)) {
                    console.log('\x1b[32m✅ Запись удалена.\x1b[0m');
                } else {
                    console.log('\x1b[31m❌ Запись не найдена.\x1b[0m');
                }
                break;
            }
            case '7': {
                tracker.listAll();
                const id = parseInt(await ask('Введите ID для редактирования: '));
                const field = await ask('Какое поле редактировать (date, weight, notes): ');
                const value = await ask('Новое значение: ');
                if (tracker.edit(id, field, value)) {
                    console.log('\x1b[32m✅ Запись обновлена.\x1b[0m');
                } else {
                    console.log('\x1b[31m❌ Не удалось обновить.\x1b[0m');
                }
                break;
            }
            case '8': {
                const start = await ask('Начальная дата (ГГГГ-ММ-ДД): ');
                const end = await ask('Конечная дата (ГГГГ-ММ-ДД): ');
                tracker.filterByDate(start, end);
                break;
            }
            case '9':
                console.log('До свидания!');
                rl.close();
                return;
            default:
                console.log('\x1b[31mНеверный выбор.\x1b[0m');
        }
    }
}

main().catch(console.error);
