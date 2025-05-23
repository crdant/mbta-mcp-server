package mbta

import (
	"math"
)

// isCloseEnough checks if two float values are close enough to be considered equal
func isCloseEnough(a, b float64) bool {
	const epsilon = 0.0001
	return math.Abs(a-b) < epsilon
}
