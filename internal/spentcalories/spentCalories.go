package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep = 0.65 // средняя длина шага.
	mInKm   = 1000 // количество метров в километре.
	minInH  = 60   // количество минут в часе.
)

func parseTraining(data string) (int, string, time.Duration, error) {
	sliceData := strings.Split(data, ",")

	if len(sliceData) != 3 {
		return 0, "", 0, errors.New("Количество входных параметров должно равнятся трем ")
	}
	steps, err := strconv.Atoi(sliceData[0])
	if err != nil {
		return 0, "", 0, err
	}
	walkingTime, err := time.ParseDuration(sliceData[2])
	if err != nil {
		return 0, "", 0, err
	}
	return steps, sliceData[1], walkingTime, nil
}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	return (float64(steps) * lenStep) / float64(mInKm)
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий(число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return distance(steps) / duration.Hours()
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {
	steps, typeTrain, timeTrain, err := parseTraining(data)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	hour := timeTrain.Hours()
	dist := distance(steps)
	speed := meanSpeed(steps, timeTrain)

	switch typeTrain {
	case "Бег":
		kcal := RunningSpentCalories(steps, weight, timeTrain)
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %f км.\nСкорость: %f км/ч\nСожгли калорий: %f\n", typeTrain, hour, dist, speed, kcal)
	case "Ходьба":
		kcal := WalkingSpentCalories(steps, weight, height, timeTrain)
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %f км.\nСкорость: %f км/ч\nСожгли калорий: %f\n", typeTrain, hour, dist, speed, kcal)
	default:
		return "неизвестный тип тренировки"
	}

}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0 // среднее количество сжигаемых калорий при беге.
)

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// steps int - количество шагов.
// weight float64 — вес пользователя.
// duration time.Duration — длительность тренировки.
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	speed := meanSpeed(steps, duration)
	return ((runningCaloriesMeanSpeedMultiplier * speed) - runningCaloriesMeanSpeedShift) * weight
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// steps int - количество шагов.
// duration time.Duration — длительность тренировки.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	speed := meanSpeed(steps, duration)
	return ((walkingCaloriesWeightMultiplier * weight) + (speed*speed/height)*walkingSpeedHeightMultiplier) * duration.Hours() * minInH
}
