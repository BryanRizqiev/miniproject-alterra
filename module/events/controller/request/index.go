package event_request

type CreateEvtReq struct {
	Title       string `json:"title" form:"title" validate:"required,min=5"`
	Location    string `json:"location" form:"location" validate:"required,min=5"`
	LocationURL string `json:"location_url" form:"location_url"`
	Description string `json:"description" form:"description"`
}

type UpdateEventStatusReq struct {
	EventId string `json:"event_id" form:"event_id" validate:"required,min=16"`
}
