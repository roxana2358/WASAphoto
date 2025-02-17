package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"wasa-photo/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

/**
* It returns a stream with following users' photos and respective information in reverse
* chronological order.
 */
func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// token extraction
	token, err := getHeaderToken(r)
	if errors.Is(err, ErrUnauthorized) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		// error on our side: log the error and send a 500 to the user
		ctx.Logger.WithError(err).Error("can't extract token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// The user ID in the path is a 64-bit unsigned integer
	userID, err := strconv.ParseUint(ps.ByName("userID"), 10, 64)
	if err != nil {
		// The value was not uint64, reject the action indicating an error on the client side
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if the user that requested the action is authorized
	if !checkAuth(token, userID) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dbUserStream, err := rt.db.GetUserStream(userID) // list of userPost in db format
	// check errors
	if err != nil {
		// error on our side: log the error and send a 500 to the user
		ctx.Logger.WithField("user", userID).WithError(err).Error("can't get stream")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// send output to user
	var userStream []Userpost
	for i := 0; i < len(dbUserStream); i++ {
		userStream = append(userStream, NewUserPostFromDatabase(dbUserStream[i]))
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(userStream)
}
