<?php
// weight_tracker.php — PHP версия

$dataFile = 'weight_data.json';

function loadEntries() {
    global $dataFile;
    if (file_exists($dataFile)) {
        $json = file_get_contents($dataFile);
        return json_decode($json, true) ?: [];
    }
    return [];
}

function saveEntries($entries) {
    global $dataFile;
    file_put_contents($dataFile, json_encode($entries, JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE));
}

$entries = loadEntries();

function color($text, $code) {
    return "\033[{$code}m{$text}\033[0m";
}

function listAll($entries) {
    if (empty($entries)) {
        echo color("Нет записей.\n", '33');
        return;
    }
    printf(color("%-4s %-12s %-12s %-30s\n", '36'), "ID", "Дата", "Вес (кг)", "Заметки");
    echo str_repeat("-", 60) . "\n";
    foreach ($entries as $e) {
        printf("%-4d %-12s %-12.1f %-30s\n", $e['id'], $e['date'], $e['weight'], $e['notes'] ?? '');
    }
}

function getStats($entries) {
    if (empty($entries)) return null;
    $weights = array_column($entries, 'weight');
    usort($entries, function($a, $b) { return strcmp($a['date'], $b['date']); });
    $first = $entries[0]['weight'];
    $last = $entries[count($entries)-1]['weight'];
    $change = $last - $first;
    $sum = array_sum($weights);
    $avg = $sum / count($weights);
    $min = min($weights);
    $max = max($weights);
    $bmi = null;
    if ($last > 0) {
        $height = 1.75;
        $bmi = $last / ($height * $height);
    }
    return ['count' => count($weights), 'current' => $last, 'average' => $avg, 'min' => $min, 'max' => $max, 'change' => $change, 'bmi' => $bmi];
}

function printStats($entries) {
    $s = getStats($entries);
    if (!$s) {
        echo color("Нет данных для статистики.\n", '33');
        return;
    }
    echo color("📊 Статистика:\n", '36');
    echo "  Всего записей: {$s['count']}\n";
    echo "  Текущий вес: " . round($s['current'], 1) . " кг\n";
    echo "  Средний вес: " . round($s['average'], 1) . " кг\n";
    echo "  Минимальный вес: " . round($s['min'], 1) . " кг\n";
    echo "  Максимальный вес: " . round($s['max'], 1) . " кг\n";
    $changeColor = $s['change'] < 0 ? '32' : ($s['change'] > 0 ? '31' : '33');
    echo "  Изменение: " . color(($s['change'] > 0 ? '+' : '') . round($s['change'], 1) . " кг", $changeColor) . "\n";
    if ($s['bmi']) {
        $bmiColor = $s['bmi'] >= 18.5 && $s['bmi'] <= 24.9 ? '32' : ($s['bmi'] >= 25 && $s['bmi'] <= 29.9 ? '33' : '31');
        echo "  ИМТ: " . color(round($s['bmi'], 1), $bmiColor) . "\n";
    }
}

function exportCSV($entries, $filename = 'weight_data.csv') {
    if (empty($entries)) {
        echo color("Нет данных для экспорта.\n", '33');
        return;
    }
    $fp = fopen($filename, 'w');
    fputcsv($fp, ['ID', 'Дата', 'Вес (кг)', 'Заметки']);
    foreach ($entries as $e) {
        fputcsv($fp, [$e['id'], $e['date'], round($e['weight'], 1), $e['notes'] ?? '']);
    }
    fclose($fp);
    echo color("💾 Экспорт CSV: $filename\n", '32');
}

function exportJSON($entries, $filename = 'weight_data_export.json') {
    if (empty($entries)) {
        echo color("Нет данных для экспорта.\n", '33');
        return;
    }
    file_put_contents($filename, json_encode($entries, JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE));
    echo color("💾 Экспорт JSON: $filename\n", '32');
}

function deleteEntry(&$entries, $id) {
    foreach ($entries as $i => $e) {
        if ($e['id'] == $id) {
            array_splice($entries, $i, 1);
            saveEntries($entries);
            return true;
        }
    }
    return false;
}

function editEntry(&$entries, $id, $field, $value) {
    foreach ($entries as &$e) {
        if ($e['id'] == $id) {
            switch ($field) {
                case 'date': $e['date'] = $value; break;
                case 'weight': $e['weight'] = (float)$value; break;
                case 'notes': $e['notes'] = $value; break;
                default: return false;
            }
            saveEntries($entries);
            return true;
        }
    }
    return false;
}

function filterByDate($entries, $start, $end) {
    $results = array_filter($entries, function($e) use ($start, $end) {
        return $e['date'] >= $start && $e['date'] <= $end;
    });
    if (empty($results)) {
        echo color("Нет записей за указанный период.\n", '33');
        return;
    }
    foreach ($results as $e) {
        echo "{$e['id']}: {$e['date']} | " . round($e['weight'], 1) . " кг | {$e['notes']}\n";
    }
}

function main() {
    global $entries;
    while (true) {
        echo "\n" . color("⚖️ Weight Tracker Pro (PHP)\n", '36');
        echo "1. Добавить запись\n";
        echo "2. Показать все записи\n";
        echo "3. Статистика\n";
        echo "4. Экспорт в CSV\n";
        echo "5. Экспорт в JSON\n";
        echo "6. Удалить запись\n";
        echo "7. Редактировать запись\n";
        echo "8. Фильтр по дате\n";
        echo "9. Выход\n";
        echo "Выберите действие: ";
        $choice = trim(fgets(STDIN));

        switch ($choice) {
            case '1':
                echo "Дата (ГГГГ-ММ-ДД, Enter для сегодня): ";
                $date = trim(fgets(STDIN));
                if (empty($date)) $date = date('Y-m-d');
                echo "Вес (кг): ";
                $weight = (float) trim(fgets(STDIN));
                echo "Заметки (опционально): ";
                $notes = trim(fgets(STDIN));
                $id = count($entries) + 1;
                $entries[] = ['id' => $id, 'date' => $date, 'weight' => $weight, 'notes' => $notes];
                saveEntries($entries);
                echo color("✅ Запись добавлена (ID: $id)\n", '32');
                break;
            case '2':
                listAll($entries);
                break;
            case '3':
                printStats($entries);
                break;
            case '4':
                echo "Имя CSV файла (по умолч. weight_data.csv): ";
                $filename = trim(fgets(STDIN));
                $filename = empty($filename) ? 'weight_data.csv' : $filename;
                exportCSV($entries, $filename);
                break;
            case '5':
                echo "Имя JSON файла (по умолч. weight_data_export.json): ";
                $filename = trim(fgets(STDIN));
                $filename = empty($filename) ? 'weight_data_export.json' : $filename;
                exportJSON($entries, $filename);
                break;
            case '6':
                listAll($entries);
                echo "Введите ID для удаления: ";
                $id = (int) trim(fgets(STDIN));
                if (deleteEntry($entries, $id)) {
                    echo color("✅ Запись удалена.\n", '32');
                } else {
                    echo color("❌ Запись не найдена.\n", '31');
                }
                break;
            case '7':
                listAll($entries);
                echo "Введите ID для редактирования: ";
                $id = (int) trim(fgets(STDIN));
                echo "Какое поле редактировать (date, weight, notes): ";
                $field = trim(fgets(STDIN));
                echo "Новое значение: ";
                $value = trim(fgets(STDIN));
                if (editEntry($entries, $id, $field, $value)) {
                    echo color("✅ Запись обновлена.\n", '32');
                } else {
                    echo color("❌ Не удалось обновить.\n", '31');
                }
                break;
            case '8':
                echo "Начальная дата (ГГГГ-ММ-ДД): ";
                $start = trim(fgets(STDIN));
                echo "Конечная дата (ГГГГ-ММ-ДД): ";
                $end = trim(fgets(STDIN));
                filterByDate($entries, $start, $end);
                break;
            case '9':
                echo "До свидания!\n";
                exit(0);
            default:
                echo color("Неверный выбор.\n", '31');
        }
    }
}

main();
?>
