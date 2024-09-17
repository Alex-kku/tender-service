package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
	"strconv"
	"tender-service/internal/entity"
	"tender-service/internal/service"
)

func (h *Handler) createBid(ctx *gin.Context) {
	var input entity.CreateBidInput

	if err := ctx.Bind(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	bid, err := h.services.Bid.CreateBid(ctx.Request.Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrTenderNotFound) {
			newErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
			return
		}
		if errors.Is(err, service.ErrNotEnoughPermissions) {
			newErrorResponse(ctx, http.StatusForbidden, err.Error())
			return
		}
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, bid)
}

func (h *Handler) getUserBids(ctx *gin.Context) {
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "5"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := limitValidate(limit); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	offset, err := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := offsetValidate(offset); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	username := ctx.Query("username")
	if err := usernameValidate(username); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	userId, err := h.services.Auth.GetUserId(ctx.Request.Context(), username)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
			return
		}
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	bids, err := h.services.Bid.GetBidsByUserId(ctx.Request.Context(), limit, offset, userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, bids)
}

func (h *Handler) getBidsForTender(ctx *gin.Context) {
	tenderId := ctx.Param("tenderId")
	if err := tenderIdValidate(tenderId); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	username := ctx.Query("username")
	if err := usernameValidate(username); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	userId, err := h.services.Auth.GetUserId(ctx.Request.Context(), username)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
			return
		}
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "5"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := limitValidate(limit); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	offset, err := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := offsetValidate(offset); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	bids, err := h.services.Bid.GetBidsForTender(ctx.Request.Context(), tenderId, userId, limit, offset)
	if err != nil {
		if errors.Is(err, service.ErrTenderNotFound) {
			newErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
			return
		}
		if errors.Is(err, service.ErrNotEnoughPermissions) {
			newErrorResponse(ctx, http.StatusForbidden, err.Error())
			return
		}
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, bids)
}

func (h *Handler) getBidStatus(ctx *gin.Context) {
	bidId := ctx.Param("bidId")
	if err := bidIdValidate(bidId); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	username := ctx.Query("username")
	if err := usernameValidate(username); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	userId, err := h.services.Auth.GetUserId(ctx.Request.Context(), username)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
			return
		}
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	status, err := h.services.Bid.GetBidStatus(ctx.Request.Context(), bidId, userId)
	if err != nil {
		//todo уточнить пришедшую ошибку и расписать коды ответа
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, status)
}

func (h *Handler) updateBidStatus(ctx *gin.Context) {
	bidId := ctx.Param("bidId")
	if err := bidIdValidate(bidId); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	status := ctx.Query("status")
	if err := bidStatusValidate(status); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	username := ctx.Query("username")
	if err := usernameValidate(username); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	userId, err := h.services.Auth.GetUserId(ctx.Request.Context(), username)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
			return
		}
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	bid, err := h.services.Bid.UpdateBidStatus(ctx.Request.Context(), bidId, status, userId)
	if err != nil {
		if errors.Is(err, service.ErrBidNotFound) {
			newErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, service.ErrNotEnoughPermissions) {
			newErrorResponse(ctx, http.StatusForbidden, err.Error())
			return
		}
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, bid)
}

func (h *Handler) editBid(ctx *gin.Context) {
	bidId := ctx.Param("bidId")
	if err := bidIdValidate(bidId); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	username := ctx.Query("username")
	if err := usernameValidate(username); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	userId, err := h.services.Auth.GetUserId(ctx.Request.Context(), username)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
			return
		}
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var input entity.EditBidInput
	if err = ctx.Bind(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if input.Name == nil && input.Description == nil {
		newErrorResponse(ctx, http.StatusBadRequest, "no parameters for edit")
		return
	}
	if input.Name != nil {
		if err := BidNameValidate(*input.Name); err != nil {
			newErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
	}
	if input.Description != nil {
		if err := BidNameValidate(*input.Description); err != nil {
			newErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
	}

	bid, err := h.services.Bid.EditBid(ctx.Request.Context(), bidId, userId, input)
	if err != nil {
		if errors.Is(err, service.ErrBidNotFound) {
			newErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, service.ErrNotEnoughPermissions) {
			newErrorResponse(ctx, http.StatusForbidden, err.Error())
			return
		}
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, bid)
}

func (h *Handler) submitBidDecision(ctx *gin.Context) {
	bidId := ctx.Param("bidId")
	if err := bidIdValidate(bidId); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	decision := ctx.Query("decision")
	if err := bidDecisionValidate(decision); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	username := ctx.Query("username")
	if err := usernameValidate(username); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	userId, err := h.services.Auth.GetUserId(ctx.Request.Context(), username)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
			return
		}
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	bid, err := h.services.Bid.SubmitBidDecision(ctx.Request.Context(), bidId, decision, userId)
	if err != nil {
		if errors.Is(err, service.ErrBidNotFound) {
			newErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, service.ErrNotEnoughPermissions) {
			newErrorResponse(ctx, http.StatusForbidden, err.Error())
			return
		}
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, bid)
}

func (h *Handler) rollbackBid(ctx *gin.Context) {
	bidId := ctx.Param("bidId")
	if err := bidIdValidate(bidId); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	version, err := strconv.Atoi(ctx.Param("version"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := versionValidate(version); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	username := ctx.Query("username")
	if err := usernameValidate(username); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	userId, err := h.services.Auth.GetUserId(ctx.Request.Context(), username)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	bid, err := h.services.Bid.RollbackBid(ctx.Request.Context(), bidId, version, userId)
	if err != nil {
		if errors.Is(err, service.ErrBidNotFound) {
			newErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, service.ErrNotEnoughPermissions) {
			newErrorResponse(ctx, http.StatusForbidden, err.Error())
			return
		}
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, bid)
}

func bidIdValidate(bidId string) error {
	if bidId == "" {
		return fmt.Errorf("bidId is empty")
	}
	if len(bidId) > 100 {
		return fmt.Errorf("bidId is too long, maxLength=100")
	}
	return nil
}

func bidStatusValidate(status string) error {
	if status == "" {
		return fmt.Errorf("status is empty")
	}
	if !slices.Contains([]string{"Created", "Published", "Canceled"}, status) {
		return fmt.Errorf("invalid status, must be 'Created'/'Published'/'Canceled'")
	}
	return nil
}

func bidDecisionValidate(decision string) error {
	if decision == "" {
		return fmt.Errorf("decision is empty")
	}
	if !slices.Contains([]string{"Approved", "Rejected"}, decision) {
		return fmt.Errorf("invalid decision, must be 'Approved'/'Rejected'")
	}
	return nil
}

func BidNameValidate(name string) error {
	if name == "" {
		return fmt.Errorf("bidName is empty")
	}
	if len(name) > 100 {
		return fmt.Errorf("bidName is too long, maxLength=100")
	}
	return nil
}

func BidDescriptionValidate(description string) error {
	if description == "" {
		return fmt.Errorf("bidDescription is empty")
	}
	if len(description) > 500 {
		return fmt.Errorf("bidDescription is too long, maxLength=500")
	}
	return nil
}
