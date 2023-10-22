package evd_res

import "time"

type EvdsPresentator struct {
	Content   string
	Image     string
	CreatedAt time.Time
}

type GetEvdsRes struct {
	Message string
	Data    []EvdsPresentator
}
