package handlers

import (
	"agnos_candidate_assignment/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StaffHandler struct {
	authService services.AuthServiceInterface
}

func NewStaffHandler(authser services.AuthServiceInterface) *StaffHandler {
	return &StaffHandler{authService: authser}
}

type registerReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (staffHandler *StaffHandler) Register(c *gin.Context) {
	hospital := c.Param("hospital")
	var req registerReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	staff, err := staffHandler.authService.Register(hospital, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"staff_id": staff.ID, "username": staff.UserName})
}

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (staffhandler *StaffHandler) Login(c *gin.Context) {
	hospital := c.Param("hospital")
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, staff, err := staffhandler.authService.Login(hospital, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "staff_id": staff.ID, "username": staff.UserName, "hospital_id": staff.HospitalID})
}
