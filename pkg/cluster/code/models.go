package code

type Commit struct {
	ID      string
	Message string
}

// SearchParams contains parameters for searching operation
type SearchParams struct {
	Filter     string
	PageNumber int
	PageSize   int
}
