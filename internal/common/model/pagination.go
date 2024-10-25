package model

type Pagination struct {
	Limit int `json:"limit"`
	Page  int `json:"page"` // 1-indexed page number
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		return 10
	}
	return p.Limit
}

func (p *Pagination) GetOffset() int {
	if p.Page == 0 {
		return 0
	}
	return (p.Page - 1) * p.Limit
}
