package daobase

const (
	DefaultPageSize  = 32
	UnKnowPageNumber = -1
	UnKnowTotalCount = -1
)

type PageMeta struct {
	LastId     IdType `json:"lastId"`
	FirstId    IdType `json:"firstId"`
	PageSize   int32  `json:"pageSize"`
	PageNumber int32  `json:"pageNumber"` //start from 1
	TotalCount int32  `json:"totalCount"`
}

func NewNextPage(lastId IdType) PageMeta {
	return PageMeta{
		LastId:     lastId,
		PageSize:   DefaultPageSize,
		PageNumber: UnKnowPageNumber,
		TotalCount: UnKnowTotalCount,
	}
}

func NewNextPageSize(lastId IdType, pageSize int32) *PageMeta {
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	return &PageMeta{
		LastId:     lastId,
		PageSize:   pageSize,
		PageNumber: UnKnowPageNumber,
		TotalCount: UnKnowTotalCount,
	}
}

func NewPrePage(firstId IdType) PageMeta {
	return PageMeta{
		FirstId:    firstId,
		PageSize:   DefaultPageSize,
		PageNumber: UnKnowPageNumber,
		TotalCount: UnKnowTotalCount,
	}
}

func NewPrePageSize(firstId IdType, pageSize int32) *PageMeta {
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	return &PageMeta{
		FirstId:    firstId,
		PageSize:   pageSize,
		PageNumber: UnKnowPageNumber,
		TotalCount: UnKnowTotalCount,
	}
}

func NewDefaultPage() PageMeta {
	return PageMeta{
		PageSize:   DefaultPageSize,
		PageNumber: UnKnowPageNumber,
		TotalCount: UnKnowTotalCount,
	}
}

func (p *PageMeta) Reverse() *PageMeta {
	p.LastId, p.FirstId = p.FirstId, p.LastId
	return p
}
