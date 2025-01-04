package paymentstype

type PaymentType struct {
	PaymentTypeID int16  `json:"payment_type_id"`
	Name          string `json:"name"`
}
