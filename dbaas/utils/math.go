package utils

func Min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func MinInt64(a, b int64) int64 {
	if a <= b {
		return a
	}
	return b
}

func MaxInt64(a, b int64) int64 {
	if a >= b {
		return a
	}
	return b
}
