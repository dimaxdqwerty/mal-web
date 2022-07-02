package models

type Recommendations struct {
	Node               Node `json:"node"`
	NumRecommendations int  `json:"num_recommendations"`
}
