package services

// Page 分页信息
type Page struct {
	Page int `json:"page"` // 当前页
	Size int `json:"size"` // 每页大小
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
