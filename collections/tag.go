package collections

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

// Tag represents a tag in the database
type Tag struct {
	ID           string `json:"id"`
	CollectionID string `json:"collection_id"`
	Name         string `json:"name"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// CreateTag creates a new tag in the database
func (t *Tag) CreateTag(db *sql.DB) error {
	result, err := db.Exec("INSERT INTO tags(collection_id, name) VALUES(?,?, ?)", t.ID, t.CollectionID, t.Name)
	if err != nil {
		log.Println(err)
		return errors.New("failed to create tag")
	}
	return nil
}

// GetTagsByCollectionID gets all tags under a collection
func GetTagsByCollectionID(db *sql.DB, collectionID int) ([]Tag, error) {
	rows, err := db.Query("SELECT id, collection_id, name FROM tags WHERE collection_id=?", collectionID)
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to get tags")
	}
	defer rows.Close()

	tags := []Tag{}
	for rows.Next() {
		var tag Tag
		err := rows.Scan(&tag.ID, &tag.CollectionID, &tag.Name)
		if err != nil {
			log.Println(err)
			return nil, errors.New("failed to get tags")
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// UpdateTag updates a tag in the database
func UpdateTag(db *sql.DB, tagID int, name string) error {
	result, err := db.Exec("UPDATE tags SET name=? WHERE id=?", name, tagID)
	if err != nil {
		log.Println(err)
		return errors.New("failed to update tag")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return errors.New("failed to update tag")
	}
	if rowsAffected == 0 {
		return fmt.Errorf("tag with id %d not found", tagID)
	}
	return nil
}

// DeleteTag deletes a tag and all data points under it from the database
func DeleteTag(db *sql.DB, tagID int) error {
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return errors.New("failed to delete tag")
	}

	_, err = tx.Exec("DELETE FROM data_points WHERE tag_id=?", tagID)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return errors.New("failed to delete tag")
	}

	result, err := tx.Exec("DELETE FROM tags WHERE id=?", tagID)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return errors.New("failed to delete tag")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return errors.New("failed to delete tag")
	}
	if rowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("tag with id %d not found", tagID)
	}

	return tx.Commit()
}
