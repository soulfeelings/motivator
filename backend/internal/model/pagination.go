package model

type PaginationRequest struct {
	Page    int `query:"page"`
	PerPage int `query:"per_page"`
}

func (p *PaginationRequest) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PerPage < 1 || p.PerPage > 100 {
		p.PerPage = 20
	}
}

func (p *PaginationRequest) Offset() int {
	return (p.Page - 1) * p.PerPage
}
