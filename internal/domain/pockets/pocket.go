package pockets

type Pocket struct {
	PocketID int64  `json:"pocket_id"`
	Name     string `json:"name"`
	Month    string `json:"month"`
	Budget   int64  `json:"budget"`
	DateInit string `json:"date_init"`
	DateEnd  string `json:"date_end"`
}
