package gin_ex

import (
	"encoding/json"
	"fmt"
	"testing"
)

type testStruct struct {
	Code Code   `json:"code"`
	Msg  string `json:"msg"`
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected testStruct
		hasError bool
	}{
		{
			name:     "测试 null 情况",
			input:    `null`,
			expected: testStruct{Code: 0, Msg: ""}, // 期望 Code 不变，Msg 为空
			hasError: false,
		},
		{
			name:     "测试合法的 uint 值",
			input:    `{"code": 2, "msg": "OK"}`,
			expected: testStruct{Code: OK, Msg: "OK"},
			hasError: false,
		},
		{
			name:     "测试 NotFound",
			input:    `{"code": 5, "msg": "Not Found"}`,
			expected: testStruct{Code: NotFound, Msg: "Not Found"},
			hasError: false,
		},
		{
			name:     "测试超出最大值的 uint 值",
			input:    `{"code": 20, "msg": "Invalid Code"}`,
			expected: testStruct{Code: 0, Msg: ""},
			hasError: true,
		},
		{
			name:     "测试合法的 string 值",
			input:    `{"code": "OK", "msg": "Success"}`,
			expected: testStruct{Code: OK, Msg: "Success"},
			hasError: false,
		},
		{
			name:     "测试不合法的 string 值",
			input:    `{"code": "InvalidCode", "msg": "Error"}`,
			expected: testStruct{Code: 0, Msg: ""},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("input = %s", tt.input), func(t *testing.T) {
			var ts testStruct
			err := json.Unmarshal([]byte(tt.input), &ts)

			if tt.hasError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.hasError && err != nil {
				t.Errorf("didn't expect error but got: %v", err)
			}
			if ts != tt.expected {
				t.Errorf("expected %+v, got %+v", tt.expected, ts)
			}
		})
	}
}
