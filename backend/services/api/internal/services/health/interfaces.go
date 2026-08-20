package health

type Service interface {
	HealthCheck() map[string]string
}

type healthServiceImpl struct{}

func NewService() Service {
	return &healthServiceImpl{}
}
