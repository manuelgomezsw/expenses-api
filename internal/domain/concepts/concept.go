package concepts

type Concept struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Value     int64  `json:"value"`
	PocketID  int    `json:"pocket_id"`
	Payed     bool   `json:"payed"`
	UpdatedAt string `json:"updated_at"`
}
