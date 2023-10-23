package evd_res

type EvdsPresentation struct {
	Content   string `json:"content"`
	Image     string `json:"image"`
	CreatedAt string `json:"created_at"`
	CreatedBy string `json:"created_by"`
	Verified  bool   `json:"verified"`
}

type GetEvdsRes struct {
	Message string
	Data    []EvdsPresentation
}
