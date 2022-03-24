package requestFilter

type FilterItem struct {
	Condition string
	Data      map[string]interface{}
}

type Filter struct {
	Filters []FilterItem
	Limit   uint
	Offset  uint
	Sort    string
	Group   string
}
