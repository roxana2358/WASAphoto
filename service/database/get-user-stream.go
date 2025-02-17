package database

import (
	"database/sql"
	"errors"
	"sort"
)

/**
* Gets the stream with 30 post from following in reverse chronological order.
 */
func (db *appdbimpl) GetUserStream(userID uint64) ([]Userpost, error) {
	var userStream []Userpost
	var userPost Userpost

	// request posts from following
	rows, err := db.c.Query(`SELECT Users.Id, Users.Username, Posts.PostId, Posts.Date, Posts.Time
							FROM Following 
							INNER JOIN Posts ON Following.FollowingId=Posts.UserId
							INNER JOIN Users ON Following.FollowingId=Users.Id
							WHERE Following.UserId=?`, userID)
	if err != nil {
		return userStream, err
	}
	defer func() { _ = rows.Close() }()

	// here we read the resultset and we build the list to be put in userStream
	var likeId uint64
	var likes []uint64
	var comment CommentOBJ
	var comments []CommentOBJ
	for rows.Next() {
		err = rows.Scan(&userPost.UserID, &userPost.Username, &userPost.PostID, &userPost.Date, &userPost.Time)
		if err != nil {
			return userStream, err
		}

		// get likes
		l, err := db.c.Query(`SELECT Likes.UserId 
							FROM Posts 
							INNER JOIN Likes ON Posts.PostId=Likes.PostId 
							WHERE Posts.PostId=?`, userPost.PostID)
		if errors.Is(err, sql.ErrNoRows) {
			// no likes
			userPost.Likes = nil
		} else if err == nil {
			// likes
			for l.Next() {
				e := l.Scan(&likeId)
				if e != nil {
					return userStream, err
				}
				likes = append(likes, likeId)
			}
			if err = l.Err(); err != nil {
				return userStream, err
			}
			userPost.Likes = likes
			likes = nil
		} else if err != nil {
			// other error
			return userStream, err
		}

		// get comments
		c, err := db.c.Query(`SELECT Users.Username, Comments.UserId, Comments.Text, Comments.CommentId
							FROM Comments
							INNER JOIN Posts ON Posts.PostId=Comments.PostId
							INNER JOIN Users ON Comments.UserId=Users.Id 
							WHERE Comments.PostId=?`, userPost.PostID)
		if errors.Is(err, sql.ErrNoRows) {
			// no comments
			userPost.Comments = nil
		} else if err == nil {
			// comments
			for c.Next() {
				err = c.Scan(&comment.Username, &comment.UserID, &comment.Comment, &comment.CommentId)
				if err != nil {
					return userStream, err
				}
				comments = append(comments, comment)
			}
			if err = c.Err(); err != nil {
				return userStream, err
			}
			userPost.Comments = comments
			comments = nil
		} else if err != nil {
			// other error
			return userStream, err
		}

		// add userPost to userStream
		userStream = append(userStream, userPost)
	}
	if err = rows.Err(); err != nil {
		return userStream, err
	}

	// sort and select
	sort.Sort(postList(userStream))
	var finalList []Userpost
	if len(userStream) > 30 {
		finalList = userStream[:30]
	} else {
		finalList = userStream[:]
	}
	return finalList, nil
}
