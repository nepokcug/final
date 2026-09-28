package nextdate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Проверка на прошлое
func afterNow(date, now time.Time) bool {
	n := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return date.After(n)
}

func nextDateByDays(now time.Time, d time.Time, q string) (string, error) {

	// 1. Парсим число дней
	days, err := strconv.Atoi(q)
	if err != nil {
		return "", errors.New("неверное число дней")
	}

	// 2. Проверяем ограничение
	if days < 1 || days > 400 {
		return "", errors.New("число дней должно быть от 1 до 400")
	}

	for {
		d = d.AddDate(0, 0, days)
		if afterNow(d, now) {
			break
		}
	}

	return d.Format("20060102"), nil
}

func nextDateByYear(now time.Time, d time.Time) (string, error) {
	for {
		d = d.AddDate(1, 0, 0)
		if afterNow(d, now) {
			break
		}
	}
	return d.Format("20060102"), nil
}

func nextDateByWeekdays(now time.Time, q string) (string, error) {

	// 1. Парсим дни недели
	var weekdays []int
	parts := strings.Split(q, ",")
	for _, p := range parts {
		day, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil || day < 1 || day > 7 {
			return "", errors.New("неверный формат дня недели")
		}
		weekdays = append(weekdays, day)
	}

	d := now.AddDate(0, 0, 1) // 3. Начинаем от завтра
	// 4. Ищем ближайщий день
	for {
		dayOfWeek := int(d.Weekday())
		if dayOfWeek == 0 {
			dayOfWeek = 7
		}
		for _, wd := range weekdays {
			if dayOfWeek == wd {
				return d.Format("20060102"), nil
			}
		}
		d = d.AddDate(0, 0, 1)
	}

}

// Функция для перевода дат в слайс
func parseNumbers(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	var numbers []int

	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		num, err := strconv.Atoi(trimmed)
		if err != nil {
			return nil, fmt.Errorf("не удалось преобразовать %q в число", trimmed)
		}
		numbers = append(numbers, num)
	}

	return numbers, nil
}

// Функция для проверки месяца
func contains(numbers []int, target int) bool {
	for _, n := range numbers {
		if n == target {
			return true
		}
	}
	return false
}

// Функция для проверки дня
func matchesDay(d time.Time, days []int) bool {
	// 1. Вычисляем последний день текущего месяца

	lastDay := time.Date(d.Year(), d.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

	// 2. Проверяем каждый день из списка
	for _, day := range days {
		// Обычный день (1..31)
		if day > 0 && d.Day() == day {
			return true
		}
		// Последний день месяца (-1)
		if day == -1 && d.Day() == lastDay {
			return true
		}
		// Предпоследний день месяца (-2)
		if day == -2 && d.Day() == lastDay-1 {
			return true
		}
	}

	// Ни один день не подошёл
	return false
}

func nextDateByMonthDays(now time.Time, d time.Time, rule []string) (string, error) {

	// 1. Парсим дни
	days, err := parseNumbers(rule[1])
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга дней: %w", err)
	}
	// Проверка: 1..31, -1, -2
	for _, day := range days {
		if day == 0 || day > 31 || day < -2 {
			return "", errors.New("день от 1 до 31, или -1, -2")
		}
	}

	// 2. Парсим месяцы
	var months []int
	if len(rule) > 2 {
		months, err = parseNumbers(rule[2])
		if err != nil {
			return "", fmt.Errorf("ошибка парсинга месяцев: %w", err)
		}
		// Проверка: 1..12
		for _, month := range months {
			if month < 1 || month > 12 {
				return "", errors.New("месяц от 1 до 12")
			}
		}
	}

	// 3. Ищем подходящий день
	for {
		d = d.AddDate(0, 0, 1) // начинаем с завтра
		// Если месяцы указаны и мы на нем прибавляем день
		if len(months) > 0 && !contains(months, int(d.Month())) {
			continue
		}
		// Выходим если день совпал
		if matchesDay(d, days) && afterNow(d, now) {
			return d.Format("20060102"), nil
		}
	}
}
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Пустое правило
	if repeat == "" {
		return "", nil // задача удаляется
	}
	// Парсим дату
	d, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", errors.New("неверный формат даты")
	}
	// Определяем правило
	rule := strings.Fields(repeat)

	switch rule[0] {
	case "d":
		if len(rule) != 2 {
			return "", errors.New("правило d требует число")
		}
		return nextDateByDays(now, d, rule[1])
	case "y":
		if len(rule) != 1 {
			return "", errors.New("правило y не требует параметров")
		}
		return nextDateByYear(now, d)
	case "w":
		if len(rule) != 2 {
			return "", errors.New("правило w требует список дней")
		}
		return nextDateByWeekdays(now, rule[1])
	case "m":
		if len(rule) < 2 {
			return "", errors.New("правило m требует дни")
		}
		return nextDateByMonthDays(now, d, rule)
	default:
		return "", fmt.Errorf("неизвестное правило: %s", rule[0])
	}
}
