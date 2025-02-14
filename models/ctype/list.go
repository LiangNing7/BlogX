package ctype

import (
	"database/sql/driver"
	"strings"
)

type List []string

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *List) Scan(value interface{}) error {
	val, ok := value.([]uint8)
	if ok {
		if string(val) == "" {
			*j = []string{}
			return nil
		}
		*j = strings.Split(string(val), ",")
	}
	return nil
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (j List) Value() (driver.Value, error) {
	return strings.Join(j, ","), nil
}
