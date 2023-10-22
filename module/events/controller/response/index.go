package evt_response

type (
	GetEventResponse struct {
		Message string            `json:"message"`
		Data    []StdPresentation `json:"data"`
	}

	StdPresentation struct {
		ID                string `json:"id"`
		Title             string `json:"title"`
		Location          string `json:"location"`
		LocationURL       string `json:"location_url"`
		Description       string `json:"description"`
		Image             string `json:"image"`
		RecommendedAction string `json:"recommended_action"`
		CreatedBy         string `json:"created_by"`
		Verified          bool   `json:"verified"`
	}
)
