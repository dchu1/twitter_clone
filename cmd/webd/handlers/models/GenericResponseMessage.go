package handlermodels

type GenericResponseMessage struct {
	Status  int               `json:",omitempty"`
	Message string            `json:",omitempty"`
	Body    map[string]string `json:",omitempty"`
}
