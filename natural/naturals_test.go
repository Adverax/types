package natural

import (
	assert "github.com/stretchr/testify/require"
	"testing"
)

func TestSign(t *testing.T) {
	type Test struct {
		name string
		a    Value
		c    int64
	}

	tests := []*Test{
		{
			name: "0/1 -> 1",
			a:    Value{Num: 0, Div: 1},
			c:    1,
		},
		{
			name: "1/3 -> 1",
			a:    Value{Num: 1, Div: 3},
			c:    1,
		},
		{
			name: "-1/3 -> -1",
			a:    Value{Num: -1, Div: 3},
			c:    -1,
		},
		{
			name: "5/3 -> 1",
			a:    Value{Num: 5, Div: 3},
			c:    1,
		},
		{
			name: "-5/3 -> -1",
			a:    Value{Num: -5, Div: 3},
			c:    -1,
		},
		{
			name: "5/-3 -> -1",
			a:    Value{Num: 5, Div: -3},
			c:    -1,
		},
		{
			name: "-5/-3 -> 1",
			a:    Value{Num: -5, Div: -3},
			c:    1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.a.Sign()
			assert.Equal(t, test.c, actual)
		})
	}
}

func TestSimplify(t *testing.T) {
	type Test struct {
		name string
		a    Value
		c    Value
	}

	tests := []*Test{
		{
			name: "0/1 -> 0/1",
			a:    Value{Num: 0, Div: 1},
			c:    Value{Num: 0, Div: 1},
		},
		{
			name: "1/3 -> 1/3",
			a:    Value{Num: 1, Div: 3},
			c:    Value{Num: 1, Div: 3},
		},
		{
			name: "2/4 -> 1/2",
			a:    Value{Num: 2, Div: 4},
			c:    Value{Num: 1, Div: 2},
		},
		{
			name: "16/4 -> 4/1",
			a:    Value{Num: 16, Div: 4},
			c:    Value{Num: 4, Div: 1},
		},
		{
			name: "5/3 -> 5/3",
			a:    Value{Num: 5, Div: 3},
			c:    Value{Num: 5, Div: 3},
		},
		{
			name: "-5/-15 -> 1/3",
			a:    Value{Num: -5, Div: -15},
			c:    Value{Num: 1, Div: 3},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.a.Simplify()
			assert.Equal(t, test.c, actual)
		})
	}

}

func TestUnary(t *testing.T) {
	type Test struct {
		name string
		a    Value
		c    Value
		act  string
	}

	tests := []*Test{
		{
			name: "0/1 -> 0/1",
			a:    Value{Num: 0, Div: 1},
			c:    Value{Num: 0, Div: 1},
			act:  "neg",
		},
		{
			name: "1/3 -> -1/3",
			a:    Value{Num: 1, Div: 3},
			c:    Value{Num: -1, Div: 3},
			act:  "neg",
		},

		{
			name: "5/3 -> 5/3",
			a:    Value{Num: 5, Div: 3},
			c:    Value{Num: 5, Div: 3},
			act:  "abs",
		},
		{
			name: "-5/3 -> 5/3",
			a:    Value{Num: -5, Div: 3},
			c:    Value{Num: 5, Div: 3},
			act:  "abs",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var actual Value
			switch test.act {
			case "neg":
				actual = test.a.Negate()
			case "abs":
				actual = test.a.Abs()
			}
			assert.Equal(t, test.c, actual)
		})
	}
}

func TestBinary(t *testing.T) {
	type Test struct {
		name string
		a, b Value
		c    Value
		act  string
	}

	tests := []*Test{
		{
			name: "0/1 + 0/10 = 0/0",
			a:    Value{Num: 0, Div: 1},
			b:    Value{Num: 0, Div: 10},
			c:    Value{Num: 0, Div: 10},
			act:  "add",
		},
		{
			name: "1/3 + 1/3 = 2/3",
			a:    Value{Num: 1, Div: 3},
			b:    Value{Num: 1, Div: 3},
			c:    Value{Num: 2, Div: 3},
			act:  "add",
		},
		{
			name: "5/3 + 2/5 = 31/15",
			a:    Value{Num: 5, Div: 3},
			b:    Value{Num: 2, Div: 5},
			c:    Value{Num: 31, Div: 15},
			act:  "add",
		},

		{
			name: "0/3 - 0/5 = 0/0",
			a:    Value{Num: 0, Div: 3},
			b:    Value{Num: 0, Div: 5},
			c:    Value{Num: 0, Div: 0},
		},
		{
			name: "5/3  - 2/5 = 31/15",
			a:    Value{Num: 5, Div: 3},
			b:    Value{Num: 2, Div: 5},
			c:    Value{Num: 19, Div: 15},
			act:  "sub",
		},

		{
			name: "0/3 * 0/5 = 0/0",
			a:    Value{Num: 0, Div: 3},
			b:    Value{Num: 0, Div: 5},
			c:    Value{Num: 0, Div: 1},
			act:  "mul",
		},
		{
			name: "5/3  * 2/5 = 10/15",
			a:    Value{Num: 5, Div: 3},
			b:    Value{Num: 2, Div: 5},
			c:    Value{Num: 2, Div: 3},
			act:  "mul",
		},

		{
			name: "0/3 / 0/5 = 0/0",
			a:    Value{Num: 0, Div: 3},
			b:    Value{Num: 0, Div: 5},
			c:    Value{Num: 0, Div: 0},
			act:  "div",
		},
		{
			name: "5/3  / 2/5 = 25/6",
			a:    Value{Num: 5, Div: 3},
			b:    Value{Num: 2, Div: 5},
			c:    Value{Num: 25, Div: 6},
			act:  "div",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var actual Value
			switch test.act {
			case "add":
				actual = test.a.Add(test.b)
			case "sub":
				actual = test.a.Subtract(test.b)
			case "mul":
				actual = test.a.Multiple(test.b)
			case "div":
				actual = test.a.Divide(test.b)
			}
			assert.Equal(t, test.c, actual)
		})
	}
}

func TestNaturalFromFloat(t *testing.T) {
	type Test struct {
		src     float64
		epsilon float64
		dst     Value
	}

	tests := map[string]Test{
		"1": {
			src: 0.2,
			dst: Value{Num: 1, Div: 5},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.epsilon == 0 {
				test.epsilon = 0.001
			}

			actual := NewFromFloat(test.src, test.epsilon)
			assert.Equal(t, test.dst, actual)
		})
	}
}

func TestNatural_Divisor(t *testing.T) {
	a := Value{1, 5}
	actual, ok := a.Divisor(3)
	assert.Equal(t, false, ok)
	assert.Equal(t, Value{1, 5}, actual)
}
