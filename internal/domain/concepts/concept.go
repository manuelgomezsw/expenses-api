package concepts

type Concept struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Value      int64  `json:"value"`
	PocketID   int    `json:"pocket_id"`
	PocketName string `json:"pocket_name"`
	Payed      bool   `json:"payed"`
	UpdatedAt  string `json:"updated_at"`
	PaymentDay int16  `json:"payment_day"`
}
