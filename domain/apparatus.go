package domain

import (
	"github.com/darkcat013/pr-kitchen/utils"
	"go.uber.org/zap"
)

type Apparatus struct {
	Id       int
	Name     string
	FoodChan chan ApparatusFoodInfo
}

type ApparatusFoodInfo struct {
	Food         *Food
	PreparedTime float64
	Priority     int
	CookId       int
	OrderId      int
}

func NewApparatus(id int, name string, channel chan ApparatusFoodInfo) *Apparatus {
	apparatus := &Apparatus{
		Id:       id,
		Name:     name,
		FoodChan: channel,
	}
	go apparatus.Start()
	return apparatus
}

func (a *Apparatus) Start() {
	for {
		af := <-a.FoodChan
		utils.Log.Info("Received food in apparatus", zap.Any("details", af), zap.Int("id", a.Id), zap.String("name", a.Name))

		sleepTime := utils.GetPriorityBasedSleepTime(af.Priority, af.Food.PreparationTime)

		if sleepTime+af.PreparedTime >= af.Food.PreparationTime {
			sleepTime = af.Food.PreparationTime - af.PreparedTime
		}

		utils.SleepFor(sleepTime)
		af.PreparedTime += sleepTime

		if af.PreparedTime >= af.Food.PreparationTime {
			cookingDetail := CookingDetail{
				FoodId: af.Food.Id,
				CookId: af.CookId,
			}

			finishedFood := FinishedFood{
				Details: cookingDetail,
				OrderId: af.OrderId,
			}

			utils.Log.Info("Finished preparing food in apparatus", zap.Any("details", af), zap.Int("id", a.Id), zap.String("name", a.Name))

			Cooks[af.CookId].ApparatusFoodsChan <- finishedFood
		} else {
			utils.Log.Info("Pass food in apparatus", zap.Any("details", af), zap.Int("id", a.Id), zap.String("name", a.Name))

			a.FoodChan <- af
		}
	}
}
