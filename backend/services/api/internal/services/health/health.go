package health

func (s *healthServiceImpl) HealthCheck() map[string]string {
	return map[string]string{"status": "ok"}
}
