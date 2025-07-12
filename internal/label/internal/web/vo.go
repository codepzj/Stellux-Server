package web

type LabelVO struct {
	ID        string `json:"id"`
	LabelType string `json:"label_type"`
	Name      string `json:"name"`
}

type LabelWithCountVO struct {
	ID        string `json:"id"`
	LabelType string `json:"label_type"`
	Name      string `json:"name"`
	Count     int    `json:"count"`
}