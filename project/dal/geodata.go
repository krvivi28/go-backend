package dao

import (
	"GOLANG/project/models"
	"GOLANG/project/utils"
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type GeoDataDao struct {
	db *sql.DB
}

// InitializeDB initializes the database and creates the users table if it doesn't exist.
func InitializeGeoDB() (*GeoDataDao, error) {
	var err error
	dbInstance, err := sql.Open("sqlite3", "geodata.db")
	if err != nil {
		return nil, err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS geodata (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        longlatdata INTEGER[][],
		email TEXT NOT NULL
    );`
	_, err = dbInstance.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return &GeoDataDao{db: dbInstance}, nil
}

func (s *GeoDataDao) Upload(ctx context.Context, geoData *models.GeoData) error {
	arrayStr := utils.ConvertToPostgresArray(geoData.LongLatData)
	_, err := s.db.Exec("INSERT INTO geodata (name, longlatdata, email) VALUES (?, ?, ?)",
		geoData.Name, arrayStr, geoData.Email)
	if err != nil {
		return err
	}

	return nil
}

func (s *GeoDataDao) GetEmailFromId(ctx context.Context, id string) (*string, error) {
	var email string
	err := s.db.QueryRow("SELECT email FROM geodata WHERE id = ?", id).Scan(&email)
	if err != nil {
		return nil, err
	}

	return &email, nil
}

func (s *GeoDataDao) List(ctx context.Context, email string) (*sql.Rows, error) {
	results, err := s.db.QueryContext(ctx, "SELECT * FROM geodata WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (s *GeoDataDao) Patch(ctx context.Context, geoId string, geoData *models.GeoData) error {
	arrayStr := utils.ConvertToPostgresArray(geoData.LongLatData)
	_, err := s.db.Exec("UPDATE geodata SET name = ?, longlatData = ? WHERE id = ?",
		geoData.Name, arrayStr, geoId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
