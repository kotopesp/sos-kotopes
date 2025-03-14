package core

import (
	"context"
	"mime/multipart"
	"time"
)

type (
	Post struct {
		ID        int        `gorm:"column:id;primaryKey"` // Unique identifier for the post
		AuthorID  int        `gorm:"column:author_id"`     // ID of the author of the post
		AnimalID  int        `gorm:"column:animal_id"`     // ID of the associated animal
		Title     string     `gorm:"column:title"`         // Title of the post
		Content   string     `gorm:"column:content"`       // Content of the post
		Photo     []byte     `gorm:"column:photo"`         // Photo animal
		Status    PostStatus `gorm:"column:status"`        // Status shows current status of post
		CreatedAt time.Time  `gorm:"column:created_at"`    // Timestamp when the post was created
		DeletedAt time.Time  `gorm:"column:deleted_at"`    // Timestamp when the post was deleted
		UpdatedAt time.Time  `gorm:"column:updated_at"`    // Timestamp when the post was last updated
	}

	// PostDetails Post Details joins post, animal, username
	PostDetails struct {
		Post     Post
		Animal   Animal
		Username string
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

	// PostForModeration structure that holds post and list of reasons why this post was reported.
	PostForModeration struct {
		Post    Post
		Reasons []string
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
		SendToModeration(ctx context.Context, postID int) (err error)
	}

	PostService interface {
		GetAllPosts(ctx context.Context, params GetAllPostsParams) ([]PostDetails, int, error)
		GetUserPosts(ctx context.Context, id int) (posts []PostDetails, count int, err error)
		GetPostByID(ctx context.Context, id int) (PostDetails, error)
		CreatePost(ctx context.Context, postDetails PostDetails, fileHeader *multipart.FileHeader) (PostDetails, error)
		UpdatePost(ctx context.Context, postUpdateRequest UpdateRequestBodyPost) (PostDetails, error)
		DeletePost(ctx context.Context, post Post) error
		PostFavouriteService
	}
)

// AmountOfPostsForModeration defines the maximum number of posts that can be fetched for moderation at once.
const AmountOfPostsForModeration = 10

// PostStatus is a custom type that represents the current state of a post.
type PostStatus string

const (
	Published    PostStatus = "published"
	Deleted      PostStatus = "deleted"
	OnModeration PostStatus = "on_moderation"
)

// Filter is a custom type that specifies the order for retrieving posts: ASC (ascending) or DESC (descending).
type Filter string

const (
	FilterDESC Filter = "DESC"
	FilterASC  Filter = "ASC"
)

func (Post) TableName() string {
	return "posts"
}
