package models

// Review represents a review from the Starscope API.
type Review struct {
	ID             int         `json:"id"`
	ExternalID     string      `json:"external_id"`
	Platform       string      `json:"platform"`
	AuthorName     string      `json:"author_name"`
	Rating         int         `json:"rating"`
	Title          *string     `json:"title"`
	Body           *string     `json:"body"`
	Language       *string     `json:"language"`
	PublishedAt    string      `json:"published_at"`
	ResponseText   *string     `json:"response_text"`
	RespondedAt    *string     `json:"responded_at"`
	Sentiment      *string     `json:"sentiment"`
	SentimentScore *FlexFloat  `json:"sentiment_score"`
	Summary        *string     `json:"summary"`
	Topics         []string    `json:"topics"`
	Keywords       []string    `json:"keywords"`
	CustomerIntent *string     `json:"customer_intent"`
	QualityScore   *FlexFloat  `json:"quality_score"`
	AnalyzedAt     *string     `json:"analyzed_at"`
	Connection     *Connection `json:"connection,omitempty"`
}
