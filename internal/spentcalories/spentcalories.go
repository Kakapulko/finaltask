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
	parts := strings.Split(data, ",")
	// проверяем, что есть шаги, тренировка и время
	if len(parts) != 3 {
		return 0, "", 0, errors.New("неправильный формат строки")
	}

	// перевод шагов в число
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, err
	}

	//проверка, что шагов больше нуля
	if steps <= 0 {
		return 0, "", 0, errors.New("шагов должно быть больше 0")
	}

	activity := parts[1]

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}
	//проверка, что время больше нуля
	if duration <= 0 {
		return 0, "", 0, errors.New("время должно быть больше 0")
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	//длина
	stepLength := height * stepLengthCoefficient
	result := float64(steps) * stepLength / mInKm

	return result
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	dist := distance(steps, height)
	speed := dist / duration.Hours()

	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	switch activity {
	case "Ходьба":
		dist := distance(steps, height)
		speed := meanSpeed(steps, height, duration)

		calories, err := WalkingSpentCalories(steps, weight, height, duration) //
		if err != nil {
			return "", err
		}
		// результат
		result := fmt.Sprintf("Тип тренировки: %s\n", activity)
		result = result + fmt.Sprintf("Длительность: %.2f ч.\n", duration.Hours())
		result = result + fmt.Sprintf("Дистанция: %.2f км.\n", dist)
		result = result + fmt.Sprintf("Скорость: %.2f км/ч\n", speed)
		result = result + fmt.Sprintf("Сожгли калорий: %.2f\n", calories)
		return result, nil

	case "Бег":
		dist := distance(steps, height)
		speed := meanSpeed(steps, height, duration)
		calories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}

		result := fmt.Sprintf("Тип тренировки: %s\n", activity)
		result = result + fmt.Sprintf("Длительность: %.2f ч.\n", duration.Hours())
		result = result + fmt.Sprintf("Дистанция: %.2f км.\n", dist)
		result = result + fmt.Sprintf("Скорость: %.2f км/ч\n", speed)
		result = result + fmt.Sprintf("Сожгли калорий: %.2f\n", calories)
		return result, nil

	default:
		return "", errors.New("неизвестный тип тренировки")
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	//проверяем, что все данные валидные
	if steps <= 0 {
		return 0, errors.New("шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, errors.New("время должно быть больше 0")
	}

	speed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()
	calories := (weight * speed * minutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, errors.New("время должно быть больше 0")
	}
	speed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()

	calories := (weight * speed * minutes) / minInH
	calories = calories * walkingCaloriesCoefficient

	return calories, nil
}
