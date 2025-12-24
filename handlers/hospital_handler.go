package handlers

import (
	"agnos_candidate_assignment/models"
	"agnos_candidate_assignment/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HospitalHandler struct {
	Repo repositories.HospitalRepositoryInterface
}

func NewHospitalHandler(repo repositories.HospitalRepositoryInterface) *HospitalHandler {
	return &HospitalHandler{Repo: repo}
}

type createHospitalRequest struct {
	Name string `json:"name" binding:"required"`
}

func (hospitalHandler *HospitalHandler) Create(c *gin.Context) {
	var req createHospitalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospital := &models.Hospital{
		Name: req.Name,
	}

	if err := hospitalHandler.Repo.Create(hospital); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create hospital"})
		return
	}
	c.JSON(http.StatusCreated, hospital)
}
