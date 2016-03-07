package tempconv

// CToF converts a Celsius temperature to Fahrenheit.
func CToF(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) }

// FToC converts a Fahrenheit temperature to Celsius.
func FToC(f Fahrenheit) Celsius { return Celsius((f - 32) * 5 / 9) }

// KtoC converts a Kelvin temperature to Celsius
func KtoC(k Kelvin) Celsius { return AbsoluteZeroC + Celsius(k / 1000) }

// CtoK converts a Celsius temperature to a Kelvin
func CtoK(c Celsius) Kelvin { return Kelvin(c - AbsoluteZeroC) * 1000 }

// CtoK converts a Celsius temperature to Kelvin
