package domain

var NewOrdersChan = make(chan Order)
var StartedOrdersChan = make(chan StartedOrder, 10)
var FinishedFoodsChan = make(chan FinishedFood)
var ApparatusesChans = make(map[string]chan ApparatusFoodInfo)
