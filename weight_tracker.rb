# weight_tracker.rb — Ruby версия

require 'json'
require 'csv'
require 'date'

DATA_FILE = 'weight_data.json'

class WeightEntry
  attr_accessor :id, :date, :weight, :notes

  def initialize(id, date, weight, notes = '')
    @id = id
    @date = date
    @weight = weight
    @notes = notes
  end

  def to_h
    { id: @id, date: @date, weight: @weight, notes: @notes }
  end

  def self.from_h(h)
    new(h[:id], h[:date], h[:weight], h[:notes] || '')
  end
end

class WeightTracker
  attr_reader :entries

  def initialize
    @entries = []
    load
  end

  def load
    if File.exist?(DATA_FILE)
      begin
        data = JSON.parse(File.read(DATA_FILE), symbolize_names: true)
        @entries = data.map { |e| WeightEntry.from_h(e) }
      rescue
        @entries = []
      end
    end
  end

  def save
    File.write(DATA_FILE, JSON.pretty_generate(@entries.map(&:to_h)))
  end

  def add(date, weight, notes = '')
    id = @entries.size + 1
    @entries << WeightEntry.new(id, date, weight, notes)
    save
    id
  end

  def list_all
    if @entries.empty?
      puts "\e[33mНет записей.\e[0m"
      return
    end
    printf "\e[36m%-4s %-12s %-12s %-30s\e[0m\n", "ID", "Дата", "Вес (кг)", "Заметки"
    puts "-" * 60
    @entries.each do |e|
      puts "%-4d %-12s %-12.1f %-30s" % [e.id, e.date, e.weight, e.notes]
    end
  end

  def stats
    return nil if @entries.empty?
    weights = @entries.map(&:weight)
    sorted = @entries.sort_by(&:date)
    first = sorted.first.weight
    last = sorted.last.weight
    change = last - first
    avg = weights.sum / weights.size.to_f
    min = weights.min
    max = weights.max
    bmi = nil
    if last > 0
      height = 1.75
      bmi = last / (height * height)
    end
    { count: weights.size, current: last, average: avg, min: min, max: max, change: change, bmi: bmi }
  end

  def print_stats
    s = stats
    if s.nil?
      puts "\e[33mНет данных для статистики.\e[0m"
      return
    end
    puts "\e[36m📊 Статистика:\e[0m"
    puts "  Всего записей: #{s[:count]}"
    puts "  Текущий вес: #{s[:current].round(1)} кг"
    puts "  Средний вес: #{s[:average].round(1)} кг"
    puts "  Минимальный вес: #{s[:min].round(1)} кг"
    puts "  Максимальный вес: #{s[:max].round(1)} кг"
    change_color = s[:change] < 0 ? "\e[32m" : s[:change] > 0 ? "\e[31m" : "\e[33m"
    puts "  Изменение: #{change_color}#{s[:change] > 0 ? '+' : ''}#{s[:change].round(1)} кг\e[0m"
    if s[:bmi]
      bmi_color = s[:bmi] >= 18.5 && s[:bmi] <= 24.9 ? "\e[32m" : s[:bmi] >= 25 && s[:bmi] <= 29.9 ? "\e[33m" : "\e[31m"
      puts "  ИМТ: #{bmi_color}#{s[:bmi].round(1)}\e[0m"
    end
  end

  def export_csv(filename = 'weight_data.csv')
    if @entries.empty?
      puts "\e[33mНет данных для экспорта.\e[0m"
      return
    end
    CSV.open(filename, 'w') do |csv|
      csv << ["ID", "Дата", "Вес (кг)", "Заметки"]
      @entries.each do |e|
        csv << [e.id, e.date, e.weight.round(1), e.notes]
      end
    end
    puts "\e[32m💾 Экспорт CSV: #{filename}\e[0m"
  end

  def export_json(filename = 'weight_data_export.json')
    if @entries.empty?
      puts "\e[33mНет данных для экспорта.\e[0m"
      return
    end
    File.write(filename, JSON.pretty_generate(@entries.map(&:to_h)))
    puts "\e[32m💾 Экспорт JSON: #{filename}\e[0m"
  end

  def delete(id)
    found = @entries.find { |e| e.id == id }
    if found
      @entries.delete(found)
      save
      true
    else
      false
    end
  end

  def edit(id, field, value)
    entry = @entries.find { |e| e.id == id }
    return false unless entry
    case field
    when 'date' then entry.date = value
    when 'weight' then entry.weight = value.to_f
    when 'notes' then entry.notes = value
    else return false
    end
    save
    true
  end

  def filter_by_date(start, end_date)
    results = @entries.select { |e| e.date >= start && e.date <= end_date }
    if results.empty?
      puts "\e[33mНет записей за указанный период.\e[0m"
      return
    end
    results.each { |e| puts "#{e.id}: #{e.date} | #{e.weight.round(1)} кг | #{e.notes}" }
  end
end

def main
  tracker = WeightTracker.new
  loop do
    puts "\n\e[36m⚖️ Weight Tracker Pro (Ruby)\e[0m"
    puts "1. Добавить запись"
    puts "2. Показать все записи"
    puts "3. Статистика"
    puts "4. Экспорт в CSV"
    puts "5. Экспорт в JSON"
    puts "6. Удалить запись"
    puts "7. Редактировать запись"
    puts "8. Фильтр по дате"
    puts "9. Выход"
    print "Выберите действие: "
    choice = gets.chomp
    case choice
    when "1"
      print "Дата (ГГГГ-ММ-ДД, Enter для сегодня): "
      date = gets.chomp
      date = Date.today.to_s if date.empty?
      print "Вес (кг): "
      weight = gets.chomp.to_f
      print "Заметки (опционально): "
      notes = gets.chomp
      id = tracker.add(date, weight, notes)
      puts "\e[32m✅ Запись добавлена (ID: #{id})\e[0m"
    when "2"
      tracker.list_all
    when "3"
      tracker.print_stats
    when "4"
      print "Имя CSV файла (по умолч. weight_data.csv): "
      filename = gets.chomp
      filename = 'weight_data.csv' if filename.empty?
      tracker.export_csv(filename)
    when "5"
      print "Имя JSON файла (по умолч. weight_data_export.json): "
      filename = gets.chomp
      filename = 'weight_data_export.json' if filename.empty?
      tracker.export_json(filename)
    when "6"
      tracker.list_all
      print "Введите ID для удаления: "
      id = gets.chomp.to_i
      if tracker.delete(id)
        puts "\e[32m✅ Запись удалена.\e[0m"
      else
        puts "\e[31m❌ Запись не найдена.\e[0m"
      end
    when "7"
      tracker.list_all
      print "Введите ID для редактирования: "
      id = gets.chomp.to_i
      print "Какое поле редактировать (date, weight, notes): "
      field = gets.chomp.downcase
      print "Новое значение: "
      value = gets.chomp
      if tracker.edit(id, field, value)
        puts "\e[32m✅ Запись обновлена.\e[0m"
      else
        puts "\e[31m❌ Не удалось обновить.\e[0m"
      end
    when "8"
      print "Начальная дата (ГГГГ-ММ-ДД): "
      start = gets.chomp
      print "Конечная дата (ГГГГ-ММ-ДД): "
      end_date = gets.chomp
      tracker.filter_by_date(start, end_date)
    when "9"
      puts "До свидания!"
      break
    else
      puts "\e[31mНеверный выбор.\e[0m"
    end
  end
end

main if __FILE__ == $0
