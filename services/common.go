package services

type Page struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

func (p Page) Offset() int {
	return (p.Page - 1) * p.Size
}
