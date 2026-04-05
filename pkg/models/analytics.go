package models

// AnalyticsOverview represents the analytics overview response.
type AnalyticsOverview struct {
	TotalReviews     int            `json:"total_reviews"`
	AverageRating    float64        `json:"average_rating"`
	ResponseRate     float64        `json:"response_rate"`
	SentimentBreakdown map[string]int `json:"sentiment_breakdown"`
	RatingDistribution map[string]int `json:"rating_distribution"`
	PlatformBreakdown  map[string]int `json:"platform_breakdown"`
}

// AnalyticsMetrics represents the analytics metrics response.
type AnalyticsMetrics struct {
	Period     int              `json:"period"`
	DataPoints []MetricDataPoint `json:"data_points"`
}

// MetricDataPoint represents a single data point in analytics metrics.
type MetricDataPoint struct {
	Date         string  `json:"date"`
	ReviewCount  int     `json:"review_count"`
	AverageRating float64 `json:"average_rating"`
	Sentiment    float64 `json:"sentiment"`
}
