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
	Food    *Food
	CookId  int
	OrderId int
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
		cf := <-a.FoodChan
		utils.Log.Info("Received food in apparatus", zap.Any("details", cf), zap.Int("id", a.Id), zap.String("name", a.Name))
		utils.SleepFor(cf.Food.PreparationTime)

		utils.Log.Info("Finished preparing food in apparatus", zap.Any("details", cf), zap.Int("id", a.Id), zap.String("name", a.Name))

		cookingDetail := CookingDetail{
			FoodId: cf.Food.Id,
			CookId: cf.CookId,
		}

		finishedFood := FinishedFood{
			Details: cookingDetail,
			OrderId: cf.OrderId,
		}

		FinishedFoodsChan <- finishedFood

	}
}
