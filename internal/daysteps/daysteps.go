package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ID-Samurai/go1fl-4-sprint-final/internal/spentcalories"
)

var (
	StepLength = 0.65 // длина шага в метрах
)

func parsePackage(data string) (int, time.Duration, error) {
	sliceData := strings.Split(data, ",")
	if len(sliceData) != 2 {
		return 0, 0, errors.New("Не коректное число входных параметров")
	}
	steps, err := strconv.Atoi(sliceData[0])
	if err != nil {
		return 0, 0, err
	}
	if steps <= 0 {
		return 0, 0, errors.New("Количество шагов должно быть больше 0")
	}
	walkingTime, err := time.ParseDuration(sliceData[1])
	if err != nil {
		return 0, 0, err
	}
	return steps, walkingTime, nil

}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {
	steps, walkingTime, err := parsePackage(data)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	if steps <= 0 {
		return ""
	}
	distance := float64(steps) * StepLength
	distance = distance / 1000
	kcal := spentcalories.WalkingSpentCalories(steps, weight, height, walkingTime)

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, kcal)
}
