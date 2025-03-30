package daysteps

import (
	"errors"
	"fmt"
	"go-4-sprit-final/internal/spentcalories"
	"strconv"
	"strings"
	"time"
)

const (
	STEP_LENGTH = 0.65 // длина шага в метрах
)

var (
	ErrorConvData =  errors.New("this data isn't correction conv")
	ErrorCorrectionData =  errors.New("this data isn't correction")
	ErrorNullStepData =  errors.New("steps = 0")
)

func parsePackage(data string) (int, time.Duration, error) {
	stepsAndTime := strings.Split(data, ",")

	if len(stepsAndTime) != 2{
		return 0, 0, ErrorCorrectionData
	}

	steps, err := strconv.Atoi(stepsAndTime[0])

	if err != nil{
		return 0, 0, ErrorConvData
	}

	if steps == 0{
		return 0, 0,ErrorNullStepData
	}

	timeDuration, err := time.ParseDuration((stepsAndTime[1]))

	if err != nil{
		return 0, 0, ErrorCorrectionData
	}

	
	return steps, timeDuration, nil
}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {
	steps, time, err := parsePackage(data)

	if err != nil{
		err = fmt.Errorf("Error: %v", err)
		fmt.Println(err)
		return ""
	}

	if steps <= 0{
		err = fmt.Errorf("Error: %v", err)
		fmt.Println(err)
		return ""
	}

	distanceM := float64(steps) * STEP_LENGTH
	distanceKM := distanceM / 1000

	calories  := spentcalories.WalkingSpentCalories(steps, weight, height, time)

	result := fmt.Sprintf("Количество шагов: %d.\n Дистанция составила %.2f км.\nВы сожгли %.4f ккал.", steps, distanceKM, calories)
	return result
}
