package sql

import "fmt"

func ConvertSliceSql(list []uint) (s string) {
	s += "("
	for i, u := range list {
		if i == len(list)-1 {
			s += fmt.Sprintf("%d", u)
			break
		}
		s += fmt.Sprintf("%d,", u)
	}
	s += ")"
	return
}
func ConvertSliceOrderSql(list []uint) (s string) {
	for i, u := range list {
		if i == len(list)-1 {
			s += fmt.Sprintf("id = %d desc", u)
			break
		}
		s += fmt.Sprintf("id = %d desc, ", u)
	}
	return
}
