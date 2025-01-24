package pockets

type Pocket struct {
	PocketID  int64  `json:"pocket_id"`
	Name      string `json:"name"`
	Status    bool   `json:"status"`
	CreatedAt string `json:"created_at"`
}
