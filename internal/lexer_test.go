package internal_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoyowazzap/expression/internal"
)

func Test_Lex(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		tokens     []internal.Token
		errMsg     string
	}{
		{
			name:       "number",
			expression: "42.1",
			tokens:     []internal.Token{{Type: internal.NUMBER, Value: 42.1}},
		},
		{
			name:       "negative number",
			expression: "-42.1",
			tokens:     []internal.Token{{Type: internal.NUMBER, Value: -42.1}},
		},
		{
			name:       "invert expression",
			expression: "-(42 + 43.2)",
			tokens: []internal.Token{
				{Type: internal.MINUS},
				{Type: internal.LEFT_PAREN},
				{Type: internal.NUMBER, Value: float64(42)},
				{Type: internal.PLUS_OP},
				{Type: internal.NUMBER, Value: 43.2},
				{Type: internal.RIGHT_PAREN},
			},
		},
		{
			name:       "path",
			expression: "$.id1.id2['id3'][1][2]",
			tokens:     []internal.Token{{Type: internal.PATH, Value: []interface{}{"id1", "id2", "id3", 1, 2}}},
		},
		{
			name:       "path with if not",
			expression: "$.id1.id2['id3'][1][2]?0",
			tokens: []internal.Token{
				{Type: internal.PATH, Value: []interface{}{"id1", "id2", "id3", 1, 2}},
				{Type: internal.IF_NOT_FOUND_OP},
				{Type: internal.NUMBER, Value: float64(0)},
			},
		},
		{
			name:       "booleans",
			expression: "true && false",
			tokens: []internal.Token{
				{Type: internal.BOOL, Value: true},
				{Type: internal.AND_OP},
				{Type: internal.BOOL, Value: false},
			},
		},
		{
			name:       "other ids",
			expression: "sum product length and or",
			tokens: []internal.Token{
				{Type: internal.SUM_WORD},
				{Type: internal.PRODUCT_WORD},
				{Type: internal.LENGTH_WORD},
				{Type: internal.AND_WORD},
				{Type: internal.OR_WORD},
			},
		},
		{
			name:       "sum expression",
			expression: "sum(1, 2, 3)",
			tokens: []internal.Token{
				{Type: internal.SUM_WORD},
				{Type: internal.LEFT_PAREN},
				{Type: internal.NUMBER, Value: float64(1)},
				{Type: internal.COMMA},
				{Type: internal.NUMBER, Value: float64(2)},
				{Type: internal.COMMA},
				{Type: internal.NUMBER, Value: float64(3)},
				{Type: internal.RIGHT_PAREN},
			},
		},
		{
			name:       "other operators",
			expression: "*/==||,",
			tokens: []internal.Token{
				{Type: internal.TIMES_OP},
				{Type: internal.DIVIDE_OP},
				{Type: internal.EQUAL_OP},
				{Type: internal.OR_OP},
				{Type: internal.COMMA},
			},
		},
		{
			name:       "number comparison operators",
			expression: "< <0 <= > >1 >=",
			tokens: []internal.Token{
				{Type: internal.LESS_THAN_OP},
				{Type: internal.LESS_THAN_OP},
				{Type: internal.NUMBER, Value: float64(0)},
				{Type: internal.LESS_THAN_OR_EQUAL_OP},
				{Type: internal.GREATER_THAN_OP},
				{Type: internal.GREATER_THAN_OP},
				{Type: internal.NUMBER, Value: float64(1)},
				{Type: internal.GREATER_THAN_OR_EQUAL_OP},
			},
		},
		{
			name:       "string with escaped characters",
			expression: " 'this is a string with \\' and \\\\ '   ",
			tokens:     []internal.Token{{Type: internal.STRING, Value: "this is a string with ' and \\ "}},
		},
		{
			name:       "complicated expression",
			expression: "!($.doAThing ? true) && ((length($.arr1) + length($.arr2)) < 3)",
			tokens: []internal.Token{
				{Type: internal.NOT_OP},
				{Type: internal.LEFT_PAREN},
				{Type: internal.PATH, Value: []interface{}{"doAThing"}},
				{Type: internal.IF_NOT_FOUND_OP},
				{Type: internal.BOOL, Value: true},
				{Type: internal.RIGHT_PAREN},
				{Type: internal.AND_OP},
				{Type: internal.LEFT_PAREN},
				{Type: internal.LEFT_PAREN},
				{Type: internal.LENGTH_WORD},
				{Type: internal.LEFT_PAREN},
				{Type: internal.PATH, Value: []interface{}{"arr1"}},
				{Type: internal.RIGHT_PAREN},
				{Type: internal.PLUS_OP},
				{Type: internal.LENGTH_WORD},
				{Type: internal.LEFT_PAREN},
				{Type: internal.PATH, Value: []interface{}{"arr2"}},
				{Type: internal.RIGHT_PAREN},
				{Type: internal.RIGHT_PAREN},
				{Type: internal.LESS_THAN_OP},
				{Type: internal.NUMBER, Value: float64(3)},
				{Type: internal.RIGHT_PAREN},
			},
		},
		{
			name:       "number with extra decimal",
			expression: "42.2.0",
			errMsg:     "unexpected token '.'",
		},
		{
			name:       "end of expression after . in path",
			expression: "$.id1.",
			errMsg:     "unexpected end of input",
		},
		{
			name:       "beginning of id in path isn't valid",
			expression: "$.id1.123",
			errMsg:     "key in path cannot start with '1'",
		},
		{
			name:       "end of expression after [ in path",
			expression: "$.id1[",
			errMsg:     "unexpected end of input",
		},
		{
			name:       "end of expression in string index in path",
			expression: "$.id1['hey there",
			errMsg:     "unexpected end of input",
		},
		{
			name:       "error converting array index in path",
			expression: "$.id1[1.1]",
			errMsg:     "error parsing number 1.1",
		},
		{
			name:       "unexpected token at start of index in path",
			expression: "$.id1[hey]",
			errMsg:     "unexpected token 'h'",
		},
		{
			name:       "end of input after string index in path (no ])",
			expression: "$.id1['key'",
			errMsg:     "unexpected end of input",
		},
		{
			name:       "bad token after string index in path",
			expression: "$.id1['key'?",
			errMsg:     "unexpected token '?'",
		},
		{
			name:       "unrecognized id",
			expression: "what(1, 2, 3)",
			errMsg:     "unexpected identifier \"what\"",
		},
		{
			name:       "end of input after &",
			expression: "true &",
			errMsg:     "unexpected end of input",
		},
		{
			name:       "only one &",
			expression: "true & false",
			errMsg:     "unexpected token ' ' after '&'",
		},
		{
			name:       "end of input after escape in string",
			expression: "'hey\\",
			errMsg:     "unexpected end of input",
		},
		{
			name:       "bad escape character in string",
			expression: "'whats up \\g'",
			errMsg:     "unexpected escaped token 'g' in string",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens, err := internal.Lex(test.expression)

			if test.errMsg == "" && err != nil {
				t.Errorf("got unexpected error: %s", err)
			} else if test.errMsg != "" && err == nil {
				t.Errorf("didn't get error, expected %s", test.errMsg)
			} else if test.errMsg != "" && err != nil && !strings.Contains(err.Error(), test.errMsg) {
				t.Errorf("error didn't contain wanted string %s: got %s", test.errMsg, err)
			}

			if !reflect.DeepEqual(tokens, test.tokens) {
				t.Errorf("tokens didn't match expected, got\n%v\nwant\n%v", tokens, test.tokens)
			}
		})
	}
}
