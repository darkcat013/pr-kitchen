package domain

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/darkcat013/pr-kitchen/config"
	"github.com/darkcat013/pr-kitchen/utils"
	"go.uber.org/zap"
)

var OrdersListChan = make(chan Order)
var pq = make(PriorityQueue, 0)

func RunOrdersHandling() {
	for {
		o := <-OrdersListChan
		//heap.Push(&pq, &Item{Order: o})
		//heap.Pop(&pq).(*Item).Order
		go SendFoodToCooks(o)
	}
}

func SendFoodToCooks(o Order) {
	orderDone := false
	cookingDetailChan := make(chan CookingDetail)

	//TODO: sort the items by preparation time

	utils.Log.Info("Start preparing order", zap.Any("order", o))
	for _, foodIndex := range o.Items {
		food := Menu[foodIndex-1]
		go SendFood(o.OrderId, food, cookingDetailChan)
	}

	distribution := Distribution{
		OrderId:    o.OrderId,
		TableId:    o.TableId,
		WaiterId:   o.WaiterId,
		Items:      o.Items,
		Priority:   o.Priority,
		MaxWait:    o.MaxWait,
		PickUpTime: o.PickUpTime,
	}

	for !orderDone {
		cookingDetail := <-cookingDetailChan
		distribution.CookingDetails = append(distribution.CookingDetails, cookingDetail)
		if len(distribution.CookingDetails) == len(o.Items) {
			orderDone = true
			distribution.CookingTime = int(time.Now().UnixMilli()) - o.PickUpTime
			utils.Log.Info("Order done", zap.Any("distribution", distribution))
			go SendDistribution(distribution)
		}
	}

}

func SendFood(orderId int, food Food, cookingDetailChan chan CookingDetail) {
	foodSent := false
	for !foodSent {
		for i := 0; i < len(Cooks); i++ {
			if Cooks[i].CanCook(food) {
				foodSent = true
				utils.Log.Info("Sending food to cook", zap.Int("orderId", orderId), zap.Int("foodIndex", food.Id), zap.Int("cookId", Cooks[i].Id))
				cookFood := CookFood{
					CookingDetailChan: cookingDetailChan,
					Food:              food,
					OrderId:           orderId,
				}
				go Cooks[i].StartCooking(cookFood)
				break
			}
		}
	}
}

func SendDistribution(d Distribution) {
	body, err := json.Marshal(d)
	if err != nil {
		utils.Log.Fatal("Failed to convert distribution to JSON ", zap.String("error", err.Error()), zap.Any("distribution", d))
	}

	utils.Log.Info("Send distribuiton to dinner hall", zap.Int("orderId", d.OrderId))

	resp, err := http.Post(config.DINING_HALL_URL, "application/json", bytes.NewBuffer(body))

	if err != nil {
		utils.Log.Error("Failed to send distribution to dinner hall", zap.String("error", err.Error()), zap.Int("orderId", d.OrderId))
	} else {
		utils.Log.Info("Response from dinner hall", zap.Int("statusCode", resp.StatusCode), zap.Int("orderId", d.OrderId))
	}
}
