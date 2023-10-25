package api

import (
	"cognivaultServer/collections"
	"cognivaultServer/database"
	"cognivaultServer/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/oklog/ulid/v2"
)

// CreateCollectionRequest represents the request body for creating a new collection.
type CreateCollectionRequest struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
	Text string `json:"text,omitempty"`
	File string `json:"file,omitempty"`
	Tag  string `json:"tag,omitempty"`
}

// CreateCollectionResponse represents the response body for creating a new collection.
type CreateCollectionResponse struct {
	ID string `json:"id"`
}

// GetCollectionRequest represents the request parameters for getting data points from a collection.
type GetCollectionRequest struct {
	CollectionName string `json:"collection_name"`
	Query          string `json:"query"`
}

// GetCollectionResponse represents the response body for getting data points from a collection.
type GetCollectionResponse struct {
	DataPoints []collections.DataPoint `json:"data_points"`
}

// UpdateTagRequest represents the request body for updating a tag.
type UpdateTagRequest struct {
	CollectionName string `json:"collection_name"`
	TagID          string `json:"tag_id"`
	NewTag         string `json:"new_tag"`
}

// UpdateCollectionRequest represents the request body for updating a collection.
type UpdateCollectionRequest struct {
	CollectionName string `json:"collection_name"`
	NewName        string `json:"new_name"`
}

// DeleteTagRequest represents the request parameters for deleting a tag.
type DeleteTagRequest struct {
	CollectionName string `json:"collection_name"`
	TagID          string `json:"tag_id"`
}

// DeleteCollectionRequest represents the request parameters for deleting a collection.
type DeleteCollectionRequest struct {
	CollectionName string `json:"collection_name"`
}

// GetTagsRequest represents the request parameters for getting tags under a collection.
type GetTagsRequest struct {
	CollectionName string `json:"collection_name"`
}

// GetTagsResponse represents the response body for getting tags under a collection.
type GetTagsResponse struct {
	Tags []collections.Tag `json:"tags"`
}

// GetDataPointsRequest represents the request parameters for getting data points under a tag.
type GetDataPointsRequest struct {
	CollectionName string `json:"collection_name"`
	TagID          string `json:"tag_id"`
}

// GetDataPointsResponse represents the response body for getting data points under a tag.
type GetDataPointsResponse struct {
	DataPoints []collections.DataPoint `json:"data_points"`
}

// CreateCollectionHandler handles the HTTP request for creating a new collection.
func CreateCollectionHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateCollectionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var data string
	switch {
	case req.URL != "":
		data, err = utils.FetchURL(req.URL)
		if err != nil {
			utils.SendResponse(w, http.StatusBadRequest, "Failed to fetch URL")
			return
		}
	case req.File != "":
		data, err = utils.ReadFile(req.File)
		if err != nil {
			utils.SendResponse(w, http.StatusBadRequest, "Failed to read file")
			return
		}
	case req.Text != "":
		data = req.Text
	default:
		utils.SendResponse(w, http.StatusBadRequest, "Missing data source")
		return
	}

	tag := req.Tag
	if tag == "" {
		tag = req.URL
	}

	collection := collections.Collection{
		ID:   ulid.Make().String(),
		Name: req.Name,
	}
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = collection.Create(db)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to create collection")
		return
	}

	tagObj := collections.Tag{
		ID:           ulid.Make().String(),
		Name:         tag,
		CollectionID: collection.ID,
	}
	err = tagObj.CreateTag(db)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to create tag")
		return
	}

	dataPoint := collections.DataPoint{
		TagID: tagObj.ID,
		Value: data,
	}
	err = dataPoint.Create()
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to create data point")
		return
	}

	resp := CreateCollectionResponse{
		ID: collection.ID,
	}
	render.JSON(w, r, resp)
}

// GetCollectionHandler handles the HTTP request for getting data points from a collection.
func GetCollectionHandler(w http.ResponseWriter, r *http.Request) {
	var req GetCollectionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	collection := collections.Collection{
		Name: req.CollectionName,
	}
	err = collection.Get()
	if err != nil {
		utils.SendResponse(w, http.StatusNotFound, "Collection not found")
		return
	}

	dataPoints, err := collection.GetDataPoints(req.Query)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to get data points")
		return
	}

	resp := GetCollectionResponse{
		DataPoints: dataPoints,
	}
	render.JSON(w, r, resp)
}

// UpdateTagHandler handles the HTTP request for updating a tag.
func UpdateTagHandler(w http.ResponseWriter, r *http.Request) {
	var req UpdateTagRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	tag := collections.Tag{
		ID:           req.TagID,
		CollectionID: collections.GetCollectionID(req.CollectionName),
		Name:         req.NewTag,
	}
	err = tag.Update()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update tag")
		return
	}

	utils.RespondWithMessage(w, http.StatusOK, "Tag updated successfully")
}

// UpdateCollectionHandler handles the HTTP request for updating a collection.
func UpdateCollectionHandler(w http.ResponseWriter, r *http.Request) {
	var req UpdateCollectionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	collection := collections.Collection{
		Name:    req.CollectionName,
		NewName: req.NewName,
	}
	err = collection.Update()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update collection")
		return
	}

	utils.RespondWithMessage(w, http.StatusOK, "Collection updated successfully")
}

// DeleteTagHandler handles the HTTP request for deleting a tag.
func DeleteTagHandler(w http.ResponseWriter, r *http.Request) {
	tagID := chi.URLParam(r, "tagID")
	collectionName := chi.URLParam(r, "collectionName")

	tag := collections.Tag{
		ID:           tagID,
		CollectionID: collections.GetCollectionID(collectionName),
	}
	err := tag.Delete()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete tag")
		return
	}

	utils.RespondWithMessage(w, http.StatusOK, "Tag deleted successfully")
}

// DeleteCollectionHandler handles the HTTP request for deleting a collection.
func DeleteCollectionHandler(w http.ResponseWriter, r *http.Request) {
	collectionName := chi.URLParam(r, "collectionName")

	collection := collections.Collection{
		Name: collectionName,
	}
	err := collection.Delete()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete collection")
		return
	}

	utils.RespondWithMessage(w, http.StatusOK, "Collection deleted successfully")
}

// GetTagsHandler handles the HTTP request for getting tags under a collection.
func GetTagsHandler(w http.ResponseWriter, r *http.Request) {
	collectionName := chi.URLParam(r, "collectionName")

	collection := collections.Collection{
		Name: collectionName,
	}
	err := collection.Get()
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Collection not found")
		return
	}

	tags, err := collection.GetTags()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get tags")
		return
	}

	resp := GetTagsResponse{
		Tags: tags,
	}
	render.JSON(w, r, resp)
}

// GetDataPointsHandler handles the HTTP request for getting data points under a tag.
func GetDataPointsHandler(w http.ResponseWriter, r *http.Request) {
	tagID := chi.URLParam(r, "tagID")
	collectionName := chi.URLParam(r, "collectionName")

	tag := collections.Tag{
		ID:           tagID,
		CollectionID: collections.GetCollectionID(collectionName),
	}
	err := tag.Get()
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Tag not found")
		return
	}

	dataPoints, err := tag.GetDataPoints()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get data points")
		return
	}

	resp := GetDataPointsResponse{
		DataPoints: dataPoints,
	}
	render.JSON(w, r, resp)
}
