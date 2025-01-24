package expenses

type Expense struct {
	ExpenseID       int    `json:"expense_id"`
	Name            string `json:"name"`
	Value           int64  `json:"value"`
	CycleID         int    `json:"cycle_id"`
	CycleName       string `json:"cycle_name"`
	PocketID        int    `json:"pocket_id"`
	PocketName      string `json:"pocket_name"`
	PaymentTypeID   int16  `json:"payment_type_id"`
	PaymentTypeName string `json:"payment_type_name"`
	CreatedAt       string `json:"created_at"`
}
