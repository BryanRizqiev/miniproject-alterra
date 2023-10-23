package evt_response

import evd_res "miniproject-alterra/module/evidence/controller/response"

type (
	GetEventResponse struct {
		Message string              `json:"message"`
		Data    []EventPresentation `json:"data"`
	}

	EventPresentation struct {
		Id                string                     `json:"id"`
		Title             string                     `json:"title"`
		Location          string                     `json:"location"`
		LocationURL       string                     `json:"location_url"`
		Description       string                     `json:"description"`
		Image             string                     `json:"image"`
		RecommendedAction string                     `json:"recommended_action"`
		CreatedBy         string                     `json:"created_by"`
		Verified          bool                       `json:"verified"`
		Evidences         []evd_res.EvdsPresentation `json:"evidences"`
	}
)
