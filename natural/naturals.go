package natural

import (
	"fmt"
	"math"
)

type Value struct {
	Num int64
	Div int64
}

// Get literal representation of natural fraction
func (a Value) String() string {
	return fmt.Sprintf("%d/%d", a.Num, a.Div)
}

// Int representation of natural fraction
func (a Value) Int() int64 {
	if a.Div == 0 {
		return 0
	}
	return a.Num / a.Div
}

// Float representation of natural fraction
func (a Value) Float() float64 {
	if a.Div == 0 {
		return 0
	}
	return float64(a.Num) / float64(a.Div)
}

// IsZero checks if natural fraction is zero
func (a Value) IsZero() bool {
	return a.Num <= 0 || a.Div == 0
}

// IsEqual checks if a is equal to b
func (a Value) IsEqual(b Value) bool {
	if a.Div == 0 && b.Div == 0 {
		return a.Num == b.Num
	}
	_, ka, kb := kab(a.Div, b.Div)
	return a.Num*ka == b.Num*kb
}

// IsLessThan checks if a is less than b
func (a Value) IsLessThan(b Value) bool {
	if a.Div == 0 && b.Div == 0 {
		return a.Num < b.Num
	}
	_, ka, kb := kab(a.Div, b.Div)
	return a.Num*ka < b.Num*kb
}

// IsGreaterThan checks if a is greater than b
func (a Value) IsGreaterThan(b Value) bool {
	if a.Div == 0 && b.Div == 0 {
		return a.Num > b.Num
	}
	_, ka, kb := kab(a.Div, b.Div)
	return a.Num*ka > b.Num*kb
}

// LessOrEqualThan checks if a is less or equal to b
func (a Value) IsLessOrEqualThan(b Value) bool {
	return !a.IsGreaterThan(b)
}

// GreaterOrEqualThan checks if a is greater or equal to b
func (a Value) IsGreaterOrEqualThan(b Value) bool {
	return !a.IsLessThan(b)
}

// IsNotEqual checks if a is not equal to b
func (a Value) IsNotEqual(b Value) bool {
	return !a.IsEqual(b)
}

// Negate returns negative natural fraction
func (a Value) Negate() Value {
	return Value{
		Num: -a.Num,
		Div: a.Div,
	}
}

// Abs returns absolute value of natural fraction
func (a Value) Abs() Value {
	return Value{
		Num: abs(a.Num),
		Div: abs(a.Div),
	}
}

// Sign returns sign of natural fraction
func (a *Value) Sign() int64 {
	if a.Num < 0 {
		if a.Div < 0 {
			return 1
		}
		return -1
	}

	if a.Div < 0 {
		return -1
	}

	return 1
}

// Scale natural fraction
func (a Value) Scale(scale int64) Value {
	return Value{
		Num: a.Num * scale,
		Div: a.Div,
	}
}

// Add natural fractions
func (a Value) Add(b Value) Value {
	if a.Div == 0 && b.Div == 0 {
		return Value{Num: a.Num + b.Num}
	}
	div, ka, kb := kab(a.Div, b.Div)
	return Value{
		Num: a.Num*ka + b.Num*kb,
		Div: div,
	}
}

// Subtract natural fractions
func (a Value) Subtract(b Value) Value {
	if a.Div == 0 && b.Div == 0 {
		return Value{Num: a.Num - b.Num, Div: 1}
	}
	div, ka, kb := kab(a.Div, b.Div)
	return Value{
		Num: a.Num*ka - b.Num*kb,
		Div: div,
	}
}

// Multiply natural fractions
func (a Value) Multiple(b Value) Value {
	if a.Div == 0 && b.Div == 0 {
		return Value{Num: a.Num * b.Num}
	}
	ad := coalesce(a.Div, 1)
	bd := coalesce(b.Div, 1)
	res := Value{
		Num: a.Num * b.Num,
		Div: ad * bd,
	}
	return res.Simplify()
}

// Divide natural fractions
func (a Value) Divide(b Value) Value {
	if a.Div == 0 && b.Div == 0 {
		return Value{Num: a.Num, Div: b.Num}
	}
	ad := coalesce(a.Div, 1)
	bd := coalesce(b.Div, 1)
	res := Value{
		Num: a.Num * bd,
		Div: b.Num * ad,
	}
	return res.Simplify()
}

// Привести число b к знаменателю div
func (a Value) Divisor(div int64) (Value, bool) {
	aa := a.Simplify()
	if aa.Div == 0 && div == 0 {
		return Value{Num: aa.Num, Div: 1}, true
	}
	ad := coalesce(aa.Div, 1)
	div = coalesce(div, 1)
	if abs(ad) == abs(div) {
		return aa, true
	}
	if abs(ad) > abs(div) {
		return aa, false
	}
	div2 := gcd(ad, div)
	if div2 != ad && div2 != div {
		return a, false
	}
	return Value{
		Num: aa.Num * div / ad,
		Div: div,
	}, true
}

// Simplify natural fraction
func (a Value) Simplify() Value {
	if a.Div == 0 {
		return Zero
	}
	div := gcd(a.Num, a.Div)
	return Value{
		Num: a.Sign() * abs(a.Num) / div,
		Div: abs(a.Div) / div,
	}
}

// Truncate return truncated natural value without fractional part
func (a Value) Truncate() Value {
	if a.Div == 0 {
		return a
	}

	return Value{
		Num: (a.Num / a.Div) * a.Div,
		Div: a.Div,
	}
}

// Fraction return fractional part of natural value
func (a Value) Fraction() Value {
	if a.Div == 0 {
		return a
	}

	return Value{
		Num: a.Num % a.Div,
		Div: a.Div,
	}
}

// Модуль числа
func abs(value int64) int64 {
	if value < 0 {
		return -value
	}
	return value
}

// Наибольший общий делитель (Алгоритм Евклида)
func gcd(a, b int64) int64 {
	if b == 0 {
		return abs(a)
	}
	return gcd(b, a%b)
}

// Наименьшее общее кратное
func lcm(a, b int64) int64 {
	return abs(a*b) / gcd(a, b)
}

func kab(a, b int64) (int64, int64, int64) {
	a = coalesce(a, 1)
	b = coalesce(b, 1)
	div := lcm(a, b)
	return div, div / a, div / b
}

func coalesce(a, b int64) int64 {
	if a == 0 {
		return b
	}
	return a
}

/*
// Алгоритм Евклида (итеративная версия)
func gcd(a, b int) int {
	for b != 0 {
		a %= b
		if a == 0 {
			return IntAbs(b)
		}
		b %= a
	}
	return intAbs(a)
}
*/

// Преобразование вещественного числа к знаменателю hasDiv
func NewFromFloatWithDivisor(want float64, hasDiv int64) Value {
	num := int64(math.Round(want * float64(hasDiv)))
	return Value{
		Num: num,
		Div: hasDiv,
	}
}

// Преобразование числа в натуральную дробь с указанной точностью
func NewFromFloat(num float64, epsilon float64) Value {
	if num < 0 {
		return newFromFloat(math.Abs(num), epsilon).Negate()
	}
	intval := int64(math.Floor(num))
	res := newFromFloat(num-float64(intval), epsilon)
	return res.Add(Value{intval * res.Div, res.Div})
}

// https://prog-cpp.ru/fraction/
func newFromFloat(num float64, epsilon float64) Value {
	var a int64 = 1
	var b int64 = 1

	if num < epsilon {
		return Value{0, 1}
	}

	// Поиск начального приближения
	c := 1.0
	for {
		b++
		c = float64(a) / float64(b)
		if (num - c) > 0 {
			break
		}
	}

	if (num - c) < epsilon {
		return Value{Num: a, Div: b}
	}

	b--
	c = float64(a) / float64(b)
	if (num - c) > -epsilon {
		return Value{Num: a, Div: b}
	}

	// Уточнение
	var mn int64 = 2 // множитель для начального приближения
	for i := 1; i < 20000; i++ {
		cc := a * mn
		zz := b * mn
		for {
			zz++
			c = float64(cc) / float64(zz)
			if (num - c) > 0 {
				break
			}
		}
		if (num - c) < epsilon {
			return Value{Num: cc, Div: zz}
		}

		zz--
		c = float64(cc) / float64(zz)
		if (num - c) > -epsilon {
			return Value{Num: cc, Div: zz}
		}
		mn++
	}

	return Value{Num: a, Div: b}
}

var Zero = Value{Num: 0, Div: 0}
