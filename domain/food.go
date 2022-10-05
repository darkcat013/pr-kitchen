package domain

type Food struct {
	Id               int     `json:"id"`
	Name             string  `json:"name"`
	PreparationTime  float64 `json:"preparation-time"`
	Complexity       int     `json:"complexity"`
	CookingApparatus string  `json:"cooking-apparatus"`
}

type StartedFood struct {
	Food  *Food
	Order *Order
}

type FinishedFood struct {
	Details CookingDetail
	OrderId int
}
