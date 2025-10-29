package types

type LoginUserRequest struct {
	SessionID string `json:"session_id"`
	Provider string `json:"provider"`
}