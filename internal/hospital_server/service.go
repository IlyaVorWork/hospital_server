package hospital_server

import (
	pkg "Server/internal/pkg/common"
	"strconv"
	"strings"
	. "time"
)

type IProvider interface {
	GetSpecs() ([]Spec, error)
	GetDoctors(specId string) ([]Doctor, error)
	GetSlots(doctorId string) ([]Appointment, error)

	MakeAppointment(date Time, time, doctorId, userId string) error
	CancelAppointment(date Time, time, doctorId string) error

	GetUserProfile(userId string) (*User, error)
	GetUserAppointments(userId string) ([]Appointment, error)

	DoesUserExist(userId string) bool
	GetAppointment(date Time, time, doctorId string) (*Appointment, error)
}

type Service struct {
	provider IProvider
}

func NewService(provider IProvider) *Service {
	return &Service{
		provider: provider,
	}
}

func (service *Service) GetSpecs() ([]Spec, error) {

	specs, err := service.provider.GetSpecs()
	if err != nil {
		return nil, err
	}

	return specs, nil
}

func (service *Service) GetDoctors(specId string) ([]Doctor, error) {

	doctors, err := service.provider.GetDoctors(specId)
	if err != nil {
		return nil, err
	}

	return doctors, nil
}

func (service *Service) GetSlots(doctorId string) ([]Appointment, error) {

	slots, err := service.provider.GetSlots(doctorId)
	if err != nil {
		return nil, err
	}

	return slots, nil
}

func (service *Service) MakeAppointment(date Time, time, doctorId, userId string) error {

	exist := service.provider.DoesUserExist(userId)
	if !exist {
		return pkg.UserNotExistErr
	}

	slot, err := service.provider.GetAppointment(date, time, doctorId)
	if err != nil {
		return pkg.SlotNotExistErr
	}

	currentDate := Now()
	hoursAndMinutes := strings.Split(slot.Time, ":")
	hours, err := strconv.Atoi(hoursAndMinutes[0])
	if err != nil {
		return err
	}
	minutes, err := strconv.Atoi(hoursAndMinutes[1])
	if err != nil {
		return err
	}

	slotDate := Date(date.Year(), date.Month(), date.Day(), hours, minutes, 0, 0, currentDate.Location())
	if currentDate.After(slotDate) {
		return pkg.SlotIsExpiredErr
	}

	err = service.provider.MakeAppointment(date, time, doctorId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (service *Service) CancelAppointment(date Time, time, doctorId, userId string) error {

	exist := service.provider.DoesUserExist(userId)
	if !exist {
		return pkg.UserNotExistErr
	}

	appointment, err := service.provider.GetAppointment(date, time, doctorId)
	if err != nil {
		return pkg.SlotNotExistErr
	}

	currentDate := Now()
	hoursAndMinutes := strings.Split(appointment.Time, ":")
	hours, err := strconv.Atoi(hoursAndMinutes[0])
	if err != nil {
		return err
	}
	minutes, err := strconv.Atoi(hoursAndMinutes[1])
	if err != nil {
		return err
	}

	appointmentDate := Date(date.Year(), date.Month(), date.Day(), hours, minutes, 0, 0, currentDate.Location())
	if currentDate.After(appointmentDate) {
		return pkg.SlotIsExpiredErr
	}

	err = service.provider.CancelAppointment(date, time, doctorId)
	if err != nil {
		return err
	}

	return nil
}

func (service *Service) GetUserProfile(userId string) (*User, error) {

	exist := service.provider.DoesUserExist(userId)
	if !exist {
		return nil, pkg.UserNotExistErr
	}

	profile, err := service.provider.GetUserProfile(userId)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (service *Service) GetUserAppointments(userId string) ([]Appointment, error) {

	exist := service.provider.DoesUserExist(userId)
	if !exist {
		return nil, pkg.UserNotExistErr
	}

	appointments, err := service.provider.GetUserAppointments(userId)
	if err != nil {
		return nil, err
	}

	return appointments, nil
}
