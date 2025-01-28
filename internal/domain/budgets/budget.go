package budgets

type Budget struct {
	CycleID    int   `json:"cycle_id"`
	Spent      int64 `json:"spent"`
	Budget     int64 `json:"budget"`
	SpentRatio int   `json:"spent_ratio"`
}
