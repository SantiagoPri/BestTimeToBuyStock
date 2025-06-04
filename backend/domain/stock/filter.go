package stock

type FilterType string

const (
	FilterExact FilterType = "exact"
	FilterILike FilterType = "ilike"
)

// ValidFilters defines which fields can be filtered and how
var ValidFilters = map[string]FilterType{
	"ticker":   FilterExact,
	"name":     FilterILike,
	"category": FilterExact,
	"price":    FilterExact,
}

// QueryParams represents all possible query parameters for stock filtering
type QueryParams struct {
	Page      int               `json:"page"`
	PageSize  int               `json:"pageSize"`
	Filters   map[string]string `json:"filters"`
	SortBy    string            `json:"sortBy"`
	SortOrder string            `json:"sortOrder"` // asc or desc
}

// NewQueryParams creates a new QueryParams with default values
func NewQueryParams() QueryParams {
	return QueryParams{
		Page:      1,
		PageSize:  10,
		Filters:   make(map[string]string),
		SortBy:    "id",
		SortOrder: "asc",
	}
}
