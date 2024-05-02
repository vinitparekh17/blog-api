package helper

import (
	"errors"
	"reflect"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jay-bhogayata/blogapi/internal/database"
	"github.com/jay-bhogayata/blogapi/internal/pb/blogservice"
)

func MarshalToGRPCTag(tag *blogservice.BlogTag) *database.Tag {
	return &database.Tag{
		ID:   tag.Id,
		Name: tag.Name,
	}
}

func MarshalToGRPCBlog(blog *database.Article) *blogservice.Blog {

	var blogId string
	blog.ArticleID.Scan(&blogId)

	var authorId string
	blog.UserID.Scan(&authorId)

	return &blogservice.Blog{
		Id:          blogId,
		Title:       blog.Title,
		Content:     blog.Content,
		AuthorId:    authorId,
		TagId:       blog.TagID.Int32,
		CreatedAt:   blog.CreatedAt.Time.Unix(),
		UpdatedAt:   blog.UpdatedAt.Time.Unix(),
		IsPublished: blog.IsPublished.Bool,
	}
}

// T type it can either be a blog or a tag type
func UnMarshalToDbType[T comparable, U comparable](obj T) (U, error) {
	// Check the type of obj and handle accordingly
	// if reflect.TypeOf(obj) == reflect.TypeOf(&blogservice.Blog{}) {

	var blogObj *blogservice.Blog

	var blogTagObj *blogservice.BlogTag

	switch reflect.TypeOf(obj).Elem().Name() {
	case "Blog":
		blogObj = reflect.ValueOf(obj).Elem().Interface().(*blogservice.Blog)

		var blogId pgtype.UUID
		blogId.Scan(blogObj.Id)

		var authorId pgtype.UUID
		authorId.Scan(blogObj.AuthorId)

		var tagId pgtype.Int4
		tagId.Scan(blogObj.TagId)

		var createdAt pgtype.Timestamp
		createdAt.Scan(blogObj.CreatedAt)

		var updatedAt pgtype.Timestamp
		updatedAt.Scan(blogObj.UpdatedAt)

		return reflect.ValueOf(&database.Article{
			ArticleID:   blogId,
			Title:       blogObj.Title,
			Content:     blogObj.Content,
			UserID:      authorId,
			TagID:       tagId,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			IsPublished: pgtype.Bool{Bool: blogObj.IsPublished},
		}).Interface().(U), nil

	case "BlogTag":
		blogTagObj = reflect.ValueOf(obj).Elem().Interface().(*blogservice.BlogTag)

		return reflect.ValueOf(&database.Tag{
			ID:   blogTagObj.Id,
			Name: blogTagObj.Name,
		}).Interface().(U), nil

	default:
		var zeroValue U
		return zeroValue, errors.New("unsupported type")
	}
}
