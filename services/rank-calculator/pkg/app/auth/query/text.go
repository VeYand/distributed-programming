package query

import (
	"errors"
	"rankcalculator/pkg/app/auth"
	"rankcalculator/pkg/app/data"
	"rankcalculator/pkg/app/query"
)

var ErrPermissionDenied = errors.New("permission denied")

func NewAuthorizedStatisticsQueryService(
	next query.StatisticsQueryService,
	checker auth.PermissionChecker,
) query.StatisticsQueryService {
	return &statisticsQueryService{
		next:    next,
		checker: checker,
	}
}

type statisticsQueryService struct {
	next    query.StatisticsQueryService
	checker auth.PermissionChecker
}

func (s *statisticsQueryService) Get(userID string, id string) (data.Statistics, error) {
	canRead, err := s.checker.CanReadText(userID, id)
	if err != nil {
		return data.Statistics{}, err
	}
	if !canRead {
		return data.Statistics{}, ErrPermissionDenied
	}

	return s.next.Get(userID, id)
}
