package domain

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/darkcat013/pr-kitchen/config"
	"github.com/darkcat013/pr-kitchen/utils"
	"go.uber.org/zap"
)

func RunOrdersHandling() {
	go startOrderTracking()
	go startOrderProcessing()
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
		utils.Log.Info("Start preparing order", zap.Any("order", o))
		Distributions[o.OrderId] = &Distribution{
			OrderId:    o.OrderId,
			TableId:    o.TableId,
			WaiterId:   o.WaiterId,
			Items:      o.Items,
			Priority:   o.Priority,
			MaxWait:    o.MaxWait,
			PickUpTime: o.PickUpTime,
		}
		go sendFoodToCooks(o)
	}
}
func sendFoodToCooks(o Order) {
	foodAmount := len(o.Items)

	for i := 0; i < foodAmount; i++ {
		go sendFood(Menu[o.Items[i]-1], &o)
	}
}

func sendFood(food *Food, order *Order) {
	utils.Log.Info("Sending food to cooks", zap.Int("foodId", food.Id), zap.Int("orderId", order.OrderId))
	FoodsChan <- &StartedFood{
		Food:  food,
		Order: order,
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
