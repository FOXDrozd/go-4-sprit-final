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

//Errors

var (
	ErrorConvData =  errors.New("this data isn't correction conv")
	ErrorCorrectionData = errors.New("this data isn't correction")
	ErrorNullStepData =  errors.New("steps = 0")
)

func parseTraining(data string) (int, string, time.Duration, error) {
	trainingData := strings.Split(data,",")

	if(len(trainingData) != 3){
		return 0,"", 0 , ErrorCorrectionData
	}

	steps, err := strconv.Atoi(trainingData[0])

	if err != nil{
		return 0,"", 0, ErrorConvData
	}

	if steps == 0{
		return 0, "", 0,ErrorNullStepData
	}

	
	timeDuration, err := time.ParseDuration((trainingData[2]))

	if err != nil{
		return 0, "",0, ErrorCorrectionData
	}

	return steps, trainingData[1], timeDuration, nil
}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	if(steps == 0){
		return 0
	}

	return (float64(steps) * lenStep) / mInKm
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий(число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {
	if(duration == 0){
		return 0
	}

	distanceUser := distance(steps)

	return distanceUser/duration.Hours()
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {
	steps, typeTraining, duration, err:= parseTraining(data)

	if err != nil {
		err = fmt.Errorf("Error: %v", err)
		fmt.Println(err)
		return ""
	}
	
	distanceUser := distance(steps)
	speed := meanSpeed(steps,duration)

	calories := 0.0

	switch typeTraining {
	case "Ходьба": {
		calories = WalkingSpentCalories(steps, weight, height, duration)
	}
case"Бег": {
	
	calories = RunningSpentCalories(steps, weight, duration)
} 
default: {
return "неизвестный тип тренировки"
}
	}
	result := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция:  %.2f км.\nСкорость:  %.2f км/ч\nСожгли калорий  %.2f", typeTraining, duration.Hours(), distanceUser, speed, calories)
	return result
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
   meanSpeedUser := meanSpeed(steps,duration)

   return  ((runningCaloriesMeanSpeedMultiplier*meanSpeedUser)-runningCaloriesMeanSpeedShift) * weight
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
	meanSpeedUser:= meanSpeed(steps, duration)

	return ((walkingCaloriesWeightMultiplier * weight) + (meanSpeedUser*meanSpeedUser/height)*walkingSpeedHeightMultiplier) * duration.Hours() * minInH

}
