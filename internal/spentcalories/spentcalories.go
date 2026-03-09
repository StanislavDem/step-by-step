package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {

	// Разделяем строку через запятую
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("invalid data format")
	}

	// Парсим кол-во шагов
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("steps count must be greater than 0")
	}

	// Вид активности — строка
	activity := parts[1]

	// Парсим продолжительность
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("duration count must be greater than 0")
	}

	// Возвращаем результат кол-ва шагов, вид активности и прод.
	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {

	// Рассчитываем длину шага опираясь на рост пользователя
	stepLen := height * stepLengthCoefficient

	// Умножаем количество шагов на длину шага в м.
	distanceMeters := float64(steps) * stepLen

	// Переводим м. в км.
	distanceKm := distanceMeters / mInKm

	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	// Проверка
	if duration <= 0 {
		return 0
	}

	// Вычисляем дистанцию в км.
	distKm := distance(steps, height)

	// Переводим продолжительность в часы
	durationHours := duration.Hours()

	// Вычисляем среднюю скорость
	speed := distKm / durationHours

	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {

	// Парсим строку
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println("ошибка парсинга:", err)
		return "", err
	}

	// Проверяем тип тренировки
	switch activity {
	case "Бег":
		// дистанция
		distKm := distance(steps, height)
		// скорость
		speed := meanSpeed(steps, height, duration)
		// калории
		calories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		// формирование строки
		return fmt.Sprintf(
			"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			activity,
			duration.Hours(),
			distKm,
			speed,
			calories,
		), nil

	case "Ходьба":
		distKm := distance(steps, height)
		speed := meanSpeed(steps, height, duration)
		calories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(
			"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			activity,
			duration.Hours(),
			distKm,
			speed,
			calories,
		), nil

	default:
		return "", errors.New("неизвестный тип тренировки")
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	// Проверка входных данных
	if steps <= 0 {
		return 0, errors.New("steps count must be greater than 0")
	}
	if weight <= 0 {
		return 0, errors.New("weight must be greater than 0")
	}
	if height <= 0 {
		return 0, errors.New("height must be greater than 0")
	}
	if duration <= 0 {
		return 0, errors.New("duration must be greater than 0")
	}

	// Рассчитываем среднюю скорость
	speed := meanSpeed(steps, height, duration)

	// Переводим продолжительность в минуты
	durationMinutes := duration.Minutes()

	// Расчёт калорий (weight * meanSpeed * durationInMinutes) / minInH
	calories := (weight * speed * durationMinutes) / float64(minInH)

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	// Проверка входных данных
	if steps <= 0 {
		return 0, errors.New("steps count must be greater than 0")
	}
	if weight <= 0 {
		return 0, errors.New("weight must be greater than 0")
	}
	if height <= 0 {
		return 0, errors.New("height must be greater than 0")
	}
	if duration <= 0 {
		return 0, errors.New("duration must be greater than 0")
	}

	// Рассчитываем среднюю скорость
	speed := meanSpeed(steps, height, duration)

	// Переводим продолжительность в минуты
	durationMinutes := duration.Minutes()

	// Расчёт калорий (weight * meanSpeed * durationInMinutes) / minInH
	calories := (weight * speed * durationMinutes) / float64(minInH)

	// Умножаем на корректирующий коэффициент walkingCaloriesCoefficient
	calories *= walkingCaloriesCoefficient

	return calories, nil
}