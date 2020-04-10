package base_api

type OperationType int

type Operation struct {
	Operate                func(arguments []float64) (result float64, asString string, err error) `json:"operate"`
	ExpectedArgumentLength int                                                                    `json:"expectedArgumentLength"`
}

type OperationResponse struct {
	Message   string    `json:"message"`
	Result    float64   `json:"result"`
	Arguments []float64 `json:"arguments"`
	AsString  string    `json:"asString"`
}
