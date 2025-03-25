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

// функция парсит строку в количество шагов, вид активности, время, ошибку
func parseTraining(data string) (int, string, time.Duration, error) {
	//убираем лишние пробелы
	data = strings.TrimSpace(data)

	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("Неправильная запись")
	}
	//очищаем каждую часть от пробелов
	part1 := strings.TrimSpace(parts[0])
	part2 := strings.TrimSpace(parts[1])
	part3 := strings.TrimSpace(parts[2])

	if part2 == "" {
		return 0, "", 0, errors.New("Тип активности не может быть пустым")
	}
//переводит строку в число
	num, err := strconv.Atoi(part1)
	if err != nil {
		return 0, "", 0, errors.New("Неправильная запись числа")
	}
	if num <= 0 {
		return 0, "", 0, nil
	}

 //парсит время
	times, err := time.ParseDuration(part3)
	if err != nil {
		return 0, "", 0, err
	}
	if times <= 0 {
		return 0, "", 0, errors.New("Длительность должна быть положительной")
	}
	return num, part2, times, nil
}



// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	distKm := (float64(steps) * lenStep) / float64(mInKm)
	return distKm
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий(число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {
	if duration <= 0 {
		return 0.0
	}
	d := distance(steps)
	return d / duration.Hours()
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {
	// Проверка входных параметров
	if weight <= 0 || height <= 0 {
		return "Ошибка: вес и рост должны быть положительными"
	}
	numOfSteps, activity, timeD, err := parseTraining(data)
	if err != nil {
		return "Ошибка: " + err.Error()
	}

	// Проверка полученных значений
	if numOfSteps <= 0 {
		return "Ошибка: количество шагов должно быть положительным"
	}
	if timeD <= 0 {
		return "Ошибка: длительность тренировки должна быть положительной"
	}

	//варианты активности
	switch activity {
	case "Ходьба":
		str1 := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.1f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f", activity, timeD.Hours(), distance(numOfSteps), meanSpeed(numOfSteps, timeD), WalkingSpentCalories(numOfSteps, weight, height, timeD))
		return str1
	case "Бег":
		str2 := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.1f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f", activity, timeD.Hours(), distance(numOfSteps), meanSpeed(numOfSteps, timeD), RunningSpentCalories(numOfSteps, weight, timeD))
		return str2
	default:
		return "неизвестный тип тренировки "
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
	averageSpeed := meanSpeed(steps, duration)
	return ((runningCaloriesMeanSpeedMultiplier * averageSpeed) - runningCaloriesMeanSpeedShift) * weight

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
	aveSpeed := meanSpeed(steps, duration)
	return ((walkingCaloriesWeightMultiplier * weight) + (aveSpeed*aveSpeed/height)*walkingSpeedHeightMultiplier) * duration.Hours() * float64(minInH)

}
