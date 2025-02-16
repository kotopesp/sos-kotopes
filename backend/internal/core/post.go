package core

import (
	"context"
	"mime/multipart"
	"time"
)

type (
	Post struct {
		ID        int       `gorm:"column:id;primaryKey"` // Unique identifier for the post
		AuthorID  int       `gorm:"column:author_id"`     // ID of the author of the post
		AnimalID  int       `gorm:"column:animal_id"`     // ID of the associated animal
		Title     string    `gorm:"column:title"`         // Title of the post
		Content   string    `gorm:"column:content"`       // Content of the post
		CreatedAt time.Time `gorm:"column:created_at"`    // Timestamp when the post was created
		UpdatedAt time.Time `gorm:"column:updated_at"`    // Timestamp when the post was last updated
		Photo     []byte    `gorm:"column:photo"`         // Photo animal
		Status    string    `gorm:"column:status"`        // Status shows current status of post e.g. Deleted | Published | OnModeration
	}

	// PostDetails Post Details joins post, animal, username
	PostDetails struct {
		Post     Post
		Animal   Animal
		Username string
	}

	// Report consists of counter for report amount and reasons why post was reported.
	Report struct {
		ID         int       `gorm:"column:id;primaryKey"`
		PostID     int       `gorm:"column:post_id"`
		ReportedAt time.Time `gorm:"column:reported_at"`
		Reason     string    `gorm:"column:reason"`
		UserID     int       `gorm:"column:user_id"`
	}

	// UpdateRequestBodyPost represents the request body for updating a post.
	UpdateRequestBodyPost struct {
		ID          *int
		AuthorID    *int
		Title       *string
		Content     *string
		Photo       *[]byte
		AnimalType  *string
		Age         *int
		Color       *string
		Gender      *string
		Description *string
		Status      *string
	}

	// GetAllPostsParams are needed for processing posts in the database
	GetAllPostsParams struct {
		Limit      *int    // Limit on the number of posts to retrieve
		Offset     *int    // Offset for pagination
		Status     *string // Filter by status of the associated animal
		AnimalType *string // Filter by type of the associated animal
		Gender     *string // Filter by gender of the associated animal
		Color      *string // Filter by color of the associated animal
		Location   *string // Filter by location of the associated animal
	}

	PostStore interface {
		GetAllPosts(ctx context.Context, params GetAllPostsParams) ([]Post, int, error)
		GetUserPosts(ctx context.Context, id int) (posts []Post, count int, err error)
		GetPostByID(ctx context.Context, id int) (Post, error)
		CreatePost(ctx context.Context, post Post) (Post, error)
		UpdatePost(ctx context.Context, post Post) (Post, error)
		DeletePost(ctx context.Context, id int) error
		ReportPost(ctx context.Context, report Report) (reportCount int, err error)
		SendToModeration(ctx context.Context, postID int) (err error)
		GetPostForModeration(ctx context.Context) (Post, error)
	}

	PostService interface {
		GetAllPosts(ctx context.Context, params GetAllPostsParams) ([]PostDetails, int, error)
		GetUserPosts(ctx context.Context, id int) (posts []PostDetails, count int, err error)
		GetPostByID(ctx context.Context, id int) (PostDetails, error)
		CreatePost(ctx context.Context, postDetails PostDetails, fileHeader *multipart.FileHeader) (PostDetails, error)
		UpdatePost(ctx context.Context, postUpdateRequest UpdateRequestBodyPost) (PostDetails, error)
		DeletePost(ctx context.Context, post Post) error
		ReportPost(ctx context.Context, post Report) (err error)
		GetPostForModeration(ctx context.Context) (Post, error)
		PostFavouriteService
	}
)

func (Post) TableName() string {
	return "posts"
}

func (Report) TableName() string { return "reports" }

const Published = "published"
const Deleted = "deleted"
const OnModeration = "on_moderation"
