package domain

import (
	"sync/atomic"
	"time"

	"github.com/darkcat013/pr-kitchen/config"
	"github.com/darkcat013/pr-kitchen/utils"
	"go.uber.org/zap"
)

type CookFood struct {
	CookingDetailChan chan CookingDetail
	Food              Food
	OrderId           int
}

type Cook struct {
	Id         int
	Info       CookInfo
	CookingNow int64
}

var Cooks []*Cook

func (c *Cook) StartCooking(cookFood CookFood) {

	atomic.AddInt64(&c.CookingNow, 1)
	utils.Log.Info("Start preparing", zap.Int("cookId", c.Id), zap.Int("orderId", cookFood.OrderId), zap.Any("food", cookFood.Food))
	time.Sleep(config.TIME_UNIT * time.Duration(cookFood.Food.PreparationTime))
	utils.Log.Info("Finished preparing", zap.Int("cookId", c.Id), zap.Int("orderId", cookFood.OrderId), zap.Any("food", cookFood.Food))
	atomic.AddInt64(&c.CookingNow, -1)
	cookFood.CookingDetailChan <- CookingDetail{
		FoodId: cookFood.Food.Id,
		CookId: c.Id,
	}
}

func (c *Cook) CanCook(food Food) bool {
	return atomic.LoadInt64(&c.CookingNow) < int64(c.Info.Proficiency) && c.Info.Rank >= food.Complexity
}
