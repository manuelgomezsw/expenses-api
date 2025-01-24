package cycles

type Cycle struct {
	CycleID    int    `json:"cycle_id"`
	PocketID   int    `json:"pocket_id"`
	PocketName string `json:"pocket_name"`
	Name       string `json:"name"`
	Budget     int64  `json:"budget"`
	DateInit   string `json:"date_init"`
	DateEnd    string `json:"date_end"`
	Status     bool   `json:"status"`
	CreatedAt  string `json:"created_at"`
}
