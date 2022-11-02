package domain

import (
	"github.com/darkcat013/pr-kitchen/utils"
	"go.uber.org/zap"
)

type Cook struct {
	Id                 int
	Info               CookInfo
	ApparatusFoodsChan chan FinishedFood
}

func NewCook(id int, info *CookInfo) *Cook {
	cook := &Cook{
		Id:                 id,
		Info:               *info,
		ApparatusFoodsChan: make(chan FinishedFood, info.Proficiency),
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
		select {
		case ff := <-c.ApparatusFoodsChan:
			utils.Log.Info("Received finished food from apparatus", zap.Int("cookId", c.Id), zap.Any("details", ff))

			FinishedFoodsChan <- ff
		case sf := <-FoodsChan:
			if !c.CanCook(sf.Food) {
				FoodsChan <- sf
			} else {
				utils.Log.Info("Start preparing", zap.Int("cookId", c.Id), zap.Float64("prepared", sf.PreparedTime), zap.Float64("preparationTime", sf.Food.PreparationTime), zap.Int("orderId", sf.Order.OrderId))

				if sf.Food.CookingApparatus != "" {
					go putFoodInApparatus(sf.Food, c.Id, sf.Order)
					continue
				}

				sleepTime := utils.GetPriorityBasedSleepTime(sf.Order.Priority, sf.Food.PreparationTime)
				if sleepTime+sf.PreparedTime >= sf.Food.PreparationTime {
					sleepTime = sf.Food.PreparationTime - sf.PreparedTime
				}

				utils.SleepFor(sleepTime)
				sf.PreparedTime += sleepTime

				if sf.PreparedTime >= sf.Food.PreparationTime {
					cookingDetail := CookingDetail{
						FoodId: sf.Food.Id,
						CookId: c.Id,
					}

					finishedFood := FinishedFood{
						Details: cookingDetail,
						OrderId: sf.Order.OrderId,
					}
					utils.Log.Info("Finished preparing", zap.Int("cookId", c.Id), zap.Int("foodId", sf.Food.Id), zap.Float64("prepared", sf.PreparedTime), zap.Float64("preparationTime", sf.Food.PreparationTime), zap.Int("orderId", sf.Order.OrderId))

					FinishedFoodsChan <- finishedFood
				} else {
					utils.Log.Info("Pass food", zap.Int("cookId", c.Id), zap.Int("foodId", sf.Food.Id), zap.Float64("prepared", sf.PreparedTime), zap.Float64("preparationTime", sf.Food.PreparationTime), zap.Int("orderId", sf.Order.OrderId))

					FoodsChan <- sf
				}
			}
		}
	}
}

func (c *Cook) CanCook(food *Food) bool {
	return c.Info.Rank >= food.Complexity
}

func putFoodInApparatus(food *Food, cookId int, order *Order) {
	ApparatusesChans[food.CookingApparatus] <- ApparatusFoodInfo{
		Food:     food,
		CookId:   cookId,
		OrderId:  order.OrderId,
		Priority: order.Priority,
	}
}
