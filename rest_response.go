package httpx

type RestResponse struct {
	Success     bool   `json:"success"`
	Description string `json:"description"`
	Payload     any    `json:"payload"`
}
