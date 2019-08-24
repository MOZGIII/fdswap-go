package fdswap

// Restorer provides a way to restore fd to it's original underlying target.
type Restorer interface {
	Restore() error
}
