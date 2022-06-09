package utils_test

import (
	"testing"

	"github.com/max-gryshin/local-chain/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetTagValue(t *testing.T) {
	tests := []struct {
		name     string
		e        interface{}
		tagName  string
		expected []interface{}
	}{
		{
			name: "",
			e: struct {
				Demo  string `db:"demo"`
				Demo1 string `json:"demo1"`
				Demo2 string `db:"demo2"`
				Demo3 string `db:"demo3"`
			}{},
			tagName:  "db",
			expected: []interface{}{"demo", "demo2", "demo3"},
		},
		{
			name: "",
			e: struct {
				Demo  string `json:"demo"`
				Demo1 string `json:"demo1"`
				Demo2 string `json:"demo2"`
				Demo3 string `json:"demo3"`
			}{},
			tagName:  "db",
			expected: []interface{}{},
		},
		{
			name: "",
			e: struct {
				Demo  string `json:"demo"`
				Demo1 string `json:"demo1"`
				Demo2 string `json:"demo2"`
				Demo3 string `json:"demo3"`
			}{},
			tagName:  "json",
			expected: []interface{}{"demo", "demo1", "demo2", "demo3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// nolint: scopelint
			actual := utils.GetTagValue(tt.e, tt.tagName)

			// nolint: scopelint
			assert.EqualValues(t, tt.expected, actual)
		})
	}
}
