package main

import (
	"fmt"
	"math"
	"time"
)

const (
	MInKm      = 1000
	MinInHours = 60
	LenStep    = 0.65
	CmInM      = 100
)

type Training struct {
	TrainingType string
	Action       int
	LenStep      float64
	Duration     time.Duration
	Weight       float64
}

func (t Training) distance() float64 {
	return float64(t.Action) * t.LenStep / float64(MInKm)
}

func (t Training) meanSpeed() float64 {
	if t.Duration == 0 {
		return 0
	}
	return t.distance() / t.Duration.Hours()
}

func (t Training) Calories() float64 {
	return 0
}

type InfoMessage struct {
	TrainingType string
	Duration     time.Duration
	Distance     float64
	Speed        float64
	Calories     float64
}

func (t Training) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: t.TrainingType,
		Duration:     t.Duration,
		Distance:     t.distance(),
		Speed:        t.meanSpeed(),
	}
}

func (i InfoMessage) String() string {
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %v мин\nДистанция: %.2f км.\nСр. скорость: %.2f км/ч\nПотрачено ккал: %.2f\n",
		i.TrainingType,
		i.Duration.Minutes(),
		i.Distance,
		i.Speed,
		i.Calories,
	)
}

type CaloriesCalculator interface {
	Calories() float64
	TrainingInfo() InfoMessage
}

const (
	CaloriesMeanSpeedMultiplier = 18
	CaloriesMeanSpeedShift      = 1.79
)

type Running struct {
	Training
}

func (r Running) Calories() float64 {
	if r.Duration == 0 {
		return 0
	}
	caloriesMean := float64(CaloriesMeanSpeedMultiplier)*r.meanSpeed() + CaloriesMeanSpeedShift
	return caloriesMean * r.Weight / float64(MInKm) * r.Duration.Hours() * float64(MinInHours)
}

func (r Running) TrainingInfo() InfoMessage {
	return r.Training.TrainingInfo()
}

const (
	CaloriesWeightMultiplier      = 0.035
	CaloriesSpeedHeightMultiplier = 0.029
	KmHInMsec                     = 0.278
)

type Walking struct {
	Training
	Height float64
}

func (w Walking) Calories() float64 {
	var heightInCm float64

	meanSpeedInMsec := w.meanSpeed() * float64(KmHInMsec)

	if CmInM == 0 {
		heightInCm = 0
	} else {
		heightInCm = w.Height / float64(CmInM)
	}

	caloriesCalc := CaloriesWeightMultiplier*w.Weight + (math.Pow(meanSpeedInMsec, 2)/heightInCm)*
		CaloriesSpeedHeightMultiplier*w.Weight
	durationMin := w.Duration.Hours() * float64(MinInHours)

	return caloriesCalc * durationMin
}

func (w Walking) TrainingInfo() InfoMessage {
	return w.Training.TrainingInfo()
}

const (
	SwimmingLenStep                  = 1.38
	SwimmingCaloriesMeanSpeedShift   = 1.1
	SwimmingCaloriesWeightMultiplier = 2
)

type Swimming struct {
	Training
	LengthPool int
	CountPool  int
}

func (s Swimming) meanSpeed() float64 {
	if s.Duration == 0 {
		return 0
	}
	return float64(s.LengthPool) * float64(s.CountPool) / MInKm / s.Duration.Hours()
}

func (s Swimming) Calories() float64 {
	return (s.meanSpeed() + SwimmingCaloriesMeanSpeedShift) * float64(SwimmingCaloriesWeightMultiplier) * float64(s.Weight) * s.Duration.Hours()
}

func (s Swimming) TrainingInfo() InfoMessage {
	return s.Training.TrainingInfo()
}

func ReadData(training CaloriesCalculator) string {
	calories := training.Calories()
	info := training.TrainingInfo()

	info.Calories = calories

	return fmt.Sprint(info)
}

func main() {

	swimming := Swimming{
		Training: Training{
			TrainingType: "Плавание",
			Action:       2000,
			LenStep:      SwimmingLenStep,
			Duration:     90 * time.Minute,
			Weight:       85,
		},
		LengthPool: 50,
		CountPool:  5,
	}

	fmt.Println(ReadData(swimming))

	walking := Walking{
		Training: Training{
			TrainingType: "Ходьба",
			Action:       20000,
			LenStep:      LenStep,
			Duration:     3*time.Hour + 45*time.Minute,
			Weight:       85,
		},
		Height: 185,
	}

	fmt.Println(ReadData(walking))

	running := Running{
		Training: Training{
			TrainingType: "Бег",
			Action:       5000,
			LenStep:      LenStep,
			Duration:     30 * time.Minute,
			Weight:       85,
		},
	}

	fmt.Println(ReadData(running))

}
