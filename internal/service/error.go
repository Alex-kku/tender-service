package service

import (
	"fmt"
)

var (
	ErrUserNotFound    = fmt.Errorf("user not found")
	ErrCannotGetUserId = fmt.Errorf("cannot get user id")

	ErrOrganizationNotFound     = fmt.Errorf("organization not found")
	ErrUserOrganizationNotFound = fmt.Errorf("user organization not found")

	ErrTenderNotFound     = fmt.Errorf("tender not found")
	ErrCannotGetTender    = fmt.Errorf("cannot get tender")
	ErrCannotUpdateTender = fmt.Errorf("cannot edit tender")

	ErrBidNotFound     = fmt.Errorf("bid not found")
	ErrCannotGetBid    = fmt.Errorf("cannot get bid")
	ErrCannotUpdateBid = fmt.Errorf("cannot edit bid")

	ErrNotEnoughPermissions = fmt.Errorf("not enough permissions")
)
