package request

type Echo struct {
	Data string `json:"data" binding:"required"`
}
