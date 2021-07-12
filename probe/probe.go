package probe

type Probe struct {
	Name     string   `json:"name"`
	Format   string   `json:"format"`
	Criteria Criteria `json:"criteria"`
}

type Criteria struct {
	Labels map[string]string `json:"labels"`
}
