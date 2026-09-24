package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONList 字符串列表 JSON 类型，兼容 MySQL JSON 列。
type JSONList []string

// Value 实现 driver.Valuer。
func (j JSONList) Value() (driver.Value, error) {
	if j == nil {
		return "[]", nil
	}
	return json.Marshal(j)
}

// Scan 实现 sql.Scanner。
func (j *JSONList) Scan(v any) error {
	if v == nil {
		*j = JSONList{}
		return nil
	}
	b, ok := v.([]byte)
	if !ok {
		return errors.New("invalid json bytes")
	}
	return json.Unmarshal(b, j)
}
