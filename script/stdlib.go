package script

import (
	"math"
	"math/rand"
)

// Loads standard functions into the world.
func (w *World) LoadStdlib() {

	// literals
	w.declare("true", boolLit(true))
	w.declare("false", boolLit(false))

	// math
	w.PureFunc("abs", math.Abs)
	w.PureFunc("acos", math.Acos)
	w.PureFunc("acosh", math.Acosh)
	w.PureFunc("asin", math.Asin)
	w.PureFunc("asinh", math.Asinh)
	w.PureFunc("atan", math.Atan)
	w.PureFunc("atanh", math.Atanh)
	w.PureFunc("cbrt", math.Cbrt)
	w.PureFunc("ceil", math.Ceil)
	w.PureFunc("cos", math.Cos)
	w.PureFunc("cosh", math.Cosh)
	w.PureFunc("erf", math.Erf)
	w.PureFunc("erfc", math.Erfc)
	w.PureFunc("exp", math.Exp)
	w.PureFunc("exp2", math.Exp2)
	w.PureFunc("expm1", math.Expm1)
	w.PureFunc("floor", math.Floor)
	w.PureFunc("gamma", math.Gamma)
	w.PureFunc("j0", math.J0)
	w.PureFunc("j1", math.J1)
	w.PureFunc("log", math.Log)
	w.PureFunc("log10", math.Log10)
	w.PureFunc("log1p", math.Log1p)
	w.PureFunc("log2", math.Log2)
	w.PureFunc("logb", math.Logb)
	w.PureFunc("sin", math.Sin)
	w.PureFunc("sinh", math.Sinh)
	w.PureFunc("sqrt", math.Sqrt)
	w.PureFunc("tan", math.Tan)
	w.PureFunc("tanh", math.Tanh)
	w.PureFunc("trunc", math.Trunc)
	w.PureFunc("y0", math.Y0)
	w.PureFunc("y1", math.Y1)
	w.PureFunc("ilogb", math.Ilogb)
	w.PureFunc("pow10", math.Pow10)
	w.PureFunc("atan2", math.Atan2)
	w.PureFunc("hypot", math.Hypot)
	w.PureFunc("remainder", math.Remainder)
	w.PureFunc("max", math.Max)
	w.PureFunc("min", math.Min)
	w.PureFunc("mod", math.Mod)
	w.PureFunc("pow", math.Pow)
	w.PureFunc("yn", math.Yn)
	w.PureFunc("jn", math.Jn)
	w.PureFunc("ldexp", math.Ldexp)
	w.PureFunc("isInf", math.IsInf)
	w.PureFunc("isNaN", math.IsNaN)
	w.PureFunc("heaviside", heaviside)
	w.declare("pi", floatLit(math.Pi))
	w.declare("inf", floatLit(math.Inf(1)))

	w.PureFunc("sinc", sinc)

	w.Func("randSeed", intseed, "Sets the random number seed")
	w.Func("rand", rand.Float64, "Random number between 0 and 1")
	w.Func("randExp", rand.ExpFloat64, "Exponentially distributed random number between 0 and +inf, mean=1")
	w.Func("randNorm", rand.NormFloat64, "Standard normal random number")
	w.Func("randInt", randInt, "Random non-negative integer")
}

// script does not know int64
func intseed(seed int) {
	rand.Seed(int64(seed))
}

func randInt(upper int) int {
	return rand.Int() % upper
}

func heaviside(x float64) float64 {
	switch {
	default:
		return 1
	case x == 0:
		return 0.5
	case x < 0:
		return 0
	}
}

func sinc(x float64) float64 {
	if x == 0 {
		return 1
	} else {
		return math.Sin(x) / x
	}
}
