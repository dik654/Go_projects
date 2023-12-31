package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/dik654/Go_projects/SNS_SERVER/controllers/dto"
	"github.com/dik654/Go_projects/SNS_SERVER/models"
	"github.com/dik654/Go_projects/SNS_SERVER/utils"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostServiceImpl struct {
	commentservice CommentService
	redisclient    *redis.Client
	postcollection *mongo.Collection
	ctx            context.Context
}

func NewPostService(commentservice CommentService, redisclient *redis.Client, postcollection *mongo.Collection, ctx context.Context) PostService {
	return &PostServiceImpl{
		commentservice: commentservice,
		redisclient:    redisclient,
		postcollection: postcollection,
		ctx:            ctx,
	}
}

func (p *PostServiceImpl) CreatePost(request *dto.CreatePostRequest, sessionInfo *dto.SessionInfo) error {
	ctx := context.Background()
	sessionDataJSON, err := p.getSessionFromRedis(ctx, sessionInfo)
	if err != nil {
		return err
	}
	var sessionData map[string]string
	if err := json.Unmarshal([]byte(sessionDataJSON), &sessionData); err != nil {
		return err
	}
	userID := sessionData["user_id"]

	var post models.Post
	post.ID = uuid.NewString()
	post.AuthorID = userID
	post.Title = request.Title
	post.Content = request.Content
	post.Judge = []models.Judge{}
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	_, err = p.postcollection.InsertOne(p.ctx, post)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostServiceImpl) GetAllPosts(paginationRequest *dto.PaginationRequest) ([]models.Post, error) {
	skip := (paginationRequest.PageNum - 1) * paginationRequest.PageSize
	limit := paginationRequest.PageSize

	cursor, err := p.postcollection.Find(
		context.Background(),
		bson.M{},
		options.Find().SetSkip(skip).SetLimit(limit),
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	var posts []models.Post
	if err = cursor.All(context.Background(), &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (p *PostServiceImpl) GetPosts(getPostRequest *dto.GetPostRequest) ([]models.Post, error) {
	filter := bson.M{}

	switch getPostRequest.Type {
	case 1:
		filter["post_title"] = bson.M{"$regex": getPostRequest.Body, "$options": "i"}
	case 2:
		filter["post_content"] = bson.M{"$regex": getPostRequest.Body, "$options": "i"}
	case 3:
		filter["post_author_id"] = bson.M{"$regex": getPostRequest.Body, "$options": "i"}
	default:
		return nil, errors.New("GET_POSTS_ERROR: invalid search type")
	}

	cursor, err := p.postcollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var posts []models.Post
	if err = cursor.All(context.Background(), &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (p *PostServiceImpl) UpdatePost(postID string, post *dto.CreatePostRequest) error {
	filter := bson.M{"post_id": postID}

	// 업데이트할 내용 정의
	update := bson.M{
		"$set": bson.M{
			"post_title":      post.Title,
			"post_content":    post.Content,
			"post_updated_at": time.Now(),
		},
	}

	// MongoDB에 업데이트 요청
	result, err := p.postcollection.UpdateOne(p.ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount != 1 {
		return errors.New("UPDATE_POST_ERROR: no matched document found for update")
	}

	return nil
}

func (p *PostServiceImpl) DeletePost(postId string) error {
	filter := bson.D{bson.E{Key: "post_id", Value: postId}}
	result, _ := p.postcollection.DeleteOne(p.ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("DELETE_USER_ERROR: no matched document found for delete")
	}
	if err := p.commentservice.DeleteComments(postId); err != nil {
		return err
	}
	return nil
}

func (p *PostServiceImpl) LikePost(postLikeRequest *dto.PostLikeRequest) error {
	filter := bson.M{"post_id": postLikeRequest.PostID}

	var update bson.M
	if postLikeRequest.Like {
		// 좋아요 수 증가
		update = bson.M{
			"$inc": bson.M{
				"post_like": 1,
			},
		}
	} else {
		// 싫어요 수 증가
		update = bson.M{
			"$inc": bson.M{
				"post_unlike": 1,
			},
		}
	}

	// MongoDB에 업데이트 요청
	result, err := p.postcollection.UpdateOne(p.ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount != 1 {
		return errors.New("UPDATE_POST_ERROR: no matched document found for update")
	}

	return nil
}

func (p *PostServiceImpl) JudgePost(postJudgeRequest models.Judge) error {
	ctx := context.Background()

	filter := bson.M{"post_id": postJudgeRequest.PostID}

	update := bson.M{
		"$push": bson.M{
			"post_judge": postJudgeRequest,
		},
	}

	result, err := p.postcollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount != 1 {
		return errors.New("UPDATE_POST_ERROR: no matched document found for update")
	}

	return nil
}

// 중복되는 메서드는 common service로 통합하여 리팩토링 필요
func (p *PostServiceImpl) getSessionFromRedis(ctx context.Context, sessionInfo *dto.SessionInfo) (string, error) {
	var sessionDataJSON string
	var err error
	switch sessionInfo.UserType {
	case "regular":
		sessionDataJSON, err = p.redisclient.Get(ctx, "regular_user_session:"+sessionInfo.SessionUUID).Result()
	case "google":
		sessionDataJSON, err = p.redisclient.Get(ctx, "google_user_session:"+sessionInfo.SessionUUID).Result()
	}
	if err != nil {
		return "", err
	}

	sessionDataJSON, _, err = utils.SeparateSessionDataAndSignature([]byte(sessionDataJSON))
	if err != nil {
		return "", err
	}

	return sessionDataJSON, nil
}

// 중복되는 메서드는 common service로 통합하여 리팩토링 필요
func (p *PostServiceImpl) CanEditPost(ctx context.Context, sessionInfo *dto.SessionInfo, postID string) (bool, error) {
	sessionDataJSON, err := p.getSessionFromRedis(ctx, sessionInfo)
	if err != nil {
		return false, err
	}
	var sessionData map[string]string
	if err := json.Unmarshal([]byte(sessionDataJSON), &sessionData); err != nil {
		return false, err
	}
	userID := sessionData["user_id"]
	var post models.Post
	if err := p.postcollection.FindOne(ctx, bson.M{"post_id": postID}).Decode(&post); err != nil {
		return false, err
	}
	return post.AuthorID == userID, nil
}
