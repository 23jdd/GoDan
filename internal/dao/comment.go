package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"godan/internal/model"
	mgodb "godan/internal/pkg/mongodb"
)

const commentCol = "comments"

func CreateComment(c *model.Comment) (string, error) {
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	c.Status = model.CommentNormal

	result, err := mgodb.Collection(commentCol).InsertOne(context.Background(), c)
	if err != nil {
		return "", err
	}
	id := result.InsertedID.(bson.ObjectID)
	c.ID = id
	return id.Hex(), nil
}

func GetCommentByID(id string) (*model.Comment, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var c model.Comment
	err = mgodb.Collection(commentCol).FindOne(context.Background(), bson.M{"_id": oid}).Decode(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func GetRootComments(videoID uint64, sort string, offset, limit int64) ([]model.Comment, int64, error) {
	filter := bson.M{"video_id": videoID, "parent_id": ""}

	total, err := mgodb.Collection(commentCol).CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	sortOpt := bson.D{{Key: "created_at", Value: -1}}
	if sort == "hot" {
		sortOpt = bson.D{{Key: "like_count", Value: -1}}
	}

	opts := options.Find().SetSort(sortOpt).SetSkip(offset).SetLimit(limit)
	cursor, err := mgodb.Collection(commentCol).Find(context.Background(), filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	var comments []model.Comment
	if err := cursor.All(context.Background(), &comments); err != nil {
		return nil, 0, err
	}
	if comments == nil {
		comments = []model.Comment{}
	}
	return comments, total, nil
}

func GetReplies(rootID string, offset, limit int64) ([]model.Comment, int64, error) {
	filter := bson.M{"root_id": rootID, "parent_id": bson.M{"$ne": ""}}

	total, err := mgodb.Collection(commentCol).CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}).SetSkip(offset).SetLimit(limit)
	cursor, err := mgodb.Collection(commentCol).Find(context.Background(), filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	var comments []model.Comment
	if err := cursor.All(context.Background(), &comments); err != nil {
		return nil, 0, err
	}
	if comments == nil {
		comments = []model.Comment{}
	}
	return comments, total, nil
}

func DeleteComment(id string) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = mgodb.Collection(commentCol).UpdateOne(
		context.Background(),
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{"status": model.CommentDeleted, "updated_at": time.Now()}},
	)
	return err
}

func IncrReplyCount(id string) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = mgodb.Collection(commentCol).UpdateOne(
		context.Background(),
		bson.M{"_id": oid},
		bson.M{"$inc": bson.M{"reply_count": 1}},
	)
	return err
}

func IncrCommentLikeCount(id string, delta int) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = mgodb.Collection(commentCol).UpdateOne(
		context.Background(),
		bson.M{"_id": oid},
		bson.M{"$inc": bson.M{"like_count": delta}},
	)
	return err
}

func IsCommentLikedByUser(commentID string, userID uint64) (bool, error) {
	var doc struct {
		ID bson.ObjectID `bson:"_id"`
	}
	err := mgodb.Collection("comment_likes").FindOne(
		context.Background(),
		bson.M{"comment_id": commentID, "user_id": userID},
	).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	return err == nil, err
}

func AddCommentLike(commentID string, userID uint64) error {
	oid, err := bson.ObjectIDFromHex(commentID)
	if err != nil {
		return err
	}
	doc := bson.M{
		"comment_id": oid.Hex(),
		"user_id":    userID,
		"created_at": time.Now(),
	}
	_, err = mgodb.Collection("comment_likes").InsertOne(context.Background(), doc)
	return err
}

func RemoveCommentLike(commentID string, userID uint64) error {
	_, err := mgodb.Collection("comment_likes").DeleteOne(
		context.Background(),
		bson.M{"comment_id": commentID, "user_id": userID},
	)
	return err
}
