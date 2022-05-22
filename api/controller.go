package api

import (
	"encoding/json"
	"net/http"
)

type Controller struct {
	auth Authenticator
	db   Database
}

func NewController(auth Authenticator, db Database) *Controller {
	return &Controller{
		auth: auth,
		db:   db,
	}
}

func (c *Controller) CreateTODO(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method is not POST", http.StatusMethodNotAllowed)
		return
	}

	// Authentication
	token := r.Header.Get("AuthToken")
	if !c.auth.IsAuthenticated(token) {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}

	// Decoding and validation
	todo := &TODO{}
	if err := json.NewDecoder(r.Body).Decode(todo); err != nil {
		http.Error(w, "invalid json: "+err.Error(), http.StatusBadRequest)
		return
	}
	if err := todo.Validate(); err != nil {
		http.Error(w, "invalid todo: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Database call
	if err := c.db.CreateTODO(todo); err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	respond(w, 201, "todo created")
}
