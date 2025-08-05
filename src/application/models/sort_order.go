package models

import "strings"

type SortOrder string

const (
	//ASC  SortOrder = "ASC"
	DESC SortOrder = "DESC"
)

func (so *SortOrder) ToUpper() SortOrder {
	return SortOrder(strings.ToUpper(string(*so)))
}
