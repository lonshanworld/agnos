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
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// Register godoc
// @Summary      Register a new staff member
// @Description  Create a new staff account for a hospital
// @Tags         staff
// @Accept       json
// @Produce      json
// @Param        hospital path string true "Hospital name"
// @Param        request body registerReq true "Staff registration request"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /:hospital/staff/create [post]
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
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// Login godoc
// @Summary      Staff login
// @Description  Authenticate staff and receive JWT token
// @Tags         staff
// @Accept       json
// @Produce      json
// @Param        hospital path string true "Hospital name"
// @Param        request body loginReq true "Login credentials"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /:hospital/staff/login [post]
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
