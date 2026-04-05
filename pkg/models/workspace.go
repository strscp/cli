package models

// Workspace represents a workspace from the Starscope API.
type Workspace struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	Slug            string  `json:"slug"`
	PlanTier        string  `json:"plan_tier"`
	ConnectionCount int     `json:"connection_count"`
	ConnectionLimit int     `json:"connection_limit"`
	TrialEndsAt     *string `json:"trial_ends_at"`
}
