package daysteps

import (
	"fmt"
	"strings"
	"time"

	"strconv"

	"github.com/Yandex-Practicum/go1fl-4-sprint-final/internal/spentcalories"
)

var (
	StepLength = 0.65 // длина шага в метрах
	Km         = 1000 // количество метров в км
)

// парсит строку в количество шагов, время в формате 0h50m, ошибку
func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("ошибка, неправильная запись")
	}
	partNum := strings.TrimSpace(parts[0])
	partTime := strings.TrimSpace(parts[1])

	num, err := strconv.Atoi(partNum)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования числа %s : %w", partNum, err)
	}
	if num <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть положительным")
	}

	times, err := time.ParseDuration(partTime)
	if err != nil {
		return 0, 0, err
	}
	if times <= 0 {
		return 0, 0, fmt.Errorf("длительность должна быть положительной")
	}
	return num, times, nil

}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {
	if weight <= 0 || height <= 0 {
		return "Ошибка: вес и рост должны быть положительными значениями"
	}
	step, timeValk, err := parsePackage(data)
	if err != nil {
		return fmt.Sprintf("Ошибка : %w", err)
	}
	if step == 0 {
		return ""
	}
	distanceInMetr := (float64(step) * StepLength) / float64(Km)

	if timeValk <= 0 {
		return "Ошибка: длительность тренировки должна быть положительной"
	}
	calories := spentcalories.WalkingSpentCalories(step, weight, height, timeValk)

	str := fmt.Sprintf("Количество шагов: %d.\nДистанция составила: %.2f км.\nВы сожгли : %.2f ккал.", step, distanceInMetr, calories)
	return str
}
