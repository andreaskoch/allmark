package repository

import "log"

// LineRange contains a Start- and a End line number.
type LineRange struct {
	Start int
	End   int
}

// NewLineRange returns a new LineRange
// with the given start and end.
func NewLineRange(start int, end int) LineRange {
	if start < 0 || end < 0 || (start > end) {
		log.Panicf("Invalid start and end values for a LineRange. Start: %v, End: %v", start, end)
	}

	return LineRange{
		Start: start,
		End:   end,
	}
}
