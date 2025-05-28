package auth

import (
	"valuator/pkg/app/model"
	"valuator/pkg/app/provider"
)

type PermissionChecker interface {
	CanReadText(userID model.AuthorID, textID model.TextID) (bool, error)
}

func NewPermissionChecker(authorProvider provider.AuthorProvider) PermissionChecker {
	return &permissionChecker{authorProvider: authorProvider}
}

type permissionChecker struct {
	authorProvider provider.AuthorProvider
}

func (p *permissionChecker) CanReadText(userID model.AuthorID, textID model.TextID) (bool, error) {
	authorIDs, err := p.authorProvider.ListTextAuthors(textID)
	if err != nil {
		return false, err
	}

	for _, authorID := range authorIDs {
		if authorID == userID {
			return true, nil
		}
	}
	return false, nil
}
