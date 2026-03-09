package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {

	// Разделяем строку через запятую
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, errors.New("invalid data format")
	}

	// Парсим кол-во шагов
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	if steps <= 0 {
		return 0, 0, errors.New("steps count must be greater than 0")
	}

	// Парсим продолжительность
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, err
	}
	if duration <= 0 {
		return 0, 0, errors.New("duration count must be greater than 0")
	}

	// Возвращаем результат кол-ва шагов
	return steps, duration, nil

}

func DayActionInfo(data string, weight, height float64) string {

	// Парсим данные о кол-ве шагов и прод. прогулки
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println("Ошибка парсинга:", err)
		return ""
	}

	// Проверяем кол-во шагов
	if steps <= 0 {
		log.Println("Ошибка парсинга: некорректное количество шагов")
		return ""
	}

	// Проверяем продолжительность
	if duration < 0 {
		log.Println("Ошибка парсинга: отрицательная продолжительность")
		return ""
	}

	// Вычисляем дистанцию в метрах
	distanceMeters := float64(steps) * stepLength

	// Переводим в километры
	distanceKm := distanceMeters / mInKm

	// Вычисляем калории через step-by-step/internal/spentcalories
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println("Ошибка вычисления калорий:", err)
		return ""
	}

	// Вывод результата
	result := fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceKm, calories,
	)

	return result
}