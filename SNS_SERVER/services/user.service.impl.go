package services

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/dik654/Go_projects/SNS_SERVER/controllers/dto"
	"github.com/dik654/Go_projects/SNS_SERVER/models"
	"github.com/dik654/Go_projects/SNS_SERVER/utils"
	"github.com/gin-contrib/sessions"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	redisclient          *redis.Client
	usercollection       *mongo.Collection
	googleusercollection *mongo.Collection
	ctx                  context.Context
}

func NewUserService(redisclient *redis.Client, usercollection *mongo.Collection, googleusercollection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImpl{
		redisclient:          redisclient,
		usercollection:       usercollection,
		googleusercollection: googleusercollection,
		ctx:                  ctx,
	}
}

func (u *UserServiceImpl) CreateUser(user *models.User) error {
	user.Password = utils.HashingPassword(user.Password)
	_, err := u.usercollection.InsertOne(u.ctx, user)
	return err
}

func (u *UserServiceImpl) GetUser(name *string) (*models.User, error) {
	var user *models.User
	query := bson.D{bson.E{Key: "user_name", Value: name}}
	err := u.usercollection.FindOne(u.ctx, query).Decode(&user)
	return user, err
}

func (u *UserServiceImpl) GetAll() ([]*models.User, error) {
	var users []*models.User
	cursor, err := u.usercollection.Find(u.ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	for cursor.Next(u.ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	cursor.Close(u.ctx)

	if len(users) == 0 {
		return nil, errors.New("GETALL_USER_ERROR: documents not found")
	}
	return users, nil
}

func (u *UserServiceImpl) UpdateUser(user *models.User) error {
	filter := bson.D{
		bson.E{
			Key: "user_name", Value: user.Name,
		},
	}
	update := bson.D{
		bson.E{
			Key: "$set", Value: bson.D{
				bson.E{
					Key: "user_name", Value: user.Name,
				},
				bson.E{
					Key: "user_age", Value: user.Age,
				},
				bson.E{
					Key: "user_address", Value: user.Address,
				},
			},
		},
	}
	result, _ := u.usercollection.UpdateOne(u.ctx, filter, update)
	if result.MatchedCount != 1 {
		return errors.New("UPDATE_USER_ERROR: no matched document found for update")
	}
	return nil
}

func (u *UserServiceImpl) DeleteUser(name *string) error {
	filter := bson.D{bson.E{Key: "user_name", Value: name}}
	result, _ := u.usercollection.DeleteOne(u.ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("DELETE_USER_ERROR: no matched document found for delete")
	}
	return nil
}

func (u *UserServiceImpl) SignIn(session sessions.Session, signInReq dto.SignInRequest) error {
	var user *dto.SignInRequest
	query := bson.D{
		{Key: "user_id", Value: signInReq.ID},
	}
	if err := u.usercollection.FindOne(u.ctx, query).Decode(&user); err != nil {
		return err
	}
	requestHashedPassword := utils.HashingPassword(signInReq.Password)
	if user.Password != requestHashedPassword {
		return errors.New("LOGIN_ERROR: invalid password")
	}

	ctx := context.Background()
	uuid := uuid.NewString()
	signature := utils.CreateSignature(uuid, os.Getenv("SECRET_KEY"))
	combinedValue := utils.CombineSessionDataAndSignature(uuid, signature)
	session.Options(sessions.Options{
		HttpOnly: true,
		Secure:   true,
		MaxAge:   int(30 * time.Minute),
	})
	session.Set("regular_user_session", combinedValue)
	session.Save()
	if cmd := u.redisclient.Set(ctx, "regular_user_session:"+uuid, combinedValue, 30*time.Minute); cmd.Err() != nil {
		return errors.New("LOGIN_ERROR: " + cmd.Err().Error())
	}

	return nil
}

func (u *UserServiceImpl) SignOut(session sessions.Session) error {
	ctx := context.Background()
	sessionValue := session.Get("regular_user_session")
	sessionString, ok := sessionValue.([]byte)
	if !ok {
		return errors.New("LOGOUT_ERROR: session value is not []byte")
	}
	sessionID, _, _ := utils.SeparateSessionDataAndSignature(sessionString)
	if sessionID == "" {
		return errors.New("LOGOUT_ERROR: Invalid session token")
	}
	u.redisclient.Del(ctx, "regular_user_session:"+sessionID)
	session.Delete("regular_user_session")
	session.Options(sessions.Options{MaxAge: -1})
	session.Save()

	return nil
}

func (u *UserServiceImpl) GoogleSignIn(session sessions.Session, userInfo *models.GoogleUser) error {
	var user *models.GoogleUser
	query := bson.D{
		bson.E{
			Key:   "google_user_id",
			Value: userInfo.ID,
		},
	}
	if err := u.googleusercollection.FindOne(u.ctx, query).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			if _, err := u.googleusercollection.InsertOne(u.ctx, userInfo); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	session.Options(sessions.Options{
		HttpOnly: true,
		Secure:   true,
		MaxAge:   int(30 * time.Minute),
	})
	ctx := context.Background()
	uuid := uuid.NewString()
	signature := utils.CreateSignature(uuid, os.Getenv("SECRET_KEY"))
	combinedValue := utils.CombineSessionDataAndSignature(uuid, signature)
	session.Set("google_user_session", combinedValue)
	if cmd := u.redisclient.Set(ctx, "google_user_session:"+uuid, combinedValue, 30*time.Minute); cmd.Err() != nil {
		return errors.New("LOGIN_ERROR: " + cmd.Err().Error())
	}

	session.Save()
	return nil
}

func (u *UserServiceImpl) GoogleSignOut(session sessions.Session) error {
	ctx := context.Background()
	sessionValue := session.Get("google_user_session")
	sessionString, ok := sessionValue.([]byte)
	if !ok {
		return errors.New("GOOGLE_LOGOUT_ERROR: session value is not []byte")
	}
	sessionID, _, _ := utils.SeparateSessionDataAndSignature(sessionString)
	if sessionID == "" {
		return errors.New("GOOGLE_LOGOUT_ERROR: Invalid session token")
	}
	u.redisclient.Del(ctx, "google_user_session:"+sessionID)
	session.Delete("google_user_session")
	session.Options(sessions.Options{MaxAge: -1})
	session.Save()
	return nil
}