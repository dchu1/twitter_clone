package handlermodels

type GenericResponseMessage struct {
	Status  int
	Message string
	Body    map[string]string
}
