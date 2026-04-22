package service

import (
	"context"
	"sort"
	"time"

	"car-service/internal/domains"
	"car-service/internal/repository"
	"car-service/internal/storage"
)

const (
	defaultRecommendationLimit      = 3
	defaultPopularRecommendationTop = 20
	defaultSimilarCarsPerSource     = 10
)

type RecommendationItem struct {
	CatalogCar
	Score *float64
}

type RecommendationsResult struct {
	Items []RecommendationItem
}

type RebuildRecommendationsResult struct {
	InteractionsCount    int
	SourceCarsCount      int
	RecommendationsCount int
	RebuiltAt            time.Time
}

type RecommendationService interface {
	GetMyRecommendations(ctx context.Context, userID uint) (RecommendationsResult, error)
	RebuildRecommendations(ctx context.Context) (RebuildRecommendationsResult, error)
}

type recommendationService struct {
	recommendationRepository repository.RecommendationRepository
	carRepository            repository.CarRepository
	imageStorage             storage.ImageStorageService
	calculator               RecommendationCalculator
}

type scoredCar struct {
	CarID uint
	Score float64
}

func NewRecommendationService(
	recommendationRepository repository.RecommendationRepository,
	carRepository repository.CarRepository,
	imageStorage storage.ImageStorageService,
	calculator RecommendationCalculator,
) RecommendationService {
	return &recommendationService{
		recommendationRepository: recommendationRepository,
		carRepository:            carRepository,
		imageStorage:             imageStorage,
		calculator:               calculator,
	}
}

func (s *recommendationService) GetMyRecommendations(
	ctx context.Context,
	userID uint,
) (RecommendationsResult, error) {
	if userID == 0 {
		return RecommendationsResult{}, validationError("user_id must be greater than zero")
	}

	interactions, err := s.recommendationRepository.ListUserAggregatedInteractions(ctx, userID)
	if err != nil {
		return RecommendationsResult{}, err
	}

	excludedCarIDs := make(map[uint]struct{}, len(interactions))
	for _, interaction := range interactions {
		excludedCarIDs[interaction.CarID] = struct{}{}
	}

	if len(interactions) == 0 {
		items, err := s.buildFallbackRecommendations(ctx, excludedCarIDs)
		if err != nil {
			return RecommendationsResult{}, err
		}

		return RecommendationsResult{Items: items}, nil
	}

	sourceCarIDs := make([]uint, 0, len(interactions))
	interactionWeights := make(map[uint]float64, len(interactions))
	for _, interaction := range interactions {
		sourceCarIDs = append(sourceCarIDs, interaction.CarID)
		interactionWeights[interaction.CarID] = interaction.Weight
	}

	precomputed, err := s.recommendationRepository.ListBySourceCarIDs(ctx, sourceCarIDs)
	if err != nil {
		return RecommendationsResult{}, err
	}

	candidateScores := make(map[uint]float64)
	for _, recommendation := range precomputed {
		if _, excluded := excludedCarIDs[recommendation.RecommendedCarID]; excluded {
			continue
		}

		candidateScores[recommendation.RecommendedCarID] +=
			interactionWeights[recommendation.SourceCarID] * recommendation.Score
	}

	if len(candidateScores) == 0 {
		items, err := s.buildFallbackRecommendations(ctx, excludedCarIDs)
		if err != nil {
			return RecommendationsResult{}, err
		}

		return RecommendationsResult{Items: items}, nil
	}

	rankedCars := rankScoredCars(candidateScores)
	if len(rankedCars) > defaultRecommendationLimit {
		rankedCars = rankedCars[:defaultRecommendationLimit]
	}

	items, err := s.loadRecommendationItems(ctx, rankedCars, true)
	if err != nil {
		return RecommendationsResult{}, err
	}

	if len(items) == 0 {
		fallbackItems, err := s.buildFallbackRecommendations(ctx, excludedCarIDs)
		if err != nil {
			return RecommendationsResult{}, err
		}

		return RecommendationsResult{Items: fallbackItems}, nil
	}

	return RecommendationsResult{Items: items}, nil
}

func (s *recommendationService) RebuildRecommendations(
	ctx context.Context,
) (RebuildRecommendationsResult, error) {
	interactions, err := s.recommendationRepository.ListAggregatedInteractions(ctx)
	if err != nil {
		return RebuildRecommendationsResult{}, err
	}

	calculated := s.calculator.BuildSimilarities(interactions, defaultSimilarCarsPerSource)
	recommendations := make([]domains.CarRecommendation, 0, len(calculated))
	for _, item := range calculated {
		recommendations = append(recommendations, domains.CarRecommendation{
			SourceCarID:      item.SourceCarID,
			RecommendedCarID: item.RecommendedCarID,
			Score:            item.Score,
			Rank:             item.Rank,
		})
	}

	if err := s.recommendationRepository.ReplaceAll(ctx, recommendations); err != nil {
		return RebuildRecommendationsResult{}, err
	}

	return RebuildRecommendationsResult{
		InteractionsCount:    len(interactions),
		SourceCarsCount:      countDistinctCars(interactions),
		RecommendationsCount: len(recommendations),
		RebuiltAt:            time.Now().UTC(),
	}, nil
}

func (s *recommendationService) buildFallbackRecommendations(
	ctx context.Context,
	excludedCarIDs map[uint]struct{},
) ([]RecommendationItem, error) {
	popularCars, err := s.recommendationRepository.ListPopularCarScores(ctx, defaultPopularRecommendationTop)
	if err != nil {
		return nil, err
	}

	rankedCars := make([]scoredCar, 0, len(popularCars))
	for _, car := range popularCars {
		if _, excluded := excludedCarIDs[car.CarID]; excluded {
			continue
		}

		rankedCars = append(rankedCars, scoredCar{
			CarID: car.CarID,
			Score: car.Score,
		})
	}

	if len(rankedCars) > defaultRecommendationLimit {
		rankedCars = rankedCars[:defaultRecommendationLimit]
	}

	return s.loadRecommendationItems(ctx, rankedCars, false)
}

func (s *recommendationService) loadRecommendationItems(
	ctx context.Context,
	rankedCars []scoredCar,
	includeScore bool,
) ([]RecommendationItem, error) {
	if len(rankedCars) == 0 {
		return []RecommendationItem{}, nil
	}

	carIDs := make([]uint, 0, len(rankedCars))
	scoreByCarID := make(map[uint]float64, len(rankedCars))
	for _, car := range rankedCars {
		carIDs = append(carIDs, car.CarID)
		scoreByCarID[car.CarID] = car.Score
	}

	cars, err := s.carRepository.ListByIDs(ctx, carIDs)
	if err != nil {
		return nil, err
	}

	carByID := make(map[uint]CatalogCar, len(cars))
	for i := range cars {
		cars[i].CarImages = catalogImages(cars[i].CarImages)

		item, err := toCatalogCar(ctx, s.imageStorage, cars[i])
		if err != nil {
			return nil, err
		}

		carByID[item.ID] = item
	}

	items := make([]RecommendationItem, 0, len(rankedCars))
	for _, rankedCar := range rankedCars {
		car, ok := carByID[rankedCar.CarID]
		if !ok {
			continue
		}

		recommendationItem := RecommendationItem{
			CatalogCar: car,
		}
		if includeScore {
			score := scoreByCarID[rankedCar.CarID]
			recommendationItem.Score = &score
		}

		items = append(items, recommendationItem)
	}

	return items, nil
}

func rankScoredCars(scores map[uint]float64) []scoredCar {
	rankedCars := make([]scoredCar, 0, len(scores))
	for carID, score := range scores {
		if score <= 0 {
			continue
		}

		rankedCars = append(rankedCars, scoredCar{
			CarID: carID,
			Score: score,
		})
	}

	sort.SliceStable(rankedCars, func(i, j int) bool {
		if rankedCars[i].Score != rankedCars[j].Score {
			return rankedCars[i].Score > rankedCars[j].Score
		}

		return rankedCars[i].CarID < rankedCars[j].CarID
	})

	return rankedCars
}

func countDistinctCars(interactions []repository.RecommendationInteraction) int {
	uniqueCars := make(map[uint]struct{}, len(interactions))
	for _, interaction := range interactions {
		uniqueCars[interaction.CarID] = struct{}{}
	}

	return len(uniqueCars)
}
