package hospital_server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type IService interface {
	GetSpecs() ([]Spec, error)
	GetDoctors(specId string) ([]Doctor, error)
	GetSlots(doctorId string) ([]Appointment, error)

	MakeAppointment(date time.Time, time, doctorId, userId string) error
	CancelAppointment(date time.Time, time, doctorId, userId string) error

	GetUserProfile(userId string) (*User, error)
	GetUserAppointments(userId string) ([]Appointment, error)
}

type Handler struct {
	service IService
}

func NewHandler(service IService) *Handler {
	return &Handler{
		service: service,
	}
}

func (handler *Handler) GetSpecs(c *gin.Context) {

	specs, err := handler.service.GetSpecs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	out := GetSpecsOut{
		Specs: specs,
	}

	c.JSON(http.StatusOK, out)
}

func (handler *Handler) GetDoctors(c *gin.Context) {

	specId := c.Query("specId")
	if specId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid specId"})
		return
	}

	doctors, err := handler.service.GetDoctors(specId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	out := GetDoctorsOut{
		Doctors: doctors,
	}

	c.JSON(http.StatusOK, out)
}

func (handler *Handler) GetSlots(c *gin.Context) {

	doctorId := c.Param("id")
	if doctorId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor id"})
		return
	}

	fmt.Println("АБОБА 1")

	slots, err := handler.service.GetSlots(doctorId)
	fmt.Println(slots[0], "АБОБА 2")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	out := GetSlotsOut{
		Slots: slots,
	}

	c.JSON(http.StatusOK, out)
}

func (handler *Handler) MakeAppointment(c *gin.Context) {

	var dto MakeAppointmentDTO
	if err := json.NewDecoder(c.Request.Body).Decode(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := handler.service.MakeAppointment(dto.Date, dto.Time, dto.DoctorId, dto.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (handler *Handler) CancelAppointment(c *gin.Context) {

	var dto CancelAppointmentDTO
	if err := json.NewDecoder(c.Request.Body).Decode(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := handler.service.CancelAppointment(dto.Date, dto.Time, dto.DoctorId, dto.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (handler *Handler) GetUserProfile(c *gin.Context) {

	userId := c.Query("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	profile, err := handler.service.GetUserProfile(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	out := GetUserProfileOut{
		Profile: *profile,
	}
	c.JSON(http.StatusOK, out)
}

func (handler *Handler) GetUserAppointments(c *gin.Context) {

	userId := c.Query("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	appointments, err := handler.service.GetUserAppointments(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	out := GetUserAppointmentsOut{
		Appointments: appointments,
	}
	c.JSON(http.StatusOK, out)
}
