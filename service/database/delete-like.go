package database

/**
* Removes like from database; returns error if the request is unsuccessfull.
 */
func (db *appdbimpl) DeleteLike(userID uint64, postID uint64) error {
	// unlike photo
	res, err := db.c.Exec(`DELETE FROM Likes WHERE PostId=? AND UserId=?`, postID, userID)
	if err != nil {
		return err
	}

	// check if it affected the database
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		return ErrLikeNotFound
	}
	return nil
}
