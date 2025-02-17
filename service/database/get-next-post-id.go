package database

import (
	"database/sql"
	"errors"
)

/**
* Gets the id for the new post.
 */
func (db *appdbimpl) GetNextPostId() (uint64, error) {
	// get next id
	row := db.c.QueryRow(`SELECT PostId FROM Posts WHERE PostId=(SELECT max(PostId) FROM Posts)`)
	var postId uint64
	err := row.Scan(&postId)
	if errors.Is(err, sql.ErrNoRows) {
		postId = 0
	} else if err != nil {
		return postId, err
	}
	postId++
	return postId, nil
}
