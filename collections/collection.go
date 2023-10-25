package collections

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
)

type Collection struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Collection) Create(db *sql.DB) error {
	c.ID = ulid.Make().String()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	_, err := db.Exec("INSERT INTO collections (id, name, created_at, updated_at) VALUES (?, ?, ?, ?)", c.ID, c.Name, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func GetAllCollections(db *sql.DB) ([]Collection, error) {
	rows, err := db.Query("SELECT * FROM collections")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []Collection
	for rows.Next() {
		var c Collection
		err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		collections = append(collections, c)
	}

	return collections, nil
}

func GetCollectionByID(db *sql.DB, id string) (*Collection, error) {
	row := db.QueryRow("SELECT * FROM collections WHERE id = ?", id)

	var c Collection
	err := row.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("collection with ID %s not found", id)
		}
		return nil, err
	}

	return &c, nil
}

func (c *Collection) Update(db *sql.DB, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	var placeholders []string
	var values []interface{}
	for k, v := range updates {
		placeholders = append(placeholders, fmt.Sprintf("%s = ?", k))
		values = append(values, v)
	}
	values = append(values, c.ID)

	query := fmt.Sprintf("UPDATE collections SET %s WHERE id = ?", joinStrings(placeholders, ", "))
	_, err := db.Exec(query, values...)
	if err != nil {
		return err
	}

	err = db.QueryRow("SELECT * FROM collections WHERE id = ?", c.ID).Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (c *Collection) Delete(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM collections WHERE id = ?", c.ID)
	if err != nil {
		return err
	}

	return nil
}

func GetTagsForCollection(db *sql.DB, collectionID string) ([]Tag, error) {
	rows, err := db.Query("SELECT * FROM tags WHERE collection_id = ?", collectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var t Tag
		err := rows.Scan(&t.ID, &t.Name, &t.CollectionID, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	return tags, nil
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}
	return strs[0] + sep + joinStrings(strs[1:], sep)
}
