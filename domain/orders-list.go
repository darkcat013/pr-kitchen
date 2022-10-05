package domain

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync/atomic"

	"github.com/darkcat013/pr-kitchen/config"
	"github.com/darkcat013/pr-kitchen/utils"
	"go.uber.org/zap"
)

func RunOrdersHandling() {
	go startOrderTracking()
	go startOrderProcessing()
	go startOrderSending()
}

func startOrderTracking() {
	for {
		ff := <-FinishedFoodsChan
		utils.Log.Info("Received finished food", zap.Any("details", ff))
		distribution := Distributions[ff.OrderId]
		distribution.CookingDetails = append(distribution.CookingDetails, ff.Details)
		if len(distribution.CookingDetails) == len(distribution.Items) {
			distribution.CookingTime = (utils.GetCurrentTimeFloat() - distribution.PickUpTime) / config.TIME_UNIT_COEFF
			go sendDistribution(*distribution)
			delete(Distributions, ff.OrderId)
		}
	}
}

func startOrderProcessing() {
	for {
		o := <-NewOrdersChan
		startedOrder := StartedOrder{
			Order:    &o,
			Items:    o.Items,
			Priority: o.Priority,
		}
		utils.Log.Info("Start preparing order", zap.Any("order", o))
		StartedOrdersChan <- startedOrder
		Distributions[o.OrderId] = &Distribution{
			OrderId:    o.OrderId,
			TableId:    o.TableId,
			WaiterId:   o.WaiterId,
			Items:      o.Items,
			Priority:   o.Priority,
			MaxWait:    o.MaxWait,
			PickUpTime: o.PickUpTime,
		}
	}
}

func startOrderSending() {
	for {
		so := <-StartedOrdersChan
		utils.Log.Info("Continue preparing order", zap.Any("order", so))

		switch so.Priority {
		case 5:
			go sendFoodToCooks(so, 1)
		case 4:
			go sendFoodToCooks(so, 1)
		case 3:
			go sendFoodToCooks(so, 2)
		case 2:
			go sendFoodToCooks(so, 2)
		case 1:
			go sendFoodToCooks(so, 3)
		}
	}
}
func sendFoodToCooks(so StartedOrder, divisor int) {

	foodAmount := len(so.Items) / divisor

	for i := 0; i < foodAmount; i++ {
		var foodId int
		foodId, so.Items = so.Items[0], so.Items[1:]
		go sendFood(Menu[foodId-1], so.Order)
	}

	if len(so.Items) == 0 {
		return
	}

	so.Priority += 2
	utils.Log.Info("Send started order back to channel", zap.Any("order", so))
	StartedOrdersChan <- so
}

func sendFood(food *Food, order *Order) {
	foodSent := false
	for !foodSent {
		for i := 0; i < len(Cooks); i++ {
			if Cooks[i].CanCook(food) {
				utils.Log.Info("Sending food to cook", zap.Int("foodId", food.Id), zap.Int("cookId", Cooks[i].Id), zap.Int64("cookingNow", atomic.LoadInt64(&Cooks[i].CookingNow)))
				Cooks[i].FoodsChan <- &StartedFood{
					Food:  food,
					Order: order,
				}
				foodSent = true
				break
			}
		}
	}
}

func sendDistribution(d Distribution) {
	body, err := json.Marshal(d)
	if err != nil {
		utils.Log.Fatal("Failed to convert distribution to JSON ", zap.String("error", err.Error()), zap.Any("distribution", d))
	}

	utils.Log.Info("Send distribuiton to dinner hall", zap.Any("distribution", d))

	resp, err := http.Post(config.DINING_HALL_URL, "application/json", bytes.NewBuffer(body))

	if err != nil {
		utils.Log.Error("Failed to send distribution to dinner hall", zap.String("error", err.Error()), zap.Any("distribution", d))
	} else {
		utils.Log.Info("Response from dinner hall", zap.Int("statusCode", resp.StatusCode), zap.Any("distribution", d))
	}
}
