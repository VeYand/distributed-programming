package query

import (
	"errors"
	"valuator/pkg/app/auth"
	"valuator/pkg/app/data"
	"valuator/pkg/app/model"
	"valuator/pkg/app/query"
)

var ErrPermissionDenied = errors.New("permission denied")

func NewAuthorizedTextQueryService(
	next query.TextQueryService,
	checker auth.PermissionChecker,
) query.TextQueryService {
	return &textQueryService{
		next:    next,
		checker: checker,
	}
}

type textQueryService struct {
	next    query.TextQueryService
	checker auth.PermissionChecker
}

func (s *textQueryService) Get(userID string, id string) (data.TextData, error) {
	canRead, err := s.checker.CanReadText(userID, model.TextID(id))
	if err != nil {
		return data.TextData{}, err
	}
	if !canRead {
		return data.TextData{}, ErrPermissionDenied
	}

	return s.next.Get(userID, id)
}
