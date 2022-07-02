package models

type RelatedAnime struct {
	Node                  Node   `json:"node"`
	RelationType          string `json:"relation_type"`
	RelationTypeFormatted string `json:"relation_type_formatted"`
}
