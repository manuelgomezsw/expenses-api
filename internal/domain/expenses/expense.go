package expenses

type Expense struct {
	ExpenseID     int    `json:"expense_id"`
	Name          string `json:"name"`
	Value         int64  `json:"value"`
	PocketID      int    `json:"pocket_id"`
	PaymentTypeID int16  `json:"payment_type_id"`
}
