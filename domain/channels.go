package domain

var NewOrdersChan = make(chan Order, 10)
var FinishedFoodsChan = make(chan FinishedFood)
var ApparatusesChans = make(map[string]chan ApparatusFoodInfo)
var FoodsChan = make(chan *StartedFood, 100)
