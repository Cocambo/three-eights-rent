package service

import (
	"math"
	"sort"

	"car-service/internal/repository"
)

type CalculatedCarRecommendation struct {
	SourceCarID      uint
	RecommendedCarID uint
	Score            float64
	Rank             int
}

type RecommendationCalculator interface {
	BuildSimilarities(interactions []repository.RecommendationInteraction, topN int) []CalculatedCarRecommendation
}

type recommendationCalculator struct{}

type weightedCarInteraction struct {
	CarID  uint
	Weight float64
}

type carPair struct {
	LeftCarID  uint
	RightCarID uint
}

func NewRecommendationCalculator() RecommendationCalculator {
	return &recommendationCalculator{}
}

func (c *recommendationCalculator) BuildSimilarities(
	interactions []repository.RecommendationInteraction,
	topN int,
) []CalculatedCarRecommendation {
	if len(interactions) == 0 || topN <= 0 {
		return nil
	}

	userVectors := make(map[uint][]weightedCarInteraction)
	normSquares := make(map[uint]float64)

	for _, interaction := range interactions {
		if interaction.UserID == 0 || interaction.CarID == 0 || interaction.Weight <= 0 {
			continue
		}

		userVectors[interaction.UserID] = append(userVectors[interaction.UserID], weightedCarInteraction{
			CarID:  interaction.CarID,
			Weight: interaction.Weight,
		})
		normSquares[interaction.CarID] += interaction.Weight * interaction.Weight
	}

	if len(userVectors) == 0 {
		return nil
	}

	dotProducts := make(map[carPair]float64)
	for _, vector := range userVectors {
		for i := 0; i < len(vector); i++ {
			for j := i + 1; j < len(vector); j++ {
				left := vector[i]
				right := vector[j]
				pair := normalizeCarPair(left.CarID, right.CarID)
				dotProducts[pair] += left.Weight * right.Weight
			}
		}
	}

	similarities := make(map[uint][]CalculatedCarRecommendation)
	for pair, dotProduct := range dotProducts {
		if dotProduct <= 0 {
			continue
		}

		leftNorm := math.Sqrt(normSquares[pair.LeftCarID])
		rightNorm := math.Sqrt(normSquares[pair.RightCarID])
		if leftNorm == 0 || rightNorm == 0 {
			continue
		}

		score := dotProduct / (leftNorm * rightNorm)
		if score <= 0 {
			continue
		}

		similarities[pair.LeftCarID] = append(similarities[pair.LeftCarID], CalculatedCarRecommendation{
			SourceCarID:      pair.LeftCarID,
			RecommendedCarID: pair.RightCarID,
			Score:            score,
		})
		similarities[pair.RightCarID] = append(similarities[pair.RightCarID], CalculatedCarRecommendation{
			SourceCarID:      pair.RightCarID,
			RecommendedCarID: pair.LeftCarID,
			Score:            score,
		})
	}

	result := make([]CalculatedCarRecommendation, 0)
	for sourceCarID, items := range similarities {
		sort.SliceStable(items, func(i, j int) bool {
			if items[i].Score != items[j].Score {
				return items[i].Score > items[j].Score
			}

			return items[i].RecommendedCarID < items[j].RecommendedCarID
		})

		if len(items) > topN {
			items = items[:topN]
		}

		for idx := range items {
			items[idx].SourceCarID = sourceCarID
			items[idx].Rank = idx + 1
			result = append(result, items[idx])
		}
	}

	sort.SliceStable(result, func(i, j int) bool {
		if result[i].SourceCarID != result[j].SourceCarID {
			return result[i].SourceCarID < result[j].SourceCarID
		}

		return result[i].Rank < result[j].Rank
	})

	return result
}

func normalizeCarPair(leftCarID, rightCarID uint) carPair {
	if leftCarID < rightCarID {
		return carPair{LeftCarID: leftCarID, RightCarID: rightCarID}
	}

	return carPair{LeftCarID: rightCarID, RightCarID: leftCarID}
}
