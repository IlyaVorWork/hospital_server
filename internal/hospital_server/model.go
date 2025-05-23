package hospital_server

import "time"

type Spec struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}
type Doctor struct {
	Id         string  `json:"id"`
	LastName   string  `json:"last_name"`
	FirstName  string  `json:"first_name"`
	SecondName *string `json:"second_name"`
	AvatarLink *string `json:"avatar_link"`
	Spec       string  `json:"spec"`
}

type User struct {
	Id                string    `json:"id"`
	Phone             string    `json:"phone"`
	LastName          string    `json:"last_name"`
	FirstName         string    `json:"first_name"`
	SecondName        *string   `json:"second_name"`
	Gender            bool      `json:"gender"`
	BirthDate         time.Time `json:"birth_date"`
	PassportSeries    string    `json:"passport_series"`
	PassportNumber    string    `json:"passport_number"`
	PassportIssuer    string    `json:"passport_issuer"`
	PassportIssueDate time.Time `json:"passport_issue_date"`
	Snils             string    `json:"snils"`
}

type Appointment struct {
	Date      time.Time `json:"date"`
	Time      string    `json:"time"`
	Cabinet   int       `json:"cabinet"`
	Doctor    Doctor    `json:"doctor"`
	PatientId string    `json:"patient_id"`
}

type GetSpecsOut struct {
	Specs []Spec `json:"specs"`
}

type GetDoctorsOut struct {
	Doctors []Doctor `json:"doctors"`
}

type GetSlotsOut struct {
	Slots []Appointment `json:"slots"`
}

type MakeAppointmentDTO struct {
	Date     time.Time `json:"date" validate:"required"`
	Time     string    `json:"time" validate:"required"`
	DoctorId string    `json:"doctor_id" validate:"required"`
	UserId   string    `json:"user_id" validate:"required"`
}

type CancelAppointmentDTO struct {
	Date     time.Time `json:"date" validate:"required"`
	Time     string    `json:"time" validate:"required"`
	DoctorId string    `json:"doctor_id" validate:"required"`
	UserId   string    `json:"user_id" validate:"required"`
}

type GetUserProfileDTO struct {
	UserId string `json:"user_id" validate:"required"`
}

type GetUserProfileOut struct {
	Profile User `json:"profile"`
}

type GetUserAppointmentsDTO struct {
	UserId string `json:"user_id" validate:"required"`
}

type GetUserAppointmentsOut struct {
	Appointments []Appointment `json:"appointments"`
}
