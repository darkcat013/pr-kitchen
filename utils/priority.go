package utils

import (
	"math"

	"go.uber.org/zap"
)

func GetPriorityBasedSleepTime(priority int, preparationTime float64) float64 {
	switch priority {
	case 5:
		return preparationTime
	case 4:
		return preparationTime
	case 3:
		return math.Ceil(preparationTime / 2.0)
	case 2:
		return math.Ceil(preparationTime / 2.0)
	case 1:
		return math.Ceil(preparationTime / 3.0)
	}
	Log.Fatal("Invalid priority", zap.Int("priority", priority))
	return float64(priority)
}
