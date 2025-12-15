package pagination

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

// Normalize clamps limit/offset to sane bounds.
func Normalize(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
