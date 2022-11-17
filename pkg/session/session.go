package session

type ServiceSession struct {
}

func CreateSession() *ServiceSession {
	return &ServiceSession{}
}
