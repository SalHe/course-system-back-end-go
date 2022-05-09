package services

type Page struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

func (p Page) ActualPage() int {
	if p.Page < 1 {
		return 1
	} else {
		return p.Page
	}
}

func (p Page) ActualSize() int {
	if p.Size <= 0 {
		return 10
	} else {
		return p.Size
	}
}

func (p Page) Offset() int {
	return (p.ActualPage() - 1) * p.ActualSize()
}
