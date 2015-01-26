package models

import (
	"database/sql"
	"time"
)

type EmailEvent struct {
	Id          sql.NullInt64
	TargetType  string
	TargetId    sql.NullInt64
	UserId      sql.NullInt64
	MassEmailId sql.NullInt64
	TemplateId  sql.NullInt64
	UniqueId    string
	From        string
	To          string
	Subject     string
	Body        string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (e *EmailEvent) Update(db *sql.DB) error {
	stmt, err := db.Prepare(`
		UPDATE
			email_events
		SET
			target_type = ?,
			target_id = ?,
			user_id = ?,
			mass_email_id = ?,
			template_id = ?,
			unique_id = ?,
			` + "`from`" + ` = ?,
			` + "`to`" + ` = ?,
			subject = ?,
			body = ?,
			status = ?,
			created_at = ?,
			updated_at = ?
		WHERE
			id = ?`)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(e.TargetType, e.TargetId, e.UserId, e.MassEmailId,
		e.TemplateId, e.UniqueId, e.From, e.To, e.Subject, e.Body, e.Status,
		e.CreatedAt, e.UpdatedAt, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func FindEmailEvent(db *sql.DB, id int) (*EmailEvent, error) {
	var e *EmailEvent

	rows, err := db.Query(`
		SELECT
			id, target_type, target_id, user_id, mass_email_id, template_id,
			unique_id, `+"`from`"+`, `+"`to`"+`, subject, body,
			status, created_at, updated_at
		FROM
			email_events
		WHERE id = ?`, id)

	if err != nil {
		return e, err
	}

	rows.Next()

	e, err = NewEmailEvent(rows)
	if err != nil {
		return e, err
	}

	return e, err
}

func FindEmailEventByUniqueId(db *sql.DB, uniqueId string) (*EmailEvent, error) {
	var e *EmailEvent

	rows, err := db.Query(`
		SELECT
			id, target_type, target_id, user_id, mass_email_id, template_id,
			unique_id, `+"`from`"+`, `+"`to`"+`, subject, body,
			status, created_at, updated_at
		FROM
			email_events
		WHERE
			unique_id = ?`, uniqueId)

	if err != nil {
		return e, err
	}

	if !rows.Next() {
		return nil, nil
	}

	e, err = NewEmailEvent(rows)
	if err != nil {
		return e, err
	}

	return e, err
}

func NewEmailEvent(rows *sql.Rows) (*EmailEvent, error) {
	var emailEvent *EmailEvent
	var id, target_id, user_id, mass_email_id, template_id sql.NullInt64
	var unique_id, target_type, from, to, subject, body, status []byte
	var created_at, updated_at time.Time

	err := rows.Scan(&id, &target_type, &target_id, &user_id, &mass_email_id,
		&template_id, &unique_id, &from, &to, &subject, &body, &status, &created_at,
		&updated_at)

	if err != nil {
		return emailEvent, err
	}

	emailEvent = &EmailEvent{
		Id:          id,
		TargetType:  string(target_type),
		TargetId:    target_id,
		UserId:      user_id,
		MassEmailId: mass_email_id,
		TemplateId:  template_id,
		UniqueId:    string(unique_id),
		From:        string(from),
		To:          string(to),
		Subject:     string(subject),
		Body:        string(body),
		Status:      string(status),
		CreatedAt:   created_at,
		UpdatedAt:   updated_at,
	}

	return emailEvent, err
}
