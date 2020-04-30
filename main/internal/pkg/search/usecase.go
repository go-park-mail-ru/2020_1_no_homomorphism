package search

import "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"

type UseCase interface {
	Search(text string, count uint) (models.SearchResult, error)
}
