package models

// Insight represents an AI-generated insight from the Starscope API.
type Insight struct {
	ID               int         `json:"id"`
	Type             string      `json:"type"`
	Title            string      `json:"title"`
	Body             string      `json:"body"`
	Severity         string      `json:"severity"`
	Confidence       float64     `json:"confidence"`
	SuggestedActions []string    `json:"suggested_actions"`
	PeriodStart      string      `json:"period_start"`
	PeriodEnd        string      `json:"period_end"`
	GeneratedAt      string      `json:"generated_at"`
	IsRead           bool        `json:"is_read"`
	Connection       *Connection `json:"connection,omitempty"`
}
