package __internal_packages

// ============================================================================
// FILE: math_complete_guide.go
// TITLE: راهنمای کامل پکیج math در Go - تمام توابع با مثال
// HOW TO RUN: go run math_complete_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - پکیج math چیست؟
// ============================================================================
//
// پکیج math توابع پایه ریاضی و ثابت‌های ریاضی را ارائه می‌دهد.
// این پکیج شامل توابع برای:
// 1. توابع مثلثاتی (Trigonometric functions)
// 2. توابع هذلولوی (Hyperbolic functions)
// 3. توابع نمایی و لگاریتمی (Exponential and logarithmic)
// 4. توابع توان و ریشه (Power and root)
// 5. توابع گرد کردن (Rounding functions)
// 6. توابع آمار پایه (Min, Max, Abs, etc.)
// 7. ثابت‌های ریاضی (Pi, E, etc.)
//
// قانون طلایی:
// "توابع math روی float64 کار می‌کنند. برای عملیات صحیح از عملیات مستقیم استفاده کن.
//  همیشه boundary conditions را چک کن (NaN, Inf, etc.)."
// ============================================================================

import (
	"fmt"
	"math"
	"strings"
)

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 COMPLETE math PACKAGE GUIDE IN GO")
	fmt.Println("All mathematical functions with examples")
	fmt.Println(strings.Repeat("=", 80))

	// ============================================================================
	// بخش 1: ثابت‌های ریاضی (Mathematical Constants)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📐 SECTION 1: MATHEMATICAL CONSTANTS")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n--- 1.1 Basic Constants ---")
	fmt.Printf("  math.Pi = %.16f\n", math.Pi)
	fmt.Printf("  math.E = %.16f\n", math.E)
	fmt.Printf("  math.Phi (Golden Ratio) = %.16f\n", math.Phi)
	fmt.Printf("  math.Sqrt2 = %.16f\n", math.Sqrt2)
	fmt.Printf("  math.SqrtE = %.16f\n", math.SqrtE)
	fmt.Printf("  math.SqrtPi = %.16f\n", math.SqrtPi)
	fmt.Printf("  math.SqrtPhi = %.16f\n", math.SqrtPhi)

	fmt.Println("\n--- 1.2 Logarithmic Constants ---")
	fmt.Printf("  math.Ln2 = %.16f (ln(2))\n", math.Ln2)
	fmt.Printf("  math.Ln10 = %.16f (ln(10))\n", math.Ln10)
	fmt.Printf("  math.Log2E = %.16f (log2(e))\n", math.Log2E)
	fmt.Printf("  math.Log10E = %.16f (log10(e))\n", math.Log10E)

	fmt.Println("\n--- 1.3 Special Constants ---")
	fmt.Printf("  math.MaxFloat32 = %g\n", math.MaxFloat32)
	fmt.Printf("  math.SmallestNonzeroFloat32 = %g\n", math.SmallestNonzeroFloat32)
	fmt.Printf("  math.MaxFloat64 = %g\n", math.MaxFloat64)
	fmt.Printf("  math.SmallestNonzeroFloat64 = %g\n", math.SmallestNonzeroFloat64)
	fmt.Printf("  math.MaxInt = %d\n", math.MaxInt)
	fmt.Printf("  math.MinInt = %d\n", math.MinInt)
	fmt.Printf("  math.MaxInt8 = %d\n", math.MaxInt8)
	fmt.Printf("  math.MinInt8 = %d\n", math.MinInt8)
	fmt.Printf("  math.MaxUint8 = %d\n", math.MaxUint8)

	fmt.Println("\n--- 1.4 Special Values ---")
	fmt.Printf("  math.NaN() = %f (Not a Number)\n", math.NaN())
	fmt.Printf("  math.Inf(1) = %f (+Infinity)\n", math.Inf(1))
	fmt.Printf("  math.Inf(-1) = %f (-Infinity)\n", math.Inf(-1))

	// ============================================================================
	// بخش 2: توابع پایه (Basic Functions)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔢 SECTION 2: BASIC FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 2.1 Abs - قدر مطلق
	fmt.Println("\n--- 2.1 math.Abs ---")
	fmt.Printf("  Abs(5.5) = %.2f\n", math.Abs(5.5))
	fmt.Printf("  Abs(-5.5) = %.2f\n", math.Abs(-5.5))
	fmt.Printf("  Abs(0) = %.2f\n", math.Abs(0))

	// 2.2 Max - بزرگترین مقدار
	fmt.Println("\n--- 2.2 math.Max ---")
	fmt.Printf("  Max(10, 20) = %.2f\n", math.Max(10, 20))
	fmt.Printf("  Max(5.5, 5.5) = %.2f\n", math.Max(5.5, 5.5))
	fmt.Printf("  Max(-10, -20) = %.2f\n", math.Max(-10, -20))

	// 2.3 Min - کوچکترین مقدار
	fmt.Println("\n--- 2.3 math.Min ---")
	fmt.Printf("  Min(10, 20) = %.2f\n", math.Min(10, 20))
	fmt.Printf("  Min(5.5, 5.5) = %.2f\n", math.Min(5.5, 5.5))
	fmt.Printf("  Min(-10, -20) = %.2f\n", math.Min(-10, -20))

	// 2.4 Dim - تفاوت مثبت (max(0, x-y))
	fmt.Println("\n--- 2.4 math.Dim ---")
	fmt.Printf("  Dim(10, 5) = %.2f (max(0, 10-5))\n", math.Dim(10, 5))
	fmt.Printf("  Dim(5, 10) = %.2f (max(0, 5-10))\n", math.Dim(5, 10))

	// 2.5 Mod - باقی‌مانده (مثل % ولی برای float64)
	fmt.Println("\n--- 2.5 math.Mod ---")
	fmt.Printf("  Mod(10, 3) = %.2f\n", math.Mod(10, 3))
	fmt.Printf("  Mod(10.5, 3.2) = %.2f\n", math.Mod(10.5, 3.2))
	fmt.Printf("  Mod(-10, 3) = %.2f\n", math.Mod(-10, 3))

	// 2.6 Remainder - باقی‌مانده IEEE 754
	fmt.Println("\n--- 2.6 math.Remainder ---")
	fmt.Printf("  Remainder(10, 3) = %.2f\n", math.Remainder(10, 3))
	fmt.Printf("  Remainder(10.5, 3.2) = %.2f\n", math.Remainder(10.5, 3.2))
	fmt.Printf("  Remainder(-10, 3) = %.2f\n", math.Remainder(-10, 3))

	// 2.7 Signbit - بررسی علامت منفی
	fmt.Println("\n--- 2.7 math.Signbit ---")
	fmt.Printf("  Signbit(5.5) = %v\n", math.Signbit(5.5))
	fmt.Printf("  Signbit(-5.5) = %v\n", math.Signbit(-5.5))
	fmt.Printf("  Signbit(0) = %v\n", math.Signbit(0))

	// 2.8 Copysign - کپی علامت
	fmt.Println("\n--- 2.8 math.Copysign ---")
	fmt.Printf("  Copysign(10, -5) = %.2f (علامت منفی از -5)\n", math.Copysign(10, -5))
	fmt.Printf("  Copysign(-10, 5) = %.2f (علامت مثبت از 5)\n", math.Copysign(-10, 5))

	// ============================================================================
	// بخش 3: توابع گرد کردن (Rounding Functions)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔄 SECTION 3: ROUNDING FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	values := []float64{1.2, 1.5, 1.8, -1.2, -1.5, -1.8, 2.0}

	// 3.1 Ceil - گرد کردن به بالا
	fmt.Println("\n--- 3.1 math.Ceil ---")
	for _, v := range values {
		fmt.Printf("  Ceil(%.1f) = %.1f\n", v, math.Ceil(v))
	}

	// 3.2 Floor - گرد کردن به پایین
	fmt.Println("\n--- 3.2 math.Floor ---")
	for _, v := range values {
		fmt.Printf("  Floor(%.1f) = %.1f\n", v, math.Floor(v))
	}

	// 3.3 Round - گرد کردن به نزدیک‌ترین عدد صحیح
	fmt.Println("\n--- 3.3 math.Round ---")
	for _, v := range values {
		fmt.Printf("  Round(%.1f) = %.1f\n", v, math.Round(v))
	}

	// 3.4 RoundToEven - گرد کردن به نزدیک‌ترین عدد زوج
	fmt.Println("\n--- 3.4 math.RoundToEven ---")
	evenValues := []float64{0.5, 1.5, 2.5, 3.5, 4.5}
	for _, v := range evenValues {
		fmt.Printf("  RoundToEven(%.1f) = %.1f\n", v, math.RoundToEven(v))
	}

	// 3.5 Trunc - حذف جزء اعشاری
	fmt.Println("\n--- 3.5 math.Trunc ---")
	for _, v := range values {
		fmt.Printf("  Trunc(%.1f) = %.1f\n", v, math.Trunc(v))
	}

	// 3.6 Modf - جدا کردن جزء صحیح و اعشاری
	fmt.Println("\n--- 3.6 math.Modf ---")
	nums := []float64{3.14, -2.71, 5.0}
	for _, v := range nums {
		intPart, fracPart := math.Modf(v)
		fmt.Printf("  Modf(%.2f) = int: %.2f, frac: %.2f\n", v, intPart, fracPart)
	}

	// ============================================================================
	// بخش 4: توابع توان و ریشه (Power and Root Functions)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⚡ SECTION 4: POWER AND ROOT FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 4.1 Pow - توان
	fmt.Println("\n--- 4.1 math.Pow ---")
	fmt.Printf("  Pow(2, 3) = %.2f\n", math.Pow(2, 3))
	fmt.Printf("  Pow(2, -2) = %.4f\n", math.Pow(2, -2))
	fmt.Printf("  Pow(4, 0.5) = %.2f\n", math.Pow(4, 0.5))
	fmt.Printf("  Pow(10, 2) = %.2f\n", math.Pow(10, 2))

	// 4.2 Pow10 - توان 10
	fmt.Println("\n--- 4.2 math.Pow10 ---")
	for i := -2; i <= 3; i++ {
		fmt.Printf("  Pow10(%d) = %.2f\n", i, math.Pow10(i))
	}

	// 4.3 Sqrt - جذر مربع
	fmt.Println("\n--- 4.3 math.Sqrt ---")
	squares := []float64{0, 1, 2, 4, 9, 16, 25}
	for _, v := range squares {
		fmt.Printf("  Sqrt(%.0f) = %.4f\n", v, math.Sqrt(v))
	}

	// 4.4 Cbrt - جذر مکعب
	fmt.Println("\n--- 4.4 math.Cbrt ---")
	cubes := []float64{0, 1, 8, 27, 64, 125}
	for _, v := range cubes {
		fmt.Printf("  Cbrt(%.0f) = %.4f\n", v, math.Cbrt(v))
	}

	// 4.5 Hypot - sqrt(a*a + b*b) (مفید برای وتر مثلث قائم‌الزاویه)
	fmt.Println("\n--- 4.5 math.Hypot ---")
	fmt.Printf("  Hypot(3, 4) = %.2f (مثلث 3-4-5)\n", math.Hypot(3, 4))
	fmt.Printf("  Hypot(5, 12) = %.2f\n", math.Hypot(5, 12))

	// ============================================================================
	// بخش 5: توابع نمایی و لگاریتمی (Exponential and Logarithmic)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📈 SECTION 5: EXPONENTIAL AND LOGARITHMIC FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 5.1 Exp - e^x
	fmt.Println("\n--- 5.1 math.Exp ---")
	expVals := []float64{0, 1, -1, 2, -2}
	for _, v := range expVals {
		fmt.Printf("  Exp(%.1f) = %.6f\n", v, math.Exp(v))
	}
	fmt.Printf("  Exp(1) = %.6f (e)\n", math.Exp(1))

	// 5.2 Exp2 - 2^x
	fmt.Println("\n--- 5.2 math.Exp2 ---")
	for _, v := range expVals {
		fmt.Printf("  Exp2(%.1f) = %.6f\n", v, math.Exp2(v))
	}

	// 5.3 Expm1 - e^x - 1 (دقیق‌تر برای x نزدیک صفر)
	fmt.Println("\n--- 5.3 math.Expm1 ---")
	smallVals := []float64{0.0001, 0.001, 0.01, 0.1}
	for _, v := range smallVals {
		fmt.Printf("  Expm1(%.4f) = %.6f\n", v, math.Expm1(v))
	}

	// 5.4 Log - لگاریتم طبیعی (ln)
	fmt.Println("\n--- 5.4 math.Log ---")
	logVals := []float64{1, math.E, 10, 100}
	for _, v := range logVals {
		fmt.Printf("  Log(%.2f) = %.6f\n", v, math.Log(v))
	}

	// 5.5 Log10 - لگاریتم پایه 10
	fmt.Println("\n--- 5.5 math.Log10 ---")
	for _, v := range logVals {
		fmt.Printf("  Log10(%.2f) = %.6f\n", v, math.Log10(v))
	}

	// 5.6 Log2 - لگاریتم پایه 2
	fmt.Println("\n--- 5.6 math.Log2 ---")
	for _, v := range logVals {
		fmt.Printf("  Log2(%.2f) = %.6f\n", v, math.Log2(v))
	}

	// 5.7 Log1p - ln(1+x) (دقیق‌تر برای x نزدیک صفر)
	fmt.Println("\n--- 5.7 math.Log1p ---")
	for _, v := range smallVals {
		fmt.Printf("  Log1p(%.4f) = %.6f\n", v, math.Log1p(v))
	}

	// 5.8 Logb - لگاریتم پایه 2 از قدر مطلق (باینری)
	fmt.Println("\n--- 5.8 math.Logb ---")
	logbVals := []float64{2, 4, 8, 16, 32, 0.5, 0.25}
	for _, v := range logbVals {
		fmt.Printf("  Logb(%.2f) = %.2f\n", v, math.Logb(v))
	}

	// 5.9 Frexp - جدا کردن fraction و exponent (x = frac × 2^exp)
	fmt.Println("\n--- 5.9 math.Frexp ---")
	frexpVals := []float64{12.5, 0.375, 1.0}
	for _, v := range frexpVals {
		frac, exp := math.Frexp(v)
		fmt.Printf("  Frexp(%.2f) = frac: %.4f, exp: %d (%.2f = %.4f × 2^%d)\n", v, frac, exp, v, frac, exp)
	}

	// 5.10 Ldexp - معکوس Frexp (frac × 2^exp)
	fmt.Println("\n--- 5.10 math.Ldexp ---")
	fmt.Printf("  Ldexp(0.78125, 4) = %.2f\n", math.Ldexp(0.78125, 4))

	// ============================================================================
	// بخش 6: توابع مثلثاتی (Trigonometric Functions)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📐 SECTION 6: TRIGONOMETRIC FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// زوایای مختلف (به رادیان)
	angles := []float64{
		0,               // 0°
		math.Pi / 6,     // 30°
		math.Pi / 4,     // 45°
		math.Pi / 3,     // 60°
		math.Pi / 2,     // 90°
		math.Pi,         // 180°
		3 * math.Pi / 2, // 270°
		2 * math.Pi,     // 360°
	}

	// 6.1 Sin - سینوس
	fmt.Println("\n--- 6.1 math.Sin ---")
	for _, angle := range angles {
		fmt.Printf("  Sin(%.2f rad) = %.4f\n", angle, math.Sin(angle))
	}

	// 6.2 Cos - کسینوس
	fmt.Println("\n--- 6.2 math.Cos ---")
	for _, angle := range angles {
		fmt.Printf("  Cos(%.2f rad) = %.4f\n", angle, math.Cos(angle))
	}

	// 6.3 Tan - تانژانت
	fmt.Println("\n--- 6.3 math.Tan ---")
	for _, angle := range angles[:5] { // فقط زوایای معتبر
		fmt.Printf("  Tan(%.2f rad) = %.4f\n", angle, math.Tan(angle))
	}

	// 6.4 SinCos - سینوس و کسینوس همزمان
	fmt.Println("\n--- 6.4 math.Sincos ---")
	for _, angle := range angles {
		sin, cos := math.Sincos(angle)
		fmt.Printf("  Sincos(%.2f rad) = sin: %.4f, cos: %.4f\n", angle, sin, cos)
	}

	// ============================================================================
	// بخش 7: توابع معکوس مثلثاتی (Inverse Trigonometric)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔄 SECTION 7: INVERSE TRIGONOMETRIC FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// مقادیر سینوس
	sinVals := []float64{-1, -0.5, 0, 0.5, 1}

	// 7.1 Asin - آرک سینوس (معکوس سینوس)
	fmt.Println("\n--- 7.1 math.Asin ---")
	for _, v := range sinVals {
		fmt.Printf("  Asin(%.2f) = %.4f rad\n", v, math.Asin(v))
	}

	// 7.2 Acos - آرک کسینوس (معکوس کسینوس)
	fmt.Println("\n--- 7.2 math.Acos ---")
	for _, v := range sinVals {
		fmt.Printf("  Acos(%.2f) = %.4f rad\n", v, math.Acos(v))
	}

	// 7.3 Atan - آرک تانژانت (معکوس تانژانت)
	fmt.Println("\n--- 7.3 math.Atan ---")
	tanVals := []float64{-10, -1, 0, 1, 10}
	for _, v := range tanVals {
		fmt.Printf("  Atan(%.2f) = %.4f rad\n", v, math.Atan(v))
	}

	// 7.4 Atan2 - آرک تانژانت دو متغیره (y/x)
	fmt.Println("\n--- 7.4 math.Atan2 ---")
	points := [][2]float64{
		{1, 1},   // 45°
		{-1, 1},  // 135°
		{-1, -1}, // -135° (225°)
		{1, -1},  // -45° (315°)
	}
	for _, p := range points {
		angle := math.Atan2(p[1], p[0])
		fmt.Printf("  Atan2(%.2f, %.2f) = %.4f rad\n", p[1], p[0], angle)
	}

	// ============================================================================
	// بخش 8: توابع هذلولوی (Hyperbolic Functions)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📈 SECTION 8: HYPERBOLIC FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	hyperVals := []float64{-2, -1, 0, 1, 2}

	// 8.1 Sinh - سینوس هذلولوی
	fmt.Println("\n--- 8.1 math.Sinh ---")
	for _, v := range hyperVals {
		fmt.Printf("  Sinh(%.1f) = %.4f\n", v, math.Sinh(v))
	}

	// 8.2 Cosh - کسینوس هذلولوی
	fmt.Println("\n--- 8.2 math.Cosh ---")
	for _, v := range hyperVals {
		fmt.Printf("  Cosh(%.1f) = %.4f\n", v, math.Cosh(v))
	}

	// 8.3 Tanh - تانژانت هذلولوی
	fmt.Println("\n--- 8.3 math.Tanh ---")
	for _, v := range hyperVals {
		fmt.Printf("  Tanh(%.1f) = %.4f\n", v, math.Tanh(v))
	}

	// 8.4 Asinh - آرک سینوس هذلولوی
	fmt.Println("\n--- 8.4 math.Asinh ---")
	for _, v := range hyperVals {
		fmt.Printf("  Asinh(%.1f) = %.4f\n", v, math.Asinh(v))
	}

	// 8.5 Acosh - آرک کسینوس هذلولوی
	fmt.Println("\n--- 8.5 math.Acosh ---")
	for _, v := range []float64{1, 2, 3, 5} {
		fmt.Printf("  Acosh(%.1f) = %.4f\n", v, math.Acosh(v))
	}

	// 8.6 Atanh - آرک تانژانت هذلولوی
	fmt.Println("\n--- 8.6 math.Atanh ---")
	for _, v := range []float64{-0.9, -0.5, 0, 0.5, 0.9} {
		fmt.Printf("  Atanh(%.1f) = %.4f\n", v, math.Atanh(v))
	}

	// ============================================================================
	// بخش 9: توابع ویژه (Special Functions)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⚡ SECTION 9: SPECIAL FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 9.1 Gamma - تابع گاما (تعمیم فاکتوریل برای اعداد حقیقی)
	fmt.Println("\n--- 9.1 math.Gamma ---")
	gammaVals := []float64{0.5, 1, 2, 3, 4, 5}
	for _, v := range gammaVals {
		fmt.Printf("  Gamma(%.1f) = %.6f\n", v, math.Gamma(v))
	}
	fmt.Printf("  Gamma(5) = %.0f (4! = 24)\n", math.Gamma(5))

	// 9.2 LogGamma - لگاریتم طبیعی تابع گاما
	fmt.Println("\n--- 9.2 math.LogGamma ---")
	for _, v := range gammaVals {
		fmt.Printf("  LogGamma(%.1f) = %.6f\n", v, math.LogGamma(v))
	}

	// 9.3 J0 - تابع بسل نوع اول مرتبه 0
	fmt.Println("\n--- 9.3 math.J0 ---")
	besselVals := []float64{0, 1, 2, 3, 4, 5}
	for _, v := range besselVals {
		fmt.Printf("  J0(%.1f) = %.6f\n", v, math.J0(v))
	}

	// 9.4 J1 - تابع بسل نوع اول مرتبه 1
	fmt.Println("\n--- 9.4 math.J1 ---")
	for _, v := range besselVals {
		fmt.Printf("  J1(%.1f) = %.6f\n", v, math.J1(v))
	}

	// 9.5 Y0 - تابع بسل نوع دوم مرتبه 0
	fmt.Println("\n--- 9.5 math.Y0 ---")
	for _, v := range []float64{1, 2, 3, 4, 5} {
		fmt.Printf("  Y0(%.1f) = %.6f\n", v, math.Y0(v))
	}

	// 9.6 Erf - تابع خطا (Error function)
	fmt.Println("\n--- 9.6 math.Erf ---")
	erfVals := []float64{-2, -1, -0.5, 0, 0.5, 1, 2}
	for _, v := range erfVals {
		fmt.Printf("  Erf(%.1f) = %.6f\n", v, math.Erf(v))
	}

	// 9.7 Erfc - تابع خطای مکمل (1 - Erf(x))
	fmt.Println("\n--- 9.7 math.Erfc ---")
	for _, v := range erfVals {
		fmt.Printf("  Erfc(%.1f) = %.6f\n", v, math.Erfc(v))
	}

	// ============================================================================
	// بخش 10: توابع کسری (Fraction Functions)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔢 SECTION 10: FRACTION AND INTEGER FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 10.1 IsNaN - بررسی Not a Number
	fmt.Println("\n--- 10.1 math.IsNaN ---")
	nan := math.NaN()
	fmt.Printf("  IsNaN(5.5) = %v\n", math.IsNaN(5.5))
	fmt.Printf("  IsNaN(NaN) = %v\n", math.IsNaN(nan))

	// 10.2 IsInf - بررسی بی‌نهایت
	fmt.Println("\n--- 10.2 math.IsInf ---")
	posInf := math.Inf(1)
	negInf := math.Inf(-1)
	fmt.Printf("  IsInf(5.5, 0) = %v (any infinity)\n", math.IsInf(5.5, 0))
	fmt.Printf("  IsInf(posInf, 1) = %v (+inf)\n", math.IsInf(posInf, 1))
	fmt.Printf("  IsInf(negInf, -1) = %v (-inf)\n", math.IsInf(negInf, -1))

	// 10.3 Inf - ایجاد بی‌نهایت
	fmt.Println("\n--- 10.3 math.Inf ---")
	fmt.Printf("  Inf(1) = %v (+∞)\n", math.Inf(1))
	fmt.Printf("  Inf(-1) = %v (-∞)\n", math.Inf(-1))

	// 10.4 NaN - ایجاد Not a Number
	fmt.Println("\n--- 10.4 math.NaN ---")
	fmt.Printf("  NaN() = %v\n", math.NaN())

	// ============================================================================
	// بخش 11: توابع سقف و کف (Ceil, Floor, etc.)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🏢 SECTION 11: CEIL, FLOOR, AND RELATED")
	fmt.Println(strings.Repeat("=", 80))

	testVals := []float64{2.1, 2.5, 2.9, -2.1, -2.5, -2.9}

	fmt.Println("\n--- 11.1 Comparison of Rounding Functions ---")
	fmt.Println("  Value   | Ceil | Floor | Round | Trunc | RoundToEven")
	fmt.Println("  --------|------|-------|-------|-------|------------")
	for _, v := range testVals {
		fmt.Printf("  %6.1f | %4.1f | %5.1f | %5.1f | %5.1f | %10.1f\n",
			v, math.Ceil(v), math.Floor(v), math.Round(v), math.Trunc(v), math.RoundToEven(v))
	}

	// ============================================================================
	// بخش 12: کاربردهای عملی
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💡 SECTION 12: PRACTICAL EXAMPLES")
	fmt.Println(strings.Repeat("=", 80))

	// 12.1 محاسبه مساحت دایره
	fmt.Println("\n--- 12.1 Circle Area ---")
	radius := 5.0
	area := math.Pi * math.Pow(radius, 2)
	fmt.Printf("  Circle with radius %.1f has area: %.2f\n", radius, area)

	// 12.2 محاسبه وتر مثلث قائم‌الزاویه
	fmt.Println("\n--- 12.2 Hypotenuse Calculation ---")
	a, b := 3.0, 4.0
	c := math.Hypot(a, b)
	fmt.Printf("  Right triangle with sides %.0f and %.0f has hypotenuse: %.2f\n", a, b, c)

	// 12.3 تبدیل درجه به رادیان و بالعکس
	fmt.Println("\n--- 12.3 Degree to Radian Conversion ---")
	deg := 180.0
	rad := deg * math.Pi / 180
	fmt.Printf("  %.0f degrees = %.4f radians\n", deg, rad)
	fmt.Printf("  %.4f radians = %.0f degrees\n", rad, rad*180/math.Pi)

	// 12.4 محاسبه فاصله اقلیدسی بین دو نقطه
	fmt.Println("\n--- 12.4 Euclidean Distance ---")
	x1, y1 := 0.0, 0.0
	x2, y2 := 3.0, 4.0
	distance := math.Hypot(x2-x1, y2-y1)
	fmt.Printf("  Distance between (%.0f,%.0f) and (%.0f,%.0f) = %.2f\n", x1, y1, x2, y2, distance)

	// 12.5 محاسبه زاویه بین دو بردار
	fmt.Println("\n--- 12.5 Angle Between Vectors ---")
	v1x, v1y := 1.0, 0.0
	v2x, v2y := 0.0, 1.0
	dot := v1x*v2x + v1y*v2y
	mag1 := math.Hypot(v1x, v1y)
	mag2 := math.Hypot(v2x, v2y)
	cosTheta := dot / (mag1 * mag2)
	angle := math.Acos(cosTheta) * 180 / math.Pi
	fmt.Printf("  Angle between vectors (%.0f,%.0f) and (%.0f,%.0f) = %.2f degrees\n", v1x, v1y, v2x, v2y, angle)

	// 12.6 نرمال کردن زاویه به بازه [0, 2π]
	fmt.Println("\n--- 12.6 Normalize Angle ---")
	rawAngle := 450.0 // degrees
	radians := rawAngle * math.Pi / 180
	normalized := math.Mod(radians, 2*math.Pi)
	fmt.Printf("  Angle %.0f° normalized to %.2f°\n", rawAngle, normalized*180/math.Pi)

	// ============================================================================
	// بخش 13: اشتباهات رایج
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("❌ SECTION 13: COMMON MISTAKES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n❌ Mistake 1: Not checking for NaN/Inf")
	fmt.Println("   result := math.Sqrt(-1)  // returns NaN")
	fmt.Println("   ✅ Check math.IsNaN(result) before using")

	fmt.Println("\n❌ Mistake 2: Comparing floats with ==")
	fmt.Println("   if a == b { ... }  // floating point precision issues")
	fmt.Println("   ✅ Use tolerance: math.Abs(a-b) < 1e-9")

	fmt.Println("\n❌ Mistake 3: Passing degrees to trig functions")
	fmt.Println("   math.Sin(90)  // expects radians, not degrees")
	fmt.Println("   ✅ Convert: math.Sin(90 * math.Pi / 180)")

	fmt.Println("\n❌ Mistake 4: Using math.Mod for negative numbers")
	fmt.Println("   math.Mod(-5, 2) = -1  // not 1")
	fmt.Println("   ✅ For positive remainder: math.Mod(math.Mod(x, m)+m, m)")

	fmt.Println("\n❌ Mistake 5: Assuming integer precision")
	fmt.Println("   math.Pow(10, 2) = 100.00000000000001")
	fmt.Println("   ✅ Round: math.Round(math.Pow(10, 2))")

	fmt.Println("\n❌ Mistake 6: Overflow in integer operations")
	fmt.Println("   int64(math.Pow(2, 60))  // may lose precision")
	fmt.Println("   ✅ Use bit shifting: 1 << 60")

	// ============================================================================
	// بخش 14: جمع‌بندی و جدول مرجع
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📚 SECTION 14: QUICK REFERENCE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ CATEGORY              │ FUNCTIONS                              │")
	fmt.Println("├───────────────────────┼────────────────────────────────────────┤")
	fmt.Println("│ Basic                 │ Abs, Max, Min, Dim, Mod, Remainder,    │")
	fmt.Println("│                       │ Signbit, Copysign                      │")
	fmt.Println("│ Rounding              │ Ceil, Floor, Round, RoundToEven, Trunc,│")
	fmt.Println("│                       │ Modf                                   │")
	fmt.Println("│ Power/Root            │ Pow, Pow10, Sqrt, Cbrt, Hypot          │")
	fmt.Println("│ Exponential/Log       │ Exp, Exp2, Expm1, Log, Log10, Log2,    │")
	fmt.Println("│                       │ Log1p, Logb, Frexp, Ldexp              │")
	fmt.Println("│ Trigonometric         │ Sin, Cos, Tan, Sincos, Asin, Acos,     │")
	fmt.Println("│                       │ Atan, Atan2                            │")
	fmt.Println("│ Hyperbolic            │ Sinh, Cosh, Tanh, Asinh, Acosh, Atanh  │")
	fmt.Println("│ Special               │ Gamma, LogGamma, J0, J1, Y0, Erf, Erfc │")
	fmt.Println("│ Constants             │ Pi, E, Phi, Sqrt2, Ln2, Ln10, MaxInt,  │")
	fmt.Println("│                       │ MaxFloat64, etc.                       │")
	fmt.Println("└───────────────────────┴────────────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. Always check for NaN/Inf when using math functions")
	fmt.Println("  2. Use tolerance for floating point comparisons")
	fmt.Println("  3. Convert degrees to radians for trig functions")
	fmt.Println("  4. Use math.Mod correctly for negative numbers")
	fmt.Println("  5. Prefer Hypot over manual sqrt for distance")
	fmt.Println("  6. Use Expm1/Log1p for values near zero")
	fmt.Println("  7. Use Modf to separate integer and fractional parts")
	fmt.Println("  8. Use IsNaN/IsInf for validation")
	fmt.Println("  9. Remember floating point is approximate")
	fmt.Println("  10. Use math.Round for rounding to nearest integer")

	fmt.Println("\n🎯 COMMON CONSTANTS:")
	fmt.Printf("  π (pi)           = %.16f\n", math.Pi)
	fmt.Printf("  e                = %.16f\n", math.E)
	fmt.Printf("  φ (phi)          = %.16f\n", math.Phi)
	fmt.Printf("  √2               = %.16f\n", math.Sqrt2)
	fmt.Printf("  ln(2)            = %.16f\n", math.Ln2)
	fmt.Printf("  ln(10)           = %.16f\n", math.Ln10)
	fmt.Printf("  log2(e)          = %.16f\n", math.Log2E)
	fmt.Printf("  log10(e)         = %.16f\n", math.Log10E)
}

// تابع کمکی برای تکرار رشته
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
