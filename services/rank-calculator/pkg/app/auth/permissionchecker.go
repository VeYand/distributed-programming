package auth

type PermissionChecker interface {
	CanReadText(userID string, textID string) (bool, error)
}
