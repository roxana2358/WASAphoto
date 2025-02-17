package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"wasa-photo/service/database"
)

var ErrUnauthorized = errors.New("user not in database")

type Username struct {
	Username string `json:"username"`
}

type ID struct {
	Id uint64 `json:"id"`
}

type Comment struct {
	Comment string `json:"comment"`
}

type Userprofile struct {
	Id             uint64   `json:"id"`
	Username       string   `json:"username"`
	Photos         []uint64 `json:"posts"`
	NumberOfPhotos int      `json:"numberOfPhotos"`
	Followers      []string `json:"followers"`
	Following      []string `json:"following"`
	Banned         []string `json:"banned"`
}

type Userpost struct {
	UserID   uint64       `json:"userId"`
	Username string       `json:"username"`
	PostID   uint64       `json:"photoId"`
	Date     string       `json:"date"`
	Time     string       `json:"time"`
	Likes    []uint64     `json:"likes"`
	Comments []CommentOBJ `json:"comments"`
}

type CommentOBJ struct {
	Username  string `json:"username"`
	UserID    uint64 `json:"userId"`
	Comment   string `json:"comment"`
	CommentId uint64 `json:"commentId"`
}

/**
* Gets token from header.
 */
func getHeaderToken(r *http.Request) (uint64, error) {
	// split authorization
	auth := strings.Split(r.Header.Get("Authorization"), " ")
	if len(auth) != 2 {
		// wrong format
		return 0, ErrUnauthorized
	}
	tokenS := auth[1]
	if strings.Compare(tokenS, "") == 0 {
		// token absent
		return 0, ErrUnauthorized
	}
	tokenI, err := strconv.ParseUint(tokenS, 10, 64)
	if err != nil {
		// token extraction failed
		return 0, err
	}
	return tokenI, nil
}

/**
* Checks if the user that requested the action is the same as the one that is supposed to do it.
 */
func checkAuth(srcUser uint64, dstUser uint64) bool {
	return srcUser == dstUser
}

/**
* UserProfileFromDatabase populates the struct with data from the database, overwriting all values.
 */
func (u *Userprofile) UserProfileFromDatabase(userprofile database.Userprofile) {
	u.Id = userprofile.Id
	u.Username = userprofile.Username
	u.Photos = userprofile.Photos
	u.NumberOfPhotos = userprofile.NumberOfPhotos
	u.Followers = userprofile.Followers
	u.Following = userprofile.Following
	u.Banned = userprofile.Banned
}

/**
* UserPostFromDatabase creates a new struct with data from the database.
 */
func NewUserPostFromDatabase(userpost database.Userpost) Userpost {
	var u Userpost
	u.UserID = userpost.UserID
	u.Username = userpost.Username
	u.PostID = userpost.PostID
	u.Date = userpost.Date
	u.Time = userpost.Time
	u.Likes = userpost.Likes
	var com CommentOBJ
	for i := 0; i < len(userpost.Comments); i++ {
		com.CommentFromDatabase(userpost.Comments[i])
		u.Comments = append(u.Comments, com)
	}
	return u
}

/**
* CommentFromDatabase populates the struct with data from the database, overwriting all values.
 */
func (c *CommentOBJ) CommentFromDatabase(comment database.CommentOBJ) {
	c.Username = comment.Username
	c.UserID = comment.UserID
	c.Comment = comment.Comment
	c.CommentId = comment.CommentId
}
