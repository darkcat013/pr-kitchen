package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/darkcat013/pr-kitchen/constants"
	"github.com/darkcat013/pr-kitchen/domain"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var log *zap.Logger
var menu []domain.Food

func SetMenu(jsonPath string) {
	file, err := os.Open(jsonPath)
	if err != nil {
		log.Fatal("Error opening " + jsonPath)
	}
	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)
	json.Unmarshal(bytes, &menu)

	if menu == nil {
		log.Fatal("Failed to decode menu from " + jsonPath)
	}
	log.Info("Menu decoded and set")
}

func PrepareOrder(o *domain.Order) {
	log.Info("Start preparing order", zap.Int("orderId", o.OrderId), zap.Int("preparationTime", o.MaxWait))
	time.Sleep(time.Duration(o.MaxWait) * constants.TIME_UNIT)
	log.Info("Order prepared", zap.Int("orderId", o.OrderId))

	d := domain.Distribution{
		OrderId:        o.OrderId,
		TableId:        o.TableId,
		WaiterId:       o.WaiterId,
		Items:          o.Items,
		Priority:       o.Priority,
		MaxWait:        o.MaxWait,
		PickUpTime:     o.PickUpTime,
		CookingTime:    o.MaxWait,
		CookingDetails: nil,
	}

	body, err := json.Marshal(d)
	if err != nil {
		log.Fatal("Failed to convert distribution to JSON ", zap.String("error", err.Error()), zap.Any("distribution", d))
	}

	log.Info("Send distribuiton to dinner hall", zap.Int("orderId", d.OrderId))

	resp, err := http.Post(constants.DINNER_HALL_URL, "application/json", bytes.NewBuffer(body))

	if err != nil {
		log.Error("Failed to send distribution to dinner hall", zap.String("error", err.Error()), zap.Int("orderId", d.OrderId))
	} else {
		log.Info("Response from dinner hall", zap.Int("statusCode", resp.StatusCode), zap.Int("orderId", d.OrderId))
	}
}

func main() {
	log = zap.NewExample()
	defer log.Sync()

	SetMenu("config/menu.json")

	router := mux.NewRouter()
	router.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		var o domain.Order
		err := json.NewDecoder(r.Body).Decode(&o)

		if err != nil {
			log.Error("Failed to decode order",
				zap.String("error", err.Error()),
			)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Info("Order decoded", zap.Any("order", o))

		w.WriteHeader(http.StatusOK)

		go PrepareOrder(&o)

	}).Methods("POST")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("Requested",
			zap.String("method", r.Method),
			zap.String("endpoint", r.URL.String()),
		)
		router.ServeHTTP(w, r)
	})

	http.Handle("/", router)
	log.Info("Started web server at port :8080")
	http.ListenAndServe(":8080", handler)
}
