package models

// Connection represents a review platform connection from the Starscope API.
type Connection struct {
	ID           int     `json:"id"`
	Platform     string  `json:"platform"`
	Name         string  `json:"name"`
	ExternalID   string  `json:"external_id"`
	ExternalURL  *string `json:"external_url"`
	IsActive     bool    `json:"is_active"`
	LastSyncedAt *string `json:"last_synced_at"`
	ReviewCount  int     `json:"review_count"`
	LogoURL      *string `json:"logo_url"`
	Rating       *float64 `json:"rating"`
	BrandColor   *string `json:"brand_color"`
}
