package service

import (
	"math"
	"testing"

	"car-service/internal/repository"
)

func TestRecommendationCalculatorBuildSimilarities(t *testing.T) {
	calculator := NewRecommendationCalculator()

	result := calculator.BuildSimilarities([]repository.RecommendationInteraction{
		{UserID: 1, CarID: 1, Weight: 3},
		{UserID: 1, CarID: 2, Weight: 1},
		{UserID: 2, CarID: 1, Weight: 3},
		{UserID: 2, CarID: 2, Weight: 1},
		{UserID: 3, CarID: 1, Weight: 1},
		{UserID: 3, CarID: 3, Weight: 3},
	}, 10)

	if len(result) != 4 {
		t.Fatalf("expected 4 recommendations, got %d", len(result))
	}

	assertRecommendation(t, result[0], 1, 2, 6/math.Sqrt(38), 1)
	assertRecommendation(t, result[1], 1, 3, 1/math.Sqrt(19), 2)
	assertRecommendation(t, result[2], 2, 1, 6/math.Sqrt(38), 1)
	assertRecommendation(t, result[3], 3, 1, 1/math.Sqrt(19), 1)
}

func TestRecommendationCalculatorBuildSimilaritiesAppliesTopN(t *testing.T) {
	calculator := NewRecommendationCalculator()

	result := calculator.BuildSimilarities([]repository.RecommendationInteraction{
		{UserID: 1, CarID: 1, Weight: 2},
		{UserID: 1, CarID: 2, Weight: 2},
		{UserID: 1, CarID: 3, Weight: 1},
		{UserID: 2, CarID: 1, Weight: 2},
		{UserID: 2, CarID: 2, Weight: 2},
	}, 1)

	for _, item := range result {
		if item.SourceCarID == 1 && item.RecommendedCarID != 2 {
			t.Fatalf("expected car 1 to keep only car 2 as top recommendation, got %#v", item)
		}
	}
}

func assertRecommendation(
	t *testing.T,
	item CalculatedCarRecommendation,
	sourceCarID, recommendedCarID uint,
	expectedScore float64,
	expectedRank int,
) {
	t.Helper()

	if item.SourceCarID != sourceCarID {
		t.Fatalf("expected source_car_id %d, got %d", sourceCarID, item.SourceCarID)
	}

	if item.RecommendedCarID != recommendedCarID {
		t.Fatalf("expected recommended_car_id %d, got %d", recommendedCarID, item.RecommendedCarID)
	}

	if math.Abs(item.Score-expectedScore) > 0.000001 {
		t.Fatalf("expected score %.6f, got %.6f", expectedScore, item.Score)
	}

	if item.Rank != expectedRank {
		t.Fatalf("expected rank %d, got %d", expectedRank, item.Rank)
	}
}
