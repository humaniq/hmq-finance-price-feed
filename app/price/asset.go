package price

type Asset struct {
	Name    string             `json:"name"`
	Prices  map[string]Value   `json:"prices_old"`
	History map[string]History `json:"history"`
}

func NewAsset(name string) *Asset {
	return &Asset{
		Name:    name,
		Prices:  make(map[string]Value),
		History: make(map[string]History),
	}
}
