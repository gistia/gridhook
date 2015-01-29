package models

import (
	"database/sql"
	"time"
)

type EmailEventStatus struct {
	Id           sql.NullInt64
	EmailEventId sql.NullInt64
	Timestamp    time.Time
	Payload      string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func FindEmailEventStatus(db *sql.DB, id int64) (*EmailEventStatus, error) {
	statuses, err := RunQuery(db, `
    SELECT
      id, email_event_id, timestamp, payload, status, created_at, updated_at
    FROM
      email_event_statuses
    WHERE id = ?`, id)

	if err != nil {
		return nil, err
	}

	if len(statuses) > 0 {
		return statuses[0], err
	}

	return nil, err
}

func FindEmailEventStatuses(db *sql.DB, emailEventId int64) ([]*EmailEventStatus, error) {
	return RunQuery(db, `
    SELECT
      id, email_event_id, timestamp, payload, status, created_at, updated_at
    FROM
      email_event_statuses
    WHERE email_event_id = ?`, emailEventId)
}

func FindLatestEmailEventStatus(db *sql.DB, emailEventId int64) (*EmailEventStatus, error) {
	statuses, err := RunQuery(db, `
    SELECT
      id, email_event_id, timestamp, payload, status, created_at, updated_at
    FROM
      email_event_statuses
    WHERE
      email_event_id = ?
    ORDER BY
      timestamp DESC
    LIMIT 1`, emailEventId)

	if err != nil {
		return nil, err
	}

	if len(statuses) > 0 {
		return statuses[0], err
	}

	return nil, err
}

func (s *EmailEventStatus) Insert(db *sql.DB) (*EmailEventStatus, error) {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()

	res, err := db.Exec(`
    INSERT INTO email_event_statuses
    (email_event_id, timestamp, payload, status, created_at, updated_at)
    VALUES (?, ?, ?, ?, ?, ?)
  `, s.EmailEventId, s.Timestamp, s.Payload, s.Status, s.CreatedAt, s.UpdatedAt)

	if err != nil {
		return s, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return s, err
	}

	status, err := FindEmailEventStatus(db, id)
	return status, err
}

func RunQuery(db *sql.DB, query string, id int64) ([]*EmailEventStatus, error) {
	var statuses []*EmailEventStatus

	rows, err := db.Query(query, id)

	if err != nil {
		return statuses, err
	}

	defer rows.Close()

	for rows.Next() {
		status, err := FillData(rows)
		if err != nil {
			return statuses, err
		}
		statuses = append(statuses, status)
	}

	return statuses, err
}

func FillData(rows *sql.Rows) (*EmailEventStatus, error) {
	var result *EmailEventStatus
	var id, email_event_id sql.NullInt64
	var timestamp, created_at, updated_at time.Time
	var payload, status []byte

	err := rows.Scan(&id, &email_event_id, &timestamp, &payload, &status,
		&created_at, &updated_at)

	if err != nil {
		return result, err
	}

	result = &EmailEventStatus{
		Id:           id,
		EmailEventId: email_event_id,
		Timestamp:    timestamp,
		Payload:      string(payload),
		Status:       string(status),
		CreatedAt:    created_at,
		UpdatedAt:    updated_at,
	}

	return result, err
}
