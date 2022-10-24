package main

import (
	"encoding/json"
	"net/http"

	"github.com/darkcat013/pr-kitchen/config"
	"github.com/darkcat013/pr-kitchen/domain"
	"github.com/darkcat013/pr-kitchen/utils"
	"go.uber.org/zap"
)

func StartServer() {
	unhandledRoutes := func(w http.ResponseWriter, r *http.Request) {

		utils.Log.Info("Requested",
			zap.String("method", r.Method),
			zap.String("endpoint", r.URL.String()),
		)

		utils.Log.Warn("Path not found", zap.Int("statusCode", http.StatusNotFound))
		http.Error(w, "404 path not found.", http.StatusNotFound)
	}

	order := func(w http.ResponseWriter, r *http.Request) {

		utils.Log.Info("Requested",
			zap.String("method", r.Method),
			zap.String("endpoint", r.URL.String()),
		)

		if r.Method != "POST" {
			utils.Log.Warn("Method not allowed", zap.Int("statusCode", http.StatusMethodNotAllowed))
			http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
			return
		}

		var o domain.Order
		err := json.NewDecoder(r.Body).Decode(&o)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			utils.Log.Fatal("Failed to decode order", zap.String("error", err.Error()))
			return
		}
		utils.Log.Info("Order decoded", zap.Any("order", o))

		domain.NewOrdersChan <- o

		w.WriteHeader(http.StatusOK)

	}

	http.HandleFunc("/", unhandledRoutes)
	http.HandleFunc("/order", order)

	utils.Log.Info("Started web server at port :" + config.PORT)

	if err := http.ListenAndServe(":"+config.PORT, nil); err != nil {
		utils.Log.Fatal("Could not start web server", zap.String("error", err.Error()))
	}
}
