package postfavourite

import (
	"context"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
    "github.com/kotopesp/sos-kotopes/pkg/logger"
    "errors"
    "time"
    "gorm.io/gorm"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.PostFavouriteStore {
	return &store{pg}
}

// GetFavouritePosts retrieves the favourite posts of a user from the database based on the GetAllPostsParams
func (s *store) GetFavouritePosts(ctx context.Context, userID int, params core.GetAllPostsParams) ([]core.Post, int, error) {
    var posts []core.Post

    query := s.DB.WithContext(ctx).Model(&core.Post{}).
        Joins("JOIN favourite_posts ON posts.id = favourite_posts.post_id").
        Where("favourite_posts.user_id = ?", userID)

    if params.Limit != nil {
        query = query.Limit(*params.Limit)
    }

    if params.Offset != nil {
        query = query.Offset(*params.Offset)
    }

    if params.Status != nil {
        query = query.Where("posts.status = ?", *params.Status)
    }

    if params.AnimalType != nil {
        query = query.Where("animals.animal_type = ?", *params.AnimalType)
    }

    if params.Gender != nil {
        query = query.Where("animals.gender = ?", *params.Gender)
    }

    if params.Color != nil {
        query = query.Where("animals.color = ?", *params.Color)
    }

    if params.SearchWord != nil && *params.SearchWord != "" {
        searchWord := "%" + *params.SearchWord + "%"
        query = query.Where("posts.title ILIKE ? OR posts.content ILIKE ?", searchWord, searchWord)
    }

    var total int64
    if err := query.Count(&total).Error; err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, 0, err
    }

    if err := query.Select("posts.*").Find(&posts).Error; err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, 0, err
    }
    
    return posts, int(total), nil
}

// GetPostFavouriteByPostAndUserID retrieves a post favourite from the database based on the post ID and user ID
func (s *store) GetPostFavouriteByPostAndUserID(ctx context.Context, postID, userID int) (core.PostFavourite, error) {
    var postFavourite core.PostFavourite

    err := s.DB.WithContext(ctx).
        Where("post_id = ? AND user_id = ?", postID, userID).
        First(&postFavourite).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return core.PostFavourite{}, core.ErrPostNotFound
        }
        return core.PostFavourite{}, err
    }

    return postFavourite, nil
}

// AddToFavourites adds a post to the user's favourites in the database 
func (s *store) AddToFavourites(ctx context.Context, postFavourite core.PostFavourite) (core.Post, error) {
    // Check if the post exists
    var post core.Post
    if err := s.DB.WithContext(ctx).Where("id = ?", postFavourite.PostID).First(&post).Error; err != nil {
        if errors.Is(err, core.ErrRecordNotFound) {
            logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
            return core.Post{}, core.ErrPostNotFound
        }
        logger.Log().Error(ctx, err.Error())
        return core.Post{}, err
    }

    // Set is_favourite flag to true
    if err := s.DB.WithContext(ctx).Model(&post).Where("id = ?", postFavourite.PostID).
        Update("is_favourite", true).Error; err != nil {
        logger.Log().Error(ctx, err.Error())
        return core.Post{}, err
    }

    // Check if the post is already in the user's favourites
    var existing core.PostFavourite
    if err := s.DB.WithContext(ctx).
        Where("post_id = ? AND user_id = ?", postFavourite.PostID, postFavourite.UserID).
        First(&existing).Error; err == nil {
        return core.Post{}, core.ErrPostAlreadyInFavourites
    }

    postFavourite.CreatedAt = time.Now()

    if err := s.DB.WithContext(ctx).Create(&postFavourite).Error; err != nil {
        logger.Log().Error(ctx, err.Error())
        return core.Post{}, err
    }

    return post, nil
}

// DeleteFromFavourites removes a post from the user's favourites in the database
func (s *store) DeleteFromFavourites(ctx context.Context, postID, userID int) error {
	if err := s.DB.WithContext(ctx).Where("post_id = ? AND user_id = ?", postID, userID).Delete(&core.PostFavourite{}).Error; err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return core.ErrPostNotFound
		}
		logger.Log().Error(ctx, err.Error())
		return err
	}

	return nil
}
