package repository

import (
	"go-herder/internal/models"
	"log"
)

type ProcessData struct {
	ID     int
	Label  *string
	Params string
}

func (r *Repository) IterProcesses() chan ProcessData {
	rows, err := r.db.Query("SELECT id, label, params FROM processes")
	if err != nil {
		return nil
	}
	ch := make(chan ProcessData)
	go func() {
		for rows.Next() {
			var pd ProcessData
			err = rows.Scan(&pd.ID, &pd.Label, &pd.Params)
			if err != nil {
				log.Println("error on IterProcesses():", err.Error())
				break
			}
			ch <- pd
		}
		close(ch)
	}()
	return ch
}

func (r *Repository) CreateSession(ip, userAgent string) (s *models.Session, err error) {
	sessionID := uuid
	return
}
func (r *Repository) GetSession(sessionID string) (s *models.Session, err error) {
	row := r.db.QueryRow("SELECT id, ip, user_agent, created_at FROM sessions WHERE id=? AND deleted_at IS NULL LIMIT 1", sessionID)
	if row.Err() != nil {
		return nil, row.Err()
	}
	s = new(models.Session)
	err = row.Scan(&s.ID, &s.IP, &s.UserAgent, &s.CreatedAt)
	return
}
func (r *Repository) DeleteSession(sessionID string) error {
	_, err := r.db.Exec(`UPDATE sessions SET deleted_at=CURRENT_TIMESTAMP WHERE id=?`, sessionID)
	return err
}
