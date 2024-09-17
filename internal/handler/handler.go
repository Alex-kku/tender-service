package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tender-service/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api")
	{
		api.GET("/ping", h.checkServer)

		tenders := api.Group("/tenders")
		{
			tenders.GET("/", h.getTenders)
			tenders.POST("/new", h.createTender)
			tenders.GET("/my", h.getUserTenders)
			tenders.GET("/:tenderId/status", h.getTenderStatus)
			tenders.PUT("/:tenderId/status", h.updateTenderStatus)
			tenders.PATCH("/:tenderId/edit", h.editTender)
			tenders.PUT("/:tenderId/rollback/:version", h.rollbackTender)
		}

		bids := api.Group("/bids")
		{
			bids.POST("/new", h.createBid)
			bids.GET("/my", h.getUserBids)
			bids.GET("/list/:tenderId", h.getBidsForTender)
			bids.GET("/:bidId/status", h.getBidStatus)
			bids.PUT("/:bidId/status", h.updateBidStatus)
			bids.PATCH("/:bidId/edit", h.editBid)
			bids.PUT("/:bidId/submit_decision", h.submitBidDecision)
			bids.PUT("/:bidId/rollback/:version", h.rollbackBid)

		}
	}
	return router
}

func (h *Handler) checkServer(ctx *gin.Context) {
	ctx.String(http.StatusOK, "ok")
}
