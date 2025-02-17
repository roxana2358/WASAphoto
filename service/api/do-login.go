package api

import (
	"encoding/json"
	"net/http"
	"wasa-photo/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

/**
* If the user does not exist, it will be created, and an identifier is returned. If the user exists,
* the user identifier is returned.
 */
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// read the username from the request body
	var username Username
	err := json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		// the body was not a parseable JSON, reject it
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := rt.db.GetUserId(username.Username)
	// if an error occurred the user is not in database, it must be added
	if err != nil {
		id, err = rt.db.CreateUser(username.Username)
		if err != nil {
			// error on our side: log the error and send a 500 to the user
			ctx.Logger.WithError(err).Error("can't create the user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// send output to user
	w.Header().Set("Content-Type", "application/json")
	var identifier ID
	identifier.Id = id
	_ = json.NewEncoder(w).Encode(identifier)
}
