package models

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/samber/lo"
	"strings"
)

type StringArrayField []string

func (o StringArrayField) Contains(item string) bool {
	return lo.Contains(o, item)
}

// Match 使用数组元素作为关键词去匹配传入的字符串
func (o StringArrayField) Match(paragraph string) bool {
	present := lo.ContainsBy[string](o, func(p string) bool {
		return strings.Contains(p, paragraph)
	})
	return present
}

func (o StringArrayField) Value() (driver.Value, error) {
	b, err := json.Marshal(o)
	return string(b), err
}

func (o *StringArrayField) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), o)
}

func (o StringArrayField) GormDataType() string {
	return "json"
}

func (o StringArrayField) ToStringArray() []string {
	r := make([]string, 0)
	for _, s := range o {
		r = append(r, s)
	}
	return r
}
