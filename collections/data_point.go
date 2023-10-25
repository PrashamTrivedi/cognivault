package collections

import (
	"database/sql"
	"fmt"
	"log"
)

type DataPoint struct {
	ID    string `json:"id"`
	TagID string `json:"tag_id"`
	Value string `json:"value"`
	db    *sql.DB
}

func (dp *DataPoint) Create() error {
	query := fmt.Sprintf("INSERT INTO data_points (tag_id, value) VALUES (%d, '%s')", dp.TagID, dp.Value)
	_, err := dp.db.Exec(query)
	if err != nil {
		log.Printf("Error creating data point: %v", err)
		return err
	}
	return nil
}

func (dp *DataPoint) Update() error {
	query := fmt.Sprintf("UPDATE data_points SET tag_id=%d, value='%s' WHERE id=%d", dp.TagID, dp.Value, dp.ID)
	_, err := dp.db.Exec(query)
	if err != nil {
		log.Printf("Error updating data point: %v", err)
		return err
	}
	return nil
}

func (dp *DataPoint) Delete() error {
	query := fmt.Sprintf("DELETE FROM data_points WHERE id=%d", dp.ID)
	_, err := dp.db.Exec(query)
	if err != nil {
		log.Printf("Error deleting data point: %v", err)
		return err
	}
	return nil
}

func GetDataPointsByTagID(db *sql.DB, tagID int) ([]DataPoint, error) {
	query := fmt.Sprintf("SELECT id, tag_id, value FROM data_points WHERE tag_id=%d", tagID)
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error getting data points by tag ID: %v", err)
		return nil, err
	}
	defer rows.Close()

	dataPoints := []DataPoint{}
	for rows.Next() {
		var dp DataPoint
		err := rows.Scan(&dp.ID, &dp.TagID, &dp.Value)
		if err != nil {
			log.Printf("Error scanning data point row: %v", err)
			return nil, err
		}
		dp.db = db
		dataPoints = append(dataPoints, dp)
	}
	return dataPoints, nil
}
