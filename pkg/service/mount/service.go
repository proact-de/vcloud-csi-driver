package mount

// Service defines the service for the mount component.
type Service struct {
	options Options
}

// NewService simply initializes a new mount service.
func NewService(opts ...Option) *Service {
	return &Service{
		options: newOptions(opts...),
	}
}
