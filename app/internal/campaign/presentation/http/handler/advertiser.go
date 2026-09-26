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

type AdvertiserHandler struct {
	uc  usecase.AdvertiserUseCases
	log appport.Logger
}

func NewAdvertiserHandler(uc usecase.AdvertiserUseCases, log appport.Logger) *AdvertiserHandler {
	return &AdvertiserHandler{uc: uc, log: log}
}

func (h *AdvertiserHandler) Create(c *gin.Context) {
	var req presdto.CreateAdvertiserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "invalid request"})
		return
	}

	out, err := h.uc.Create.Execute(
		c.Request.Context(),
		appdto.CreateAdvertiserInput{
			Name:    req.Name,
			Country: req.Country,
		},
	)
	if err != nil {
		if errors.Is(err, domain.ErrBadInput) {
			c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "bad input"})
			return
		}
		h.log.Error(c.Request.Context(), "create advertiser failed", "error", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, presdto.CreateAdvertiserResponse{
		ID:        out.AdvertiserID.String(),
		Name:      out.Name,
		Country:   out.Country,
		CreatedAt: out.CreatedAt,
		UpdatedAt: out.UpdatedAt,
	})
}

func (h *AdvertiserHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("advertiser_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "bad input"})
		return
	}

	out, err := h.uc.Get.Execute(
		c.Request.Context(),
		appdto.GetAdvertiserInput{
			AdvertiserID: id,
		},
	)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		h.log.Error(c.Request.Context(), "get advertiser failed", "error", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, presdto.GetAdvertiserResponse{
		ID:        out.AdvertiserID.String(),
		Name:      out.Name,
		Country:   out.Country,
		CreatedAt: out.CreatedAt,
		UpdatedAt: out.UpdatedAt,
	})
}

func (h *AdvertiserHandler) Patch(c *gin.Context) {
	id, err := uuid.Parse(c.Param("advertiser_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "bad input"})
		return
	}

	var req presdto.PatchAdvertiserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "invalid request"})
		return
	}

	out, err := h.uc.Patch.Execute(
		c.Request.Context(),
		appdto.PatchAdvertiserInput{
			AdvertiserID: id,
			Name:         req.Name,
			Country:      req.Country,
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
		h.log.Error(c.Request.Context(), "patch advertiser failed", "error", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, presdto.PatchAdvertiserResponse{
		ID:        out.AdvertiserID.String(),
		Name:      out.Name,
		Country:   out.Country,
		CreatedAt: out.CreatedAt,
		UpdatedAt: out.UpdatedAt,
	})
}

func (h *AdvertiserHandler) Put(c *gin.Context) {
	id, err := uuid.Parse(c.Param("advertiser_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "bad input"})
		return
	}

	var req presdto.PutAdvertiserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "invalid request"})
		return
	}

	_, err = h.uc.Put.Execute(
		c.Request.Context(),
		appdto.PutAdvertiserInput{
			AdvertiserID: id,
			Name:         req.Name,
			Country:      req.Country,
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
		h.log.Error(c.Request.Context(), "put advertiser failed", "error", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *AdvertiserHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("advertiser_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, presdto.ErrorResponse{Error: "bad input"})
		return
	}

	_, err = h.uc.Delete.Execute(
		c.Request.Context(),
		appdto.DeleteAdvertiserInput{
			AdvertiserID: id,
		},
	)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		h.log.Error(c.Request.Context(), "delete advertiser failed", "error", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusNoContent)
}
