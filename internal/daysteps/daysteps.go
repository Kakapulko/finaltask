package daysteps

import (
	"errors"
	"fmt"
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
	parts := strings.Split(data, ",")

	// должны быть шаги и время
	if len(parts) != 2 {
		return 0, 0, errors.New("неправильный формат строки")
	}

	// перевод шаги в число
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	// шаги> нуля
	if steps <= 0 {
		return 0, 0, errors.New("шагов должно быть больше 0")
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, err
	}

	//время> нуля
	if duration <= 0 {
		return 0, 0, errors.New("время должно быть больше 0")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println("ошибка:", err)
		return ""
	}
	if steps <= 0 {
		return ""
	}

	distance := float64(steps) * stepLength
	distanceKm := distance / mInKm //в км

	//калории
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println("ошибка:", err)
		return ""
	}

	result := fmt.Sprintf("Количество шагов: %d.\n", steps)
	result = result + fmt.Sprintf("Дистанция составила %.2f км.\n", distanceKm)
	result = result + fmt.Sprintf("Вы сожгли %.2f ккал.\n", calories)

	return result
}
