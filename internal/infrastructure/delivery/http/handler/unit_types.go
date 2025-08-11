package handler

type Unit struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetUnitsResponse struct {
	Pageable
	Units []Unit `json:"units"`
}

type CreateUnitRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Image       string `json:"image"`
}
