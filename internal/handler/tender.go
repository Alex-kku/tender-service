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

func (h *Handler) getTenders(ctx *gin.Context) {
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

	serviceType, ok := ctx.GetQueryArray("service_type")
	if ok {
		if err := serviceTypeValidate(serviceType); err != nil {
			newErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
	}

	tenders, err := h.services.Tender.GetTenders(ctx.Request.Context(), limit, offset, serviceType)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, tenders)
}

func (h *Handler) createTender(ctx *gin.Context) {
	var input entity.CreateTenderInput

	if err := ctx.Bind(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := h.services.Auth.GetUserId(ctx.Request.Context(), input.CreatorUsername)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
			return
		}
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	ok, err := h.services.Auth.CheckResponsibility(ctx.Request.Context(), userId, input.OrganizationId)
	if !ok {
		if errors.Is(err, service.ErrOrganizationNotFound) {
			newErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
		newErrorResponse(ctx, http.StatusForbidden, err.Error())
		return
	}

	tender, err := h.services.Tender.CreateTender(ctx.Request.Context(), userId, input)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, tender)
}

func (h *Handler) getUserTenders(ctx *gin.Context) {
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

	tenders, err := h.services.Tender.GetTendersByUserId(ctx.Request.Context(), limit, offset, userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, tenders)
}

func (h *Handler) getTenderStatus(ctx *gin.Context) {
	tenderId := ctx.Param("tenderId")
	if err := tenderIdValidate(tenderId); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	username := ctx.Query("username")
	var userId *string
	if username != "" {
		id, err := h.services.Auth.GetUserId(ctx.Request.Context(), username)
		if err != nil {
			newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
			return
		}
		userId = &id
	}

	status, err := h.services.Tender.GetTenderStatus(ctx.Request.Context(), tenderId, userId)
	if err != nil {
		if errors.Is(err, service.ErrTenderNotFound) {
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

	ctx.JSON(http.StatusOK, status)
}

func (h *Handler) updateTenderStatus(ctx *gin.Context) {
	tenderId := ctx.Param("tenderId")
	if err := tenderIdValidate(tenderId); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	status := ctx.Query("status")
	if err := tenderStatusValidate(status); err != nil {
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

	tender, err := h.services.Tender.UpdateTenderStatus(ctx.Request.Context(), tenderId, status, userId)
	if err != nil {
		if errors.Is(err, service.ErrTenderNotFound) {
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

	ctx.JSON(http.StatusOK, tender)
}

func (h *Handler) editTender(ctx *gin.Context) {
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

	var input entity.EditTenderInput
	if err = ctx.Bind(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if input.Name == nil && input.Description == nil && input.ServiceType == nil {
		newErrorResponse(ctx, http.StatusBadRequest, "no parameters for edit")
		return
	}
	if input.Name != nil {
		if err := tenderNameValidate(*input.Name); err != nil {
			newErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
	}
	if input.Description != nil {
		if err := tenderDescriptionValidate(*input.Description); err != nil {
			newErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
	}
	if input.ServiceType != nil {
		if err := serviceTypeValidate([]string{*input.ServiceType}); err != nil {
			newErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
	}

	tender, err := h.services.Tender.EditTender(ctx.Request.Context(), tenderId, userId, input)
	if err != nil {
		if errors.Is(err, service.ErrTenderNotFound) {
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

	ctx.JSON(http.StatusOK, tender)
}

func (h *Handler) rollbackTender(ctx *gin.Context) {
	tenderId := ctx.Param("tenderId")
	if err := tenderIdValidate(tenderId); err != nil {
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
		if errors.Is(err, service.ErrUserNotFound) {
			newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
			return
		}
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	tender, err := h.services.Tender.RollbackTender(ctx.Request.Context(), tenderId, version, userId)
	if err != nil {
		if errors.Is(err, service.ErrTenderNotFound) {
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

	ctx.JSON(http.StatusOK, tender)
}

func limitValidate(limit int) error {
	if limit < 0 || limit > 50 {
		return fmt.Errorf("invalid limit, min = 0, max = 50")
	}
	return nil
}

func offsetValidate(offset int) error {
	if offset < 0 {
		return fmt.Errorf("invalid offset, min = 0")
	}
	return nil
}

func tenderNameValidate(name string) error {
	if name == "" {
		return fmt.Errorf("tenderName is empty")
	}
	if len(name) > 100 {
		return fmt.Errorf("tenderName is too long, maxLength=100")
	}
	return nil
}

func tenderDescriptionValidate(description string) error {
	if description == "" {
		return fmt.Errorf("tenderDescription is empty")
	}
	if len(description) > 500 {
		return fmt.Errorf("tenderDescription is too long, maxLength=500")
	}
	return nil
}

func serviceTypeValidate(serviceType []string) error {
	for _, s := range serviceType {
		if !slices.Contains([]string{"Construction", "Delivery", "Manufacture"}, s) {
			return fmt.Errorf("invalid service_type, must be 'Construction'/'Delivery'/'Manufacture'")
		}
	}
	return nil
}

func usernameValidate(username string) error {
	if username == "" {
		return fmt.Errorf("username is empty")
	}
	return nil
}

func tenderIdValidate(tenderId string) error {
	if tenderId == "" {
		return fmt.Errorf("tenderId is empty")
	}
	if len(tenderId) > 100 {
		return fmt.Errorf("tenderId is too long, maxLength=100")
	}
	return nil
}

func tenderStatusValidate(status string) error {
	if status == "" {
		return fmt.Errorf("status is empty")
	}
	if !slices.Contains([]string{"Created", "Published", "Closed"}, status) {
		return fmt.Errorf("invalid status, must be 'Created'/'Published'/'Closed'")
	}
	return nil
}

func versionValidate(version int) error {
	if version < 1 {
		return fmt.Errorf("invalid version, min = 1")
	}
	return nil
}
