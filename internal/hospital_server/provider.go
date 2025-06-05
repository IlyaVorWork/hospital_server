package hospital_server

import (
	"database/sql"
	"fmt"
	"time"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (provider *PostgresRepository) GetSpecs() ([]Spec, error) {

	var specs []Spec

	rows, err := provider.db.Query("SELECT s.id, s.title, count(a.doctor_id) AS slots_count FROM specialization s JOIN public.doctor d on s.id = d.specialization_id JOIN public.appointment a on d.id = a.doctor_id AND (a.date > CURRENT_DATE OR (a.date = CURRENT_DATE AND a.time > to_char(now(), 'HH24:MI'))) AND a.profile_id IS NULL GROUP BY s.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var spec Spec
		if err := rows.Scan(&spec.Id, &spec.Title, &spec.SlotsCount); err != nil {
			return nil, err
		}
		specs = append(specs, spec)
	}

	return specs, nil
}

func (provider *PostgresRepository) GetDoctors(specId string) ([]Doctor, error) {

	var doctors []Doctor

	_, err := provider.db.Exec("SET TIMEZONE = 'Asia/Novosibirsk'")
	if err != nil {
		return nil, err
	}

	rows, err := provider.db.Query("SELECT d.id, d.last_name, d.first_name, d.second_name, d.avatar_link, COUNT(a.doctor_id) as slots_count FROM doctor d LEFT JOIN appointment a ON a.doctor_id = d.id AND (a.date > CURRENT_DATE OR (a.date = CURRENT_DATE AND a.time > to_char(CURRENT_TIMESTAMP, 'HH24:MI'))) AND a.profile_id IS NULL WHERE d.specialization_id = $1 GROUP BY d.id, d.last_name, d.first_name, d.second_name, d.avatar_link", specId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var doctor Doctor
		if err := rows.Scan(&doctor.Id, &doctor.LastName, &doctor.FirstName, &doctor.SecondName, &doctor.AvatarLink, &doctor.SlotsCount); err != nil {
			return nil, err
		}
		doctors = append(doctors, doctor)
	}

	return doctors, nil
}

func (provider *PostgresRepository) GetSlots(doctorId string) ([]Appointment, error) {

	var doctor Doctor
	var slots []Appointment

	_, err := provider.db.Exec("SET TIMEZONE = 'Asia/Novosibirsk'")
	if err != nil {
		return nil, err
	}

	err = provider.db.QueryRow("SELECT d.id, d.last_name, d.first_name, d.second_name, d.avatar_link FROM doctor d JOIN specialization s ON s.id = d.specialization_id WHERE d.id = $1", doctorId).Scan(&doctor.Id, &doctor.LastName, &doctor.FirstName, &doctor.SecondName, &doctor.AvatarLink)
	if err != nil {
		return nil, err
	}

	rows, err := provider.db.Query("SELECT a.date, a.time, c.number, a.profile_id FROM appointment a JOIN cabinet c ON c.id = a.cabinet_id WHERE a.doctor_id = $1 AND a.profile_id IS NULL AND (a.date > CURRENT_DATE OR (a.date = CURRENT_DATE AND a.time > to_char(CURRENT_TIMESTAMP, 'HH24:MI')))", doctorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	fmt.Println("АБОБА 3")
	for rows.Next() {
		var slot Appointment
		if err := rows.Scan(&slot.Date, &slot.Time, &slot.Cabinet, &slot.PatientId); err != nil {
			return nil, err
		}
		fmt.Println(slot)
		slot.Doctor = doctor
		slots = append(slots, slot)
	}

	return slots, nil
}

func (provider *PostgresRepository) MakeAppointment(date time.Time, time, doctorId, userId string) error {

	_, err := provider.db.Exec("UPDATE appointment SET profile_id = $1 WHERE date = $2 AND time = $3 AND doctor_id = $4", userId, date, time, doctorId)
	if err != nil {
		return err
	}

	return nil
}

func (provider *PostgresRepository) CancelAppointment(date time.Time, time, doctorId string) error {

	_, err := provider.db.Exec("UPDATE appointment SET profile_id = NULL WHERE date = $1 AND time = $2 AND doctor_id = $3", date, time, doctorId)
	if err != nil {
		return err
	}

	return nil
}

func (provider *PostgresRepository) GetUserProfile(userId string) (*User, error) {

	var profile User

	err := provider.db.QueryRow("SELECT id, phone, last_name, first_name, second_name, gender, birth_date, passport_series, passport_number, passport_issuer, passport_issue_date, snils FROM profile WHERE id = $1", userId).Scan(&profile.Id, &profile.Phone, &profile.LastName, &profile.FirstName, &profile.SecondName, &profile.Gender, &profile.BirthDate, &profile.PassportSeries, &profile.PassportNumber, &profile.PassportIssuer, &profile.PassportIssueDate, &profile.Snils)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (provider *PostgresRepository) GetUserAppointments(userId string) ([]Appointment, error) {

	var appointments []Appointment

	rows, err := provider.db.Query("SELECT a.date, a.time, c.number, a.doctor_id, a.profile_id FROM appointment a JOIN cabinet c ON c.id = a.cabinet_id WHERE a.profile_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var appointment Appointment
		var doctor Doctor
		if err := rows.Scan(&appointment.Date, &appointment.Time, &appointment.Cabinet, &doctor.Id, &appointment.PatientId); err != nil {
			return nil, err
		}

		err := provider.db.QueryRow("SELECT d.id, d.last_name, d.first_name, d.second_name, d.avatar_link, s.title FROM doctor d JOIN specialization s ON s.id = d.specialization_id WHERE d.id = $1", doctor.Id).Scan(&doctor.Id, &doctor.LastName, &doctor.FirstName, &doctor.SecondName, &doctor.AvatarLink, &doctor.Spec)
		if err != nil {
			return nil, err
		}

		appointment.Doctor = doctor
		appointments = append(appointments, appointment)
	}

	return appointments, nil
}

func (provider *PostgresRepository) DoesUserExist(userId string) bool {

	var id string

	err := provider.db.QueryRow("SELECT id FROM profile WHERE id = $1", userId).Scan(&id)

	return err == nil
}

func (provider *PostgresRepository) GetAppointment(date time.Time, time, doctorId string) (*Appointment, error) {

	var appointment Appointment
	err := provider.db.QueryRow("SELECT a.date, a.time, c.number, a.profile_id FROM appointment a JOIN cabinet c ON c.id = a.cabinet_id WHERE a.date = $1 AND a.time = $2 AND a.doctor_id = $3", date, time, doctorId).Scan(&appointment.Date, &appointment.Time, &appointment.Cabinet, &appointment.PatientId)
	if err != nil {
		return nil, err
	}

	return &appointment, nil
}
