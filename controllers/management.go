package controllers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
	"wishy/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const managerSessionDuration = 12 * time.Hour

var ErrInvalidCredentials = errors.New("invalid username or password")

type ManagerPageData struct {
	Authenticated bool
	User          models.User
	Categories    []models.Category
	Wishes        []models.Wish
	Error         string
}

type ManagerError struct {
	StatusCode int
	Err        error
}

func (e *ManagerError) Error() string {
	return e.Err.Error()
}

func (e *ManagerError) Unwrap() error {
	return e.Err
}

func managerError(status int, message string, err error) error {
	if err != nil {
		return &ManagerError{StatusCode: status, Err: fmt.Errorf("%s: %w", message, err)}
	}
	return &ManagerError{StatusCode: status, Err: errors.New(message)}
}

func LoginManager(ctx context.Context, username, password string, db *mongo.Database) (string, error) {
	var user models.User
	err := db.Collection(user.DBCollectionName()).FindOne(ctx, bson.M{"username": strings.TrimSpace(username)}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", ErrInvalidCredentials
		}
		return "", managerError(500, "find manager account", err)
	}
	if password == "" || !hmac.Equal([]byte(user.Password), []byte(password)) {
		return "", ErrInvalidCredentials
	}
	return signManagerSession(user, time.Now().Add(managerSessionDuration)), nil
}

func ManagerFromSession(ctx context.Context, token string, db *mongo.Database) (models.User, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return models.User{}, managerError(401, "invalid manager session", nil)
	}
	userID, err := primitive.ObjectIDFromHex(parts[0])
	if err != nil {
		return models.User{}, managerError(401, "invalid manager session", err)
	}
	expiresAt, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return models.User{}, managerError(401, "expired or invalid manager session", err)
	}
	if time.Now().Unix() >= expiresAt {
		return models.User{}, managerError(401, "expired manager session", nil)
	}

	var user models.User
	err = db.Collection(user.DBCollectionName()).FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, managerError(401, "invalid manager session", err)
		}
		return models.User{}, managerError(500, "find manager account", err)
	}
	if user.Password == "" {
		return models.User{}, managerError(401, "invalid manager session", nil)
	}
	expected := signManagerSession(user, time.Unix(expiresAt, 0))
	if !hmac.Equal([]byte(expected), []byte(token)) {
		return models.User{}, managerError(401, "invalid manager session", nil)
	}
	return user, nil
}

func signManagerSession(user models.User, expiresAt time.Time) string {
	id := user.ID.Hex()
	expiry := strconv.FormatInt(expiresAt.Unix(), 10)
	payload := id + "." + expiry
	mac := hmac.New(sha256.New, []byte(user.Password))
	_, _ = mac.Write([]byte(payload))
	return payload + "." + hex.EncodeToString(mac.Sum(nil))
}

func LoadManagerPage(ctx context.Context, user models.User, db *mongo.Database) (ManagerPageData, error) {
	page := ManagerPageData{
		Authenticated: true,
		User:          user,
		Categories:    []models.Category{},
		Wishes:        []models.Wish{},
	}
	categoryCursor, err := db.Collection(models.Category{}.DBCollectionName()).Find(ctx, bson.M{"uid": user.ID}, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))
	if err != nil {
		return page, managerError(500, "load categories", err)
	}
	if err := categoryCursor.All(ctx, &page.Categories); err != nil {
		return page, managerError(500, "read categories", err)
	}
	wishCursor, err := db.Collection(models.Wish{}.DBCollectionName()).Find(ctx, bson.M{"uid": user.ID}, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))
	if err != nil {
		return page, managerError(500, "load wishes", err)
	}
	if err := wishCursor.All(ctx, &page.Wishes); err != nil {
		return page, managerError(500, "read wishes", err)
	}
	return page, nil
}

func ManageData(ctx context.Context, action string, form url.Values, user models.User, db *mongo.Database) error {
	switch action {
	case "category_create":
		name, err := categoryName(form.Get("name"))
		if err != nil {
			return err
		}
		category := models.Category{BaseCategory: models.BaseCategory{ID: primitive.NewObjectID(), Name: name}, UserId: user.ID}
		if _, err := db.Collection(category.DBCollectionName()).InsertOne(ctx, category); err != nil {
			return managerError(500, "create category", err)
		}
		return nil
	case "category_update":
		id, err := objectIDField(form.Get("category_id"), "category")
		if err != nil {
			return err
		}
		name, err := categoryName(form.Get("name"))
		if err != nil {
			return err
		}
		result, err := db.Collection(models.Category{}.DBCollectionName()).UpdateOne(ctx, bson.M{"_id": id, "uid": user.ID}, bson.M{"$set": bson.M{"name": name}})
		if err != nil {
			return managerError(500, "update category", err)
		}
		if result.MatchedCount == 0 {
			return managerError(404, "category not found", nil)
		}
		_, err = db.Collection(models.Wish{}.DBCollectionName()).UpdateMany(ctx, bson.M{"uid": user.ID, "cat._id": id}, bson.M{"$set": bson.M{"cat.name": name}})
		if err != nil {
			return managerError(500, "update category in wishes", err)
		}
		return nil
	case "category_delete":
		id, err := objectIDField(form.Get("category_id"), "category")
		if err != nil {
			return err
		}
		filter := bson.M{"uid": user.ID, "cat._id": id}
		count, err := db.Collection(models.Wish{}.DBCollectionName()).CountDocuments(ctx, filter)
		if err != nil {
			return managerError(500, "check category wishes", err)
		}
		if count > 0 {
			return managerError(409, "move or delete the wishes assigned to this category first", nil)
		}
		result, err := db.Collection(models.Category{}.DBCollectionName()).DeleteOne(ctx, bson.M{"_id": id, "uid": user.ID})
		if err != nil {
			return managerError(500, "delete category", err)
		}
		if result.DeletedCount == 0 {
			return managerError(404, "category not found", nil)
		}
		return nil
	case "wish_create":
		wish, err := wishFromForm(ctx, form, user, db)
		if err != nil {
			return err
		}
		wish.ID = primitive.NewObjectID()
		wish.Ts = time.Now()
		wish.Active = form.Get("active") == "on"
		if _, err := db.Collection(wish.DBCollectionName()).InsertOne(ctx, wish); err != nil {
			return managerError(500, "create wish", err)
		}
		return updateUserWishTimestamp(ctx, user, db)
	case "wish_update":
		id, err := objectIDField(form.Get("wish_id"), "wish")
		if err != nil {
			return err
		}
		wish, err := wishFromForm(ctx, form, user, db)
		if err != nil {
			return err
		}
		wish.Active = form.Get("active") == "on"
		result, err := db.Collection(wish.DBCollectionName()).UpdateOne(ctx, bson.M{"_id": id, "uid": user.ID}, bson.M{"$set": bson.M{
			"name": wish.Name, "link": wish.Link, "preference": wish.Preference, "cat": wish.Category, "active": wish.Active,
		}})
		if err != nil {
			return managerError(500, "update wish", err)
		}
		if result.MatchedCount == 0 {
			return managerError(404, "wish not found", nil)
		}
		return updateUserWishTimestamp(ctx, user, db)
	case "wish_delete":
		id, err := objectIDField(form.Get("wish_id"), "wish")
		if err != nil {
			return err
		}
		result, err := db.Collection(models.Wish{}.DBCollectionName()).DeleteOne(ctx, bson.M{"_id": id, "uid": user.ID})
		if err != nil {
			return managerError(500, "delete wish", err)
		}
		if result.DeletedCount == 0 {
			return managerError(404, "wish not found", nil)
		}
		return updateUserWishTimestamp(ctx, user, db)
	default:
		return managerError(400, "unknown management action", nil)
	}
}

func categoryName(value string) (string, error) {
	name := strings.TrimSpace(value)
	if name == "" || len(name) > 100 {
		return "", managerError(400, "category name must be between 1 and 100 characters", nil)
	}
	return name, nil
}

func wishFromForm(ctx context.Context, form url.Values, user models.User, db *mongo.Database) (models.Wish, error) {
	name := strings.TrimSpace(form.Get("name"))
	if name == "" || len(name) > 200 {
		return models.Wish{}, managerError(400, "wish name must be between 1 and 200 characters", nil)
	}
	link := strings.TrimSpace(form.Get("link"))
	if len(link) > 2048 {
		return models.Wish{}, managerError(400, "link is too long", nil)
	}
	if link != "" {
		parsed, err := url.ParseRequestURI(link)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
			return models.Wish{}, managerError(400, "enter a valid HTTP or HTTPS link", err)
		}
	}
	preference, err := strconv.Atoi(form.Get("preference"))
	if err != nil || preference < int(models.Low) || preference > int(models.High) {
		return models.Wish{}, managerError(400, "invalid preference", err)
	}
	categoryID, err := objectIDField(form.Get("category_id"), "category")
	if err != nil {
		return models.Wish{}, err
	}
	var category models.Category
	err = db.Collection(category.DBCollectionName()).FindOne(ctx, bson.M{"_id": categoryID, "uid": user.ID}).Decode(&category)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Wish{}, managerError(400, "category not found", err)
		}
		return models.Wish{}, managerError(500, "find category", err)
	}
	return models.Wish{
		BaseWish: models.BaseWish{
			Name:       name,
			Link:       link,
			Preference: models.PreferenceType(preference),
		},
		Category: category.BaseCategory,
		UserId:   user.ID,
	}, nil
}

func objectIDField(value, field string) (primitive.ObjectID, error) {
	id, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		return primitive.NilObjectID, managerError(400, "invalid "+field+" id", err)
	}
	return id, nil
}

func updateUserWishTimestamp(ctx context.Context, user models.User, db *mongo.Database) error {
	_, err := db.Collection(user.DBCollectionName()).UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"last_wish_update": time.Now()}})
	if err != nil {
		return managerError(500, "update wishlist timestamp", err)
	}
	return nil
}
