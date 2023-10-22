package evd_req

type CreateEvdReq struct {
	Content string `json:"content" form:"content" validate:"required,min=5"`
	EventId string `json:"event_id" form:"event_id" validate:"required,min=16"`
}
