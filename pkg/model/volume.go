package model

// Volume just represents a volume resource.
type Volume struct {
	ID      string
	Name    string
	Size    int64
	Profile string
	BusType string
	SubType string
	Device  string
}
