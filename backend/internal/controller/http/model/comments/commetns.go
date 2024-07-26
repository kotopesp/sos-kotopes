package comments

import "time"

type Comments struct {
	ID         int       `json:"id" db:"id"`
	Content    string    `json:"content" db:"content"`
	Author_id  int       `json:"autorID" db:"author_id"`
	Posts_id   int       `json:"postsID" db:"posts_id"`
	Is_deleted bool      `json:"isDELETED" db:"is_deleted"`
	Deleted_at time.Time `json:"deletedAT" db:"deleted_at"`
	Created_at time.Time `json:"createdAT" db:"created_at"`
	Updated_at time.Time `json:"updatedAT" db:"updated_at"`
	Parent_id  int       `json:"parent_id " db:"parent_id"`
	Reply_id   int       `json:"reply_id" db:"reply_id"`
}

type GetAllCommentsParams struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}
