package utils

import "mal/models"

func ContainsOneOfGenres(genres []models.Genres, formGenres []string) bool {
	for _, genre := range genres {
		for _, formGenre := range formGenres {
			if genre.Name == formGenre {
				return true
			}
		}
	}
	return false
}
