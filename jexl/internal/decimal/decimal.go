// Copyright (c) 2026 Harness, Inc.
// Copyright (c) 2015 Spring, Inc.
// Copyright (c) 2013 Oguz Bilgic
// Use of this source code is governed by the MIT license.

package decimal

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// DivisionPrecision is the number of decimal places in the result when it
// doesn't divide exactly.
var DivisionPrecision = 16

// PowPrecisionNegativeExponent specifies the maximum precision of the result
// when calculating decimal power with a negative exponent.
var PowPrecisionNegativeExponent = 16

// Zero is the zero value Decimal.
var Zero = New(0, 1)

var zeroInt = big.NewInt(0)
var oneInt = big.NewInt(1)
var twoInt = big.NewInt(2)
var fiveInt = big.NewInt(5)
var tenInt = big.NewInt(10)

// Decimal represents a fixed-point decimal. It is immutable.
// number = value * 10 ^ exp
type Decimal struct {
	value *big.Int
	exp   int32
}

// New returns a new fixed-point decimal, value * 10 ^ exp.
func New(value int64, exp int32) Decimal {
	return Decimal{
		value: big.NewInt(value),
		exp:   exp,
	}
}

// NewFromInt converts an int64 to Decimal.
func NewFromInt(value int64) Decimal {
	return Decimal{
		value: big.NewInt(value),
		exp:   0,
	}
}

// NewFromBigInt returns a new Decimal from a big.Int, value * 10 ^ exp.
func NewFromBigInt(value *big.Int, exp int32) Decimal {
	return Decimal{
		value: new(big.Int).Set(value),
		exp:   exp,
	}
}

// NewFromString returns a new Decimal from a string representation.
func NewFromString(value string) (Decimal, error) {
	originalInput := value
	var intString string
	var exp int64

	eIndex := strings.IndexAny(value, "Ee")
	if eIndex != -1 {
		expInt, err := strconv.ParseInt(value[eIndex+1:], 10, 32)
		if err != nil {
			if e, ok := err.(*strconv.NumError); ok && e.Err == strconv.ErrRange {
				return Decimal{}, fmt.Errorf("can't convert %s to decimal: fractional part too long", value)
			}
			return Decimal{}, fmt.Errorf("can't convert %s to decimal: exponent is not numeric", value)
		}
		value = value[:eIndex]
		exp = expInt
	}

	pIndex := -1
	vLen := len(value)
	for i := 0; i < vLen; i++ {
		if value[i] == '.' {
			if pIndex > -1 {
				return Decimal{}, fmt.Errorf("can't convert %s to decimal: too many .s", value)
			}
			pIndex = i
		}
	}

	if pIndex == -1 {
		intString = value
	} else {
		if pIndex+1 < vLen {
			intString = value[:pIndex] + value[pIndex+1:]
		} else {
			intString = value[:pIndex]
		}
		expInt := -len(value[pIndex+1:])
		exp += int64(expInt)
	}

	var dValue *big.Int
	if len(intString) <= 18 {
		parsed64, err := strconv.ParseInt(intString, 10, 64)
		if err != nil {
			return Decimal{}, fmt.Errorf("can't convert %s to decimal", value)
		}
		dValue = big.NewInt(parsed64)
	} else {
		dValue = new(big.Int)
		_, ok := dValue.SetString(intString, 10)
		if !ok {
			return Decimal{}, fmt.Errorf("can't convert %s to decimal", value)
		}
	}

	if exp < math.MinInt32 || exp > math.MaxInt32 {
		return Decimal{}, fmt.Errorf("can't convert %s to decimal: fractional part too long", originalInput)
	}

	return Decimal{
		value: dValue,
		exp:   int32(exp),
	}, nil
}

// NewFromFloat converts a float64 to Decimal using strconv for shortest representation.
func NewFromFloat(value float64) Decimal {
	if value == 0 {
		return New(0, 0)
	}
	if math.IsNaN(value) || math.IsInf(value, 0) {
		panic(fmt.Sprintf("Cannot create a Decimal from %v", value))
	}
	str := strconv.FormatFloat(value, 'f', -1, 64)
	d, err := NewFromString(str)
	if err != nil {
		panic(fmt.Sprintf("Cannot create a Decimal from %v: %v", value, err))
	}
	return d
}

// Copy returns a copy of decimal with the same value and exponent.
func (d Decimal) Copy() Decimal {
	d.ensureInitialized()
	return Decimal{
		value: new(big.Int).Set(d.value),
		exp:   d.exp,
	}
}

func (d Decimal) rescale(exp int32) Decimal {
	d.ensureInitialized()

	if d.exp == exp {
		return Decimal{
			new(big.Int).Set(d.value),
			d.exp,
		}
	}

	diff := math.Abs(float64(exp) - float64(d.exp))
	value := new(big.Int).Set(d.value)

	expScale := new(big.Int).Exp(tenInt, big.NewInt(int64(diff)), nil)
	if exp > d.exp {
		value = value.Quo(value, expScale)
	} else if exp < d.exp {
		value = value.Mul(value, expScale)
	}

	return Decimal{
		value: value,
		exp:   exp,
	}
}

// Abs returns the absolute value of the decimal.
func (d Decimal) Abs() Decimal {
	if !d.IsNegative() {
		return d
	}
	d.ensureInitialized()
	d2Value := new(big.Int).Abs(d.value)
	return Decimal{
		value: d2Value,
		exp:   d.exp,
	}
}

// Add returns d + d2.
func (d Decimal) Add(d2 Decimal) Decimal {
	rd, rd2 := RescalePair(d, d2)
	d3Value := new(big.Int).Add(rd.value, rd2.value)
	return Decimal{
		value: d3Value,
		exp:   rd.exp,
	}
}

// Sub returns d - d2.
func (d Decimal) Sub(d2 Decimal) Decimal {
	rd, rd2 := RescalePair(d, d2)
	d3Value := new(big.Int).Sub(rd.value, rd2.value)
	return Decimal{
		value: d3Value,
		exp:   rd.exp,
	}
}

// Neg returns -d.
func (d Decimal) Neg() Decimal {
	d.ensureInitialized()
	val := new(big.Int).Neg(d.value)
	return Decimal{
		value: val,
		exp:   d.exp,
	}
}

// Mul returns d * d2.
func (d Decimal) Mul(d2 Decimal) Decimal {
	d.ensureInitialized()
	d2.ensureInitialized()

	expInt64 := int64(d.exp) + int64(d2.exp)
	if expInt64 > math.MaxInt32 || expInt64 < math.MinInt32 {
		panic(fmt.Sprintf("exponent %v overflows an int32!", expInt64))
	}

	d3Value := new(big.Int).Mul(d.value, d2.value)
	return Decimal{
		value: d3Value,
		exp:   int32(expInt64),
	}
}

// Div returns d / d2. If it doesn't divide exactly, the result will have
// DivisionPrecision digits after the decimal point.
func (d Decimal) Div(d2 Decimal) Decimal {
	return d.DivRound(d2, int32(DivisionPrecision))
}

// QuoRem does division with remainder.
func (d Decimal) QuoRem(d2 Decimal, precision int32) (Decimal, Decimal) {
	d.ensureInitialized()
	d2.ensureInitialized()
	if d2.value.Sign() == 0 {
		panic("decimal division by 0")
	}
	scale := -precision
	e := int64(d.exp) - int64(d2.exp) - int64(scale)
	if e > math.MaxInt32 || e < math.MinInt32 {
		panic("overflow in decimal QuoRem")
	}
	var aa, bb, expo big.Int
	var scalerest int32
	if e < 0 {
		aa = *d.value
		expo.SetInt64(-e)
		bb.Exp(tenInt, &expo, nil)
		bb.Mul(d2.value, &bb)
		scalerest = d.exp
	} else {
		expo.SetInt64(e)
		aa.Exp(tenInt, &expo, nil)
		aa.Mul(d.value, &aa)
		bb = *d2.value
		scalerest = scale + d2.exp
	}
	var q, r big.Int
	q.QuoRem(&aa, &bb, &r)
	dq := Decimal{value: &q, exp: scale}
	dr := Decimal{value: &r, exp: scalerest}
	return dq, dr
}

// DivRound divides and rounds to a given precision.
func (d Decimal) DivRound(d2 Decimal, precision int32) Decimal {
	q, r := d.QuoRem(d2, precision)
	var rv2 big.Int
	rv2.Abs(r.value)
	rv2.Lsh(&rv2, 1)
	r2 := Decimal{value: &rv2, exp: r.exp + precision}
	var c = r2.Cmp(d2.Abs())

	if c < 0 {
		return q
	}

	if d.value.Sign()*d2.value.Sign() < 0 {
		return q.Sub(New(1, -precision))
	}

	return q.Add(New(1, -precision))
}

// Mod returns d % d2.
func (d Decimal) Mod(d2 Decimal) Decimal {
	_, r := d.QuoRem(d2, 0)
	return r
}

// Pow returns d to the power of d2.
func (d Decimal) Pow(d2 Decimal) Decimal {
	baseSign := d.Sign()
	expSign := d2.Sign()

	if baseSign == 0 {
		if expSign == 0 {
			return Decimal{}
		}
		if expSign == 1 {
			return Decimal{zeroInt, 0}
		}
		if expSign == -1 {
			return Decimal{}
		}
	}

	if expSign == 0 {
		return Decimal{oneInt, 0}
	}

	one := Decimal{oneInt, 0}
	expIntPart, expFracPart := d2.QuoRem(one, 0)

	if baseSign == -1 && !expFracPart.IsZero() {
		return Decimal{}
	}

	intPartPow, _ := d.PowBigInt(expIntPart.value)

	if expFracPart.value.Sign() == 0 {
		return intPartPow
	}

	digitsBase := d.NumDigits()
	digitsExponent := d2.NumDigits()

	precision := digitsBase
	if digitsExponent > precision {
		precision += digitsExponent
	}
	precision += 6

	fracPartPow, err := d.Abs().Ln(-d.exp + int32(precision))
	if err != nil {
		return Decimal{}
	}

	fracPartPow = fracPartPow.Mul(expFracPart)

	fracPartPow, err = fracPartPow.ExpTaylor(-d.exp + int32(precision))
	if err != nil {
		return Decimal{}
	}

	return intPartPow.Mul(fracPartPow)
}

// PowBigInt returns d to the power of exp, where exp is *big.Int.
func (d Decimal) PowBigInt(exp *big.Int) (Decimal, error) {
	return d.powBigIntWithPrecision(exp, int32(PowPrecisionNegativeExponent))
}

func (d Decimal) powBigIntWithPrecision(exp *big.Int, precision int32) (Decimal, error) {
	if d.IsZero() && exp.Sign() == 0 {
		return Decimal{}, fmt.Errorf("cannot represent undefined value of 0**0")
	}

	tmpExp := new(big.Int).Set(exp)
	isExpNeg := exp.Sign() < 0

	if isExpNeg {
		tmpExp.Abs(tmpExp)
	}

	n, result := d, New(1, 0)

	for tmpExp.Sign() > 0 {
		if tmpExp.Bit(0) == 1 {
			result = result.Mul(n)
		}
		tmpExp.Rsh(tmpExp, 1)

		if tmpExp.Sign() > 0 {
			n = n.Mul(n)
		}
	}

	if isExpNeg {
		return New(1, 0).DivRound(result, precision), nil
	}

	return result, nil
}

// ExpTaylor calculates e^d using Taylor series.
func (d Decimal) ExpTaylor(precision int32) (Decimal, error) {
	if d.IsZero() {
		return Decimal{oneInt, 0}.Round(precision), nil
	}

	var epsilon Decimal
	var divPrecision int32
	if precision < 0 {
		epsilon = New(1, -1)
		divPrecision = 8
	} else {
		epsilon = New(1, -precision-1)
		divPrecision = precision + 1
	}

	decAbs := d.Abs()
	pow := d.Abs()
	factorial := New(1, 0)
	factorials := []Decimal{New(1, 0)}

	result := New(1, 0)

	for i := int64(1); ; {
		step := pow.DivRound(factorial, divPrecision)
		result = result.Add(step)

		if step.Cmp(epsilon) < 0 {
			break
		}

		pow = pow.Mul(decAbs)
		i++

		if len(factorials) >= int(i) && !factorials[i-1].IsZero() {
			factorial = factorials[i-1]
		} else {
			factorial = factorials[i-2].Mul(New(i, 0))
			factorials = append(factorials, Zero)
			factorials[i-1] = factorial
		}
	}

	if d.Sign() < 0 {
		result = New(1, 0).DivRound(result, precision+1)
	}

	result = result.Round(precision)
	return result, nil
}

// Ln calculates natural logarithm of d.
func (d Decimal) Ln(precision int32) (Decimal, error) {
	if d.IsNegative() {
		return Decimal{}, fmt.Errorf("cannot calculate natural logarithm for negative decimals")
	}
	if d.IsZero() {
		return Decimal{}, fmt.Errorf("cannot represent natural logarithm of 0, result: -infinity")
	}

	calcPrecision := precision + 2
	z := d.Copy()

	var comp1, comp3, comp2, comp4, reduceAdjust Decimal
	comp1 = z.Sub(Decimal{oneInt, 0})
	comp3 = Decimal{oneInt, -1}

	usePowerSeries := false

	if comp1.Abs().Cmp(comp3) <= 0 {
		usePowerSeries = true
	} else {
		expDelta := int32(z.NumDigits()) + z.exp
		z.exp -= expDelta

		ln10val := ln10.withPrecision(calcPrecision)
		reduceAdjust = NewFromInt(int64(expDelta))
		reduceAdjust = reduceAdjust.Mul(ln10val)

		comp1 = z.Sub(Decimal{oneInt, 0})

		if comp1.Abs().Cmp(comp3) <= 0 {
			usePowerSeries = true
		} else {
			zFloat := z.InexactFloat64()
			comp1 = NewFromFloat(math.Log(zFloat))
		}
	}

	epsilon := Decimal{oneInt, -calcPrecision}

	if usePowerSeries {
		comp2 = comp1.Add(Decimal{twoInt, 0})
		comp3 = comp1.DivRound(comp2, calcPrecision)
		comp1 = comp3.Add(comp3)
		comp2 = comp1.Copy()

		for n := 1; ; n++ {
			comp2 = comp2.Mul(comp3).Mul(comp3)
			comp4 = NewFromInt(int64(2*n + 1))
			comp4 = comp2.DivRound(comp4, calcPrecision)
			comp1 = comp1.Add(comp4)

			if comp4.Abs().Cmp(epsilon) <= 0 {
				break
			}
		}
	} else {
		var prevStep Decimal
		maxIters := calcPrecision*2 + 10

		for i := int32(0); i < maxIters; i++ {
			comp3, _ = comp1.ExpTaylor(calcPrecision)
			comp2 = comp3.Sub(z)
			comp2 = comp2.Add(comp2)
			comp4 = comp3.Add(z)
			comp3 = comp2.DivRound(comp4, calcPrecision)
			comp1 = comp1.Sub(comp3)

			if prevStep.Add(comp3).IsZero() {
				break
			}

			if comp3.Abs().Cmp(epsilon) <= 0 {
				break
			}

			prevStep = comp3
		}
	}

	comp1 = comp1.Add(reduceAdjust)
	return comp1.Round(precision), nil
}

// NumDigits returns the number of digits of the decimal coefficient.
func (d Decimal) NumDigits() int {
	if d.value == nil {
		return 1
	}

	if d.value.IsInt64() {
		i64 := d.value.Int64()
		if i64 <= (1<<53) && i64 >= -(1<<53) {
			if i64 == 0 {
				return 1
			}
			return int(math.Log10(math.Abs(float64(i64)))) + 1
		}
	}

	estimatedNumDigits := int(float64(d.value.BitLen()) / math.Log2(10))
	digitsBigInt := big.NewInt(int64(estimatedNumDigits))
	errorCorrectionUnit := digitsBigInt.Exp(tenInt, digitsBigInt, nil)

	if d.value.CmpAbs(errorCorrectionUnit) >= 0 {
		return estimatedNumDigits + 1
	}

	return estimatedNumDigits
}

// Cmp compares the numbers represented by d and d2 and returns:
//
//	-1 if d <  d2
//	 0 if d == d2
//	+1 if d >  d2
func (d Decimal) Cmp(d2 Decimal) int {
	d.ensureInitialized()
	d2.ensureInitialized()

	if d.exp == d2.exp {
		return d.value.Cmp(d2.value)
	}

	rd, rd2 := RescalePair(d, d2)
	return rd.value.Cmp(rd2.value)
}

// Equal returns whether the numbers represented by d and d2 are equal.
func (d Decimal) Equal(d2 Decimal) bool {
	return d.Cmp(d2) == 0
}

// GreaterThan returns true when d is greater than d2.
func (d Decimal) GreaterThan(d2 Decimal) bool {
	return d.Cmp(d2) == 1
}

// LessThan returns true when d is less than d2.
func (d Decimal) LessThan(d2 Decimal) bool {
	return d.Cmp(d2) == -1
}

// LessThanOrEqual returns true when d is less than or equal to d2.
func (d Decimal) LessThanOrEqual(d2 Decimal) bool {
	cmp := d.Cmp(d2)
	return cmp == -1 || cmp == 0
}

// Sign returns:
//
//	-1 if d <  0
//	 0 if d == 0
//	+1 if d >  0
func (d Decimal) Sign() int {
	if d.value == nil {
		return 0
	}
	return d.value.Sign()
}

// IsNegative returns true if d < 0.
func (d Decimal) IsNegative() bool {
	return d.Sign() == -1
}

// IsZero returns true if d == 0.
func (d Decimal) IsZero() bool {
	return d.Sign() == 0
}

// Exponent returns the exponent component of the decimal.
func (d Decimal) Exponent() int32 {
	return d.exp
}

// IntPart returns the integer component of the decimal.
func (d Decimal) IntPart() int64 {
	scaledD := d.rescale(0)
	return scaledD.value.Int64()
}

// Float64 returns the nearest float64 value for d.
func (d Decimal) Float64() (f float64, exact bool) {
	return d.Rat().Float64()
}

// InexactFloat64 returns the nearest float64 value for d.
func (d Decimal) InexactFloat64() float64 {
	f, _ := d.Float64()
	return f
}

// Rat returns a rational number representation of the decimal.
func (d Decimal) Rat() *big.Rat {
	d.ensureInitialized()
	if d.exp <= 0 {
		denom := new(big.Int).Exp(tenInt, big.NewInt(-int64(d.exp)), nil)
		return new(big.Rat).SetFrac(d.value, denom)
	}

	mul := new(big.Int).Exp(tenInt, big.NewInt(int64(d.exp)), nil)
	num := new(big.Int).Mul(d.value, mul)
	return new(big.Rat).SetFrac(num, oneInt)
}

// String returns the string representation of the decimal.
func (d Decimal) String() string {
	return d.string(true)
}

// Round rounds the decimal to places decimal places.
func (d Decimal) Round(places int32) Decimal {
	if d.exp == -places {
		return d
	}
	ret := d.rescale(-places - 1)

	if ret.value.Sign() < 0 {
		ret.value.Sub(ret.value, fiveInt)
	} else {
		ret.value.Add(ret.value, fiveInt)
	}

	_, m := ret.value.DivMod(ret.value, tenInt, new(big.Int))
	ret.exp++
	if ret.value.Sign() < 0 && m.Cmp(zeroInt) != 0 {
		ret.value.Add(ret.value, oneInt)
	}

	return ret
}

// Truncate truncates off digits from the number, without rounding.
func (d Decimal) Truncate(precision int32) Decimal {
	d.ensureInitialized()
	if precision >= 0 && -precision > d.exp {
		return d.rescale(-precision)
	}
	return d
}

func (d Decimal) string(trimTrailingZeros bool) string {
	if d.exp >= 0 {
		return d.rescale(0).value.String()
	}

	absVal := new(big.Int).Abs(d.value)
	str := absVal.String()

	var intPart, fractionalPart string

	dExpInt := int(d.exp)
	if len(str) > -dExpInt {
		intPart = str[:len(str)+dExpInt]
		fractionalPart = str[len(str)+dExpInt:]
	} else {
		intPart = "0"
		num0s := -dExpInt - len(str)
		fractionalPart = strings.Repeat("0", num0s) + str
	}

	if trimTrailingZeros {
		i := len(fractionalPart) - 1
		for ; i >= 0; i-- {
			if fractionalPart[i] != '0' {
				break
			}
		}
		fractionalPart = fractionalPart[:i+1]
	}

	number := intPart
	if len(fractionalPart) > 0 {
		number += "." + fractionalPart
	}

	if d.value.Sign() < 0 {
		return "-" + number
	}

	return number
}

func (d *Decimal) ensureInitialized() {
	if d.value == nil {
		d.value = new(big.Int)
	}
}

// RescalePair rescales two decimals to common exponential value.
func RescalePair(d1 Decimal, d2 Decimal) (Decimal, Decimal) {
	d1.ensureInitialized()
	d2.ensureInitialized()

	if d1.exp < d2.exp {
		return d1, d2.rescale(d1.exp)
	} else if d1.exp > d2.exp {
		return d1.rescale(d2.exp), d2
	}

	return d1, d2
}

// ln10 is used internally by Ln.
var ln10 = newConstApproximation("2.302585092994045684017991454684364207601101488628772976033327900967572609677352480235997205089598298341967784042286248633409525465082806756666287369098781689482907208325554680843799894826233198528393505308965377732628846163366222287698219886746543667474404243274365155048934314939391479619404400222105101714174800368808401264708068556774321622835522011480466371565912137345074785694768346361679210180644507064800027750268491674655058685693567342067058113642922455440575892572420824131469568901675894025677631135691929203337658714166023010570308963457207544037084746994016826928280848118428931484852494864487192780967627127577539702766860595249671667418348570442250719796500471495105049221477656763693866297697952211071826454973477266242570942932258279850258550978526538320760672631716430950599508780752371033310119785754733154142180842754386359177811705430982748238504564801909561029929182431823752535770975053956518769751037497088869218020518933950723853920514463419726528728696511086257149219884997874887377134568620916705849807828059751193854445009978131146915934666241071846692310107598438319191292230792503747298650929009880391941702654416816335727555703151596113564846546190897042819763365836983716328982174407366009162177850541779276367731145041782137660111010731042397832521894898817597921798666394319523936855916447118246753245630912528778330963604262982153040874560927760726641354787576616262926568298704957954913954918049209069438580790032763017941503117866862092408537949861264933479354871737451675809537088281067452440105892444976479686075120275724181874989395971643105518848195288330746699317814634930000321200327765654130472621883970596794457943468343218395304414844803701305753674262153675579814770458031413637793236291560128185336498466942261465206459942072917119370602444929358037007718981097362533224548366988505528285966192805098447175198503666680874970496982273220244823343097169111136813588418696549323714996941979687803008850408979618598756579894836445212043698216415292987811742973332588607915912510967187510929248475023930572665446276200923068791518135803477701295593646298412366497023355174586195564772461857717369368404676577047874319780573853271810933883496338813069945569399346101090745616033312247949360455361849123333063704751724871276379140924398331810164737823379692265637682071706935846394531616949411701841938119405416449466111274712819705817783293841742231409930022911502362192186723337268385688273533371925103412930705632544426611429765388301822384091026198582888433587455960453004548370789052578473166283701953392231047527564998119228742789713715713228319641003422124210082180679525276689858180956119208391760721080919923461516952599099473782780648128058792731993893453415320185969711021407542282796298237068941764740642225757212455392526179373652434440560595336591539160312524480149313234572453879524389036839236450507881731359711238145323701508413491122324390927681724749607955799151363982881058285740538000653371655553014196332241918087621018204919492651483892")

type constApproximation struct {
	exact          Decimal
	approximations []Decimal
}

func newConstApproximation(value string) constApproximation {
	parts := strings.Split(value, ".")
	coeff, fractional := parts[0], parts[1]

	coeffLen := len(coeff)
	maxPrecision := len(fractional)

	var approximations []Decimal
	for p := 1; p < maxPrecision; p *= 2 {
		r, err := NewFromString(value[:coeffLen+p])
		if err != nil {
			panic(err)
		}
		approximations = append(approximations, r)
	}

	exact, err := NewFromString(value)
	if err != nil {
		panic(err)
	}

	return constApproximation{
		exact,
		approximations,
	}
}

func (c constApproximation) withPrecision(precision int32) Decimal {
	i := 0

	if precision >= 1 {
		i++
	}

	for precision >= 16 {
		precision /= 16
		i += 4
	}

	for precision >= 2 {
		precision /= 2
		i++
	}

	if i >= len(c.approximations) {
		return c.exact
	}

	return c.approximations[i]
}
