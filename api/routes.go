package api

import (
	"net/http"

	"github.com/go-chi/chi"
)

// SetRoutes sets up the routes for the API endpoints using the chi router.
func SetRoutes(r *chi.Mux) http.Handler {

	// Create a new collection
	r.Post("/collections", CreateCollectionHandler)

	// Get data points from a collection
	r.Get("/collections/{collectionName}/datapoints", GetDataPointsHandler)

	// Update a tag
	r.Put("/collections/{collectionName}/tags/{tagName}", UpdateTagHandler)

	// Delete a tag
	r.Delete("/collections/{collectionName}/tags/{tagName}", DeleteTagHandler)

	// Update a collection
	r.Put("/collections/{collectionName}", UpdateCollectionHandler)

	// Delete a collection
	r.Delete("/collections/{collectionName}", DeleteCollectionHandler)

	// Get tags under a collection
	r.Get("/collections/{collectionName}/tags", GetTagsHandler)

	// Get data points under a tag
	r.Get("/collections/{collectionName}/tags/{tagName}/datapoints", GetDataPointsByTagHandler)

	return r
}
