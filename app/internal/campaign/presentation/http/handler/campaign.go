package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appdto "github.com/iPatrushevSergey/adpace/app/internal/campaign/application/dto"
	appport "github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/usecase"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain"
	presdto "github.com/iPatrushevSergey/adpace/app/internal/campaign/presentation/http/dto"
)

type CampaignHandler struct {
	uc  usecase.CampaignUseCases
	log appport.Logger
}

func NewCampaignHandler(uc usecase.CampaignUseCases, log appport.Logger) *CampaignHandler {
	return &CampaignHandler{uc: uc, log: log}
}

func (h *CampaignHandler) Create(c *gin.Context) {
	var req presdto.CreateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "invalid request"})
		return
	}

	out, err := h.uc.Create.Execute(
		c.Request.Context(),
		appdto.CreateCampaignInput{
			AdvertiserID: req.AdvertiserID,
			Name:         req.Name,
			BudgetTotal:  req.BudgetTotal,
			BudgetDaily:  req.BudgetDaily,
		},
	)
	if err != nil {
		if errors.Is(err, domain.ErrBadInput) {
			c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "bad input"})
			return
		}
		h.log.Error(c.Request.Context(), "create campaign failed", "error", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, presdto.CreateCampaignResponse{
		ID:           out.CampaignID,
		AdvertiserID: out.AdvertiserID,
		Name:         out.Name,
		BudgetTotal:  out.BudgetTotal,
		BudgetDaily:  out.BudgetDaily,
		SpendTotal:   out.SpendTotal,
		SpendToday:   out.SpendToday,
		Status:       out.Status,
		PauseReason:  out.PauseReason,
		CreatedAt:    out.CreatedAt,
		UpdatedAt:    out.UpdatedAt,
	})
}

func (h *CampaignHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("campaign_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "bad input"})
		return
	}

	out, err := h.uc.Get.Execute(
		c.Request.Context(),
		appdto.GetCampaignInput{
			CampaignID: id.String(),
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBadInput):
			c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "bad input"})
			return
		case errors.Is(err, domain.ErrNotFound):
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		h.log.Error(c.Request.Context(), "get campaign failed", "error", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, presdto.GetCampaignResponse{
		ID:           out.CampaignID,
		AdvertiserID: out.AdvertiserID,
		Name:         out.Name,
		BudgetTotal:  out.BudgetTotal,
		BudgetDaily:  out.BudgetDaily,
		SpendTotal:   out.SpendTotal,
		SpendToday:   out.SpendToday,
		Status:       out.Status,
		PauseReason:  out.PauseReason,
		CreatedAt:    out.CreatedAt,
		UpdatedAt:    out.UpdatedAt,
	})
}

func (h *CampaignHandler) Patch(c *gin.Context) {
	id, err := uuid.Parse(c.Param("campaign_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "bad input"})
		return
	}

	var req presdto.PatchCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "invalid request"})
		return
	}

	out, err := h.uc.Patch.Execute(
		c.Request.Context(),
		appdto.PatchCampaignInput{
			CampaignID:  id.String(),
			Name:        req.Name,
			BudgetTotal: req.BudgetTotal,
			BudgetDaily: req.BudgetDaily,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBadInput):
			c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "bad input"})
			return
		case errors.Is(err, domain.ErrNotFound):
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		h.log.Error(c.Request.Context(), "patch campaign failed", "error", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, presdto.PatchCampaignResponse{
		ID:           out.CampaignID,
		AdvertiserID: out.AdvertiserID,
		Name:         out.Name,
		BudgetTotal:  out.BudgetTotal,
		BudgetDaily:  out.BudgetDaily,
		SpendTotal:   out.SpendTotal,
		SpendToday:   out.SpendToday,
		Status:       out.Status,
		PauseReason:  out.PauseReason,
		CreatedAt:    out.CreatedAt,
		UpdatedAt:    out.UpdatedAt,
	})
}

func (h *CampaignHandler) Put(c *gin.Context) {
	id, err := uuid.Parse(c.Param("campaign_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "bad input"})
		return
	}

	var req presdto.PutCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "invalid request"})
		return
	}

	_, err = h.uc.Put.Execute(
		c.Request.Context(),
		appdto.PutCampaignInput{
			CampaignID:  id.String(),
			Name:        req.Name,
			BudgetTotal: req.BudgetTotal,
			BudgetDaily: req.BudgetDaily,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBadInput):
			c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "bad input"})
			return
		case errors.Is(err, domain.ErrNotFound):
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		h.log.Error(c.Request.Context(), "put campaign failed", "error", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *CampaignHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("campaign_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "bad input"})
		return
	}

	_, err = h.uc.Delete.Execute(
		c.Request.Context(),
		appdto.DeleteCampaignInput{
			CampaignID: id.String(),
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBadInput):
			c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "bad input"})
			return
		case errors.Is(err, domain.ErrNotFound):
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		h.log.Error(c.Request.Context(), "delete campaign failed", "error", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *CampaignHandler) Pause(c *gin.Context) {
	id, err := uuid.Parse(c.Param("campaign_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "bad input"})
		return
	}

	var req presdto.PauseCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "invalid request"})
		return
	}

	out, err := h.uc.Pause.Execute(
		c.Request.Context(),
		appdto.PauseCampaignInput{
			CampaignID: id.String(),
			Reason:     req.Reason,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBadInput):
			c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "bad input"})
			return
		case errors.Is(err, domain.ErrNotFound):
			c.AbortWithStatus(http.StatusNotFound)
			return
		case errors.Is(err, domain.ErrConflict):
			c.AbortWithStatus(http.StatusConflict)
			return
		}
		h.log.Error(c.Request.Context(), "pause campaign failed", "error", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, presdto.PauseCampaignResponse{
		ID:          out.CampaignID,
		Status:      out.Status,
		PauseReason: out.PauseReason,
		UpdatedAt:   out.UpdatedAt,
	})
}

func (h *CampaignHandler) Resume(c *gin.Context) {
	id, err := uuid.Parse(c.Param("campaign_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "bad input"})
		return
	}

	out, err := h.uc.Resume.Execute(
		c.Request.Context(),
		appdto.ResumeCampaignInput{
			CampaignID: id.String(),
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBadInput):
			c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "bad input"})
			return
		case errors.Is(err, domain.ErrNotFound):
			c.AbortWithStatus(http.StatusNotFound)
			return
		case errors.Is(err, domain.ErrConflict):
			c.AbortWithStatus(http.StatusConflict)
			return
		}
		h.log.Error(c.Request.Context(), "resume campaign failed", "error", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, presdto.ResumeCampaignResponse{
		ID:        out.CampaignID,
		Status:    out.Status,
		UpdatedAt: out.UpdatedAt,
	})
}
