package streamlabs

type Donation struct {
	Amount   int
	Currency string
}

type ReSub struct {
	Months int
}

type Follow struct {
}

type ev struct {
	For     string `json:"for"`
	Type    string `json:"type"`
	Message []struct {
		Amount   interface{} `json:"amount"`
		Currency string      `json:"currency"`
	} `json:"message"`
}
