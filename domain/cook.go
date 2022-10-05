package domain

import (
	"sync/atomic"

	"github.com/darkcat013/pr-kitchen/utils"
	"go.uber.org/zap"
)

type Cook struct {
	Id         int
	Info       CookInfo
	CookingNow int64
	FoodsChan  chan *StartedFood
}

func NewCook(id int, info *CookInfo) *Cook {
	cook := &Cook{
		Id:         id,
		Info:       *info,
		CookingNow: 0,
		FoodsChan:  make(chan *StartedFood, info.Proficiency),
	}

	go cook.StartWorking()
	return cook
}

func (c *Cook) StartWorking() {
	for i := 0; i < c.Info.Proficiency; i++ {
		go c.StartProeficientCooking()
	}
}

func (c *Cook) StartProeficientCooking() {
	for {
		sf := <-c.FoodsChan
		atomic.AddInt64(&c.CookingNow, 1)
		utils.Log.Info("Start preparing", zap.Int("cookId", c.Id), zap.Int("foodId", sf.Food.Id), zap.Int("orderd", sf.Order.OrderId), zap.Int64("cookingNow", atomic.LoadInt64(&c.CookingNow)))

		if sf.Food.CookingApparatus != "" {
			atomic.AddInt64(&c.CookingNow, -1)
			go putFoodInApparatus(sf.Food, c.Id, sf.Order.OrderId)
			continue
		}

		utils.SleepFor(sf.Food.PreparationTime)
		utils.Log.Info("Finished preparing", zap.Int("cookId", c.Id), zap.Int("foodId", sf.Food.Id), zap.Int("orderd", sf.Order.OrderId), zap.Int64("cookingNow", atomic.LoadInt64(&c.CookingNow)))

		cookingDetail := CookingDetail{
			FoodId: sf.Food.Id,
			CookId: c.Id,
		}

		finishedFood := FinishedFood{
			Details: cookingDetail,
			OrderId: sf.Order.OrderId,
		}
		FinishedFoodsChan <- finishedFood
		atomic.AddInt64(&c.CookingNow, -1)

	}
}

func (c *Cook) CanCook(food *Food) bool {
	return atomic.LoadInt64(&c.CookingNow) < int64(c.Info.Proficiency) && c.Info.Rank >= food.Complexity
}

func putFoodInApparatus(food *Food, cookId, orderId int) {
	ApparatusesChans[food.CookingApparatus] <- ApparatusFoodInfo{
		Food:    food,
		CookId:  cookId,
		OrderId: orderId,
	}
}
