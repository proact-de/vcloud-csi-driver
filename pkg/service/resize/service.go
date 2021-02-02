package resize

// Service defines the service for the resize component.
type Service struct {
	options Options
}

// NewService simply initializes a new resize service.
func NewService(opts ...Option) *Service {
	return &Service{
		options: newOptions(opts...),
	}
}
