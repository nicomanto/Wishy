package controllers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"wishy/models"
	"wishy/templates"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const managerSessionDuration = 12 * time.Hour

var errInvalidCredentials = errors.New("invalid username or password")

type managerPageData struct {
	Authenticated bool
	User          models.User
	Categories    []models.Category
	Wishes        []models.Wish
	Error         string
	Path          string
}

type managementHTTPError struct {
	StatusCode int
	Err        error
}

func (e *managementHTTPError) Error() string {
	return e.Err.Error()
}

func (e *managementHTTPError) Unwrap() error {
	return e.Err
}

func managerError(status int, message string, err error) error {
	if err != nil {
		return &managementHTTPError{StatusCode: status, Err: fmt.Errorf("%s: %w", message, err)}
	}
	return &managementHTTPError{StatusCode: status, Err: errors.New(message)}
}

// authenticateManager verifies bcrypt credentials, records the login, and creates a session token.
func authenticateManager(ctx context.Context, username, password string, db *mongo.Database) (string, error) {
	var user models.User
	err := db.Collection(user.DBCollectionName()).FindOne(ctx, bson.M{"uname": strings.TrimSpace(username)}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", errInvalidCredentials
		}
		return "", managerError(500, "find manager account", err)
	}
	if password == "" {
		return "", errInvalidCredentials
	}
	if !verifyManagerPassword(user.Password, password) {
		return "", errInvalidCredentials
	}

	now := time.Now()
	filter := bson.M{"_id": user.ID}
	updates := bson.M{"last_login": now}
	result, err := db.Collection(user.DBCollectionName()).UpdateOne(ctx, filter, bson.M{"$set": updates})
	if err != nil {
		return "", managerError(500, "update last login", err)
	}
	if result.MatchedCount == 0 {
		return "", managerError(401, "manager account no longer exists", nil)
	}
	return signManagerSession(user, time.Now().Add(managerSessionDuration)), nil
}

func verifyManagerPassword(storedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)) == nil
}

// managerFromSession validates the token signature and reloads its user from MongoDB.
func managerFromSession(ctx context.Context, token string, db *mongo.Database) (models.User, error) {
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

// loadManagerPage fetches the signed-in user's categories and wishes for the template.
func loadManagerPage(ctx context.Context, user models.User, db *mongo.Database) (managerPageData, error) {
	page := managerPageData{
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
	wishCursor, err := db.Collection(models.Wish{}.DBCollectionName()).Find(ctx, bson.M{"uid": user.ID}, options.Find().SetSort(bson.D{{Key: "ts", Value: -1}, {Key: "name", Value: 1}}))
	if err != nil {
		return page, managerError(500, "load wishes", err)
	}
	if err := wishCursor.All(ctx, &page.Wishes); err != nil {
		return page, managerError(500, "read wishes", err)
	}
	return page, nil
}

// manageData dispatches a submitted form action to its category or wish operation.
func manageData(ctx context.Context, action string, form url.Values, user models.User, db *mongo.Database) error {
	switch action {
	case "category_create":
		return createCategory(ctx, form, user, db)
	case "category_update":
		return updateCategory(ctx, form, user, db)
	case "category_delete":
		return deleteCategory(ctx, form, user, db)
	case "wish_create":
		return createWish(ctx, form, user, db)
	case "wish_update":
		return updateWish(ctx, form, user, db)
	case "wish_delete":
		return deleteWish(ctx, form, user, db)
	default:
		return managerError(400, "unknown management action", nil)
	}
}

// createCategory validates and inserts a category owned by the signed-in user.
func createCategory(ctx context.Context, form url.Values, user models.User, db *mongo.Database) error {
	name, err := categoryName(form.Get("name"))
	if err != nil {
		return err
	}
	category := models.Category{BaseCategory: models.BaseCategory{ID: primitive.NewObjectID(), Name: name}, UserId: user.ID}
	if _, err := db.Collection(category.DBCollectionName()).InsertOne(ctx, category); err != nil {
		return managerError(500, "create category", err)
	}
	return nil
}

// updateCategory renames a category and synchronizes the embedded category in its wishes.
func updateCategory(ctx context.Context, form url.Values, user models.User, db *mongo.Database) error {
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
	now := time.Now()
	wishResult, err := db.Collection(models.Wish{}.DBCollectionName()).UpdateMany(ctx, bson.M{"uid": user.ID, "cat._id": id}, bson.M{"$set": bson.M{"cat.name": name, "ts": now}})
	if err != nil {
		return managerError(500, "update category in wishes", err)
	}
	if wishResult.ModifiedCount > 0 {
		return updateUserWishTimestamp(ctx, user, db)
	}
	return nil
}

// deleteCategory refuses to remove categories that still contain wishes.
func deleteCategory(ctx context.Context, form url.Values, user models.User, db *mongo.Database) error {
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
}

// createWish validates the form, inserts the wish, and updates the user's wishlist timestamp.
func createWish(ctx context.Context, form url.Values, user models.User, db *mongo.Database) error {
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
}

// updateWish replaces the editable fields and refreshes the wish and user timestamps.
func updateWish(ctx context.Context, form url.Values, user models.User, db *mongo.Database) error {
	id, err := objectIDField(form.Get("wish_id"), "wish")
	if err != nil {
		return err
	}
	wish, err := wishFromForm(ctx, form, user, db)
	if err != nil {
		return err
	}
	wish.Active = form.Get("active") == "on"
	wish.Ts = time.Now()
	result, err := db.Collection(wish.DBCollectionName()).UpdateOne(ctx, bson.M{"_id": id, "uid": user.ID}, bson.M{"$set": bson.M{
		"name": wish.Name, "link": wish.Link, "preference": wish.Preference, "cat": wish.Category, "active": wish.Active, "ts": wish.Ts,
	}})
	if err != nil {
		return managerError(500, "update wish", err)
	}
	if result.MatchedCount == 0 {
		return managerError(404, "wish not found", nil)
	}
	return updateUserWishTimestamp(ctx, user, db)
}

// deleteWish removes a wish owned by the user and updates the wishlist timestamp.
func deleteWish(ctx context.Context, form url.Values, user models.User, db *mongo.Database) error {
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
	result, err := db.Collection(user.DBCollectionName()).UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"last_wish_update": time.Now()}})
	if err != nil {
		return managerError(500, "update wishlist timestamp", err)
	}
	if result.MatchedCount == 0 {
		return managerError(404, "user not found while updating wishlist timestamp", nil)
	}
	return nil
}

// HandleManagement routes the management page's HTTP methods to their handlers.
func HandleManagement(ctx context.Context, request events.APIGatewayProxyRequest, db *mongo.Database) (events.APIGatewayProxyResponse, error) {
	templates.InitHtmlTpls()

	switch request.HTTPMethod {
	case http.MethodGet:
		// GET renders either the sign-in form or the authenticated management page.
		return getManagerPage(ctx, request, db)
	case http.MethodPost:
		// POST dispatches sign-in, sign-out, and authenticated data-management actions.
		return postManagerForm(ctx, request, db)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Headers:    map[string]string{"Allow": "GET, POST"},
		}, nil
	}
}

func getManagerPage(ctx context.Context, request events.APIGatewayProxyRequest, db *mongo.Database) (events.APIGatewayProxyResponse, error) {
	user, err := managerUserFromRequest(ctx, request, db)
	if err != nil {
		if managerStatus(err) != http.StatusUnauthorized {
			logrus.Errorln(err)
			return renderManagerPageForRequest(request, managerPageData{Error: "An error occurred. Please try again later."}, http.StatusInternalServerError, nil)
		}
		return renderManagerPageForRequest(request, managerPageData{}, http.StatusOK, map[string]string{"Set-Cookie": expiredManagerCookie()})
	}

	page, err := loadManagerPage(ctx, user, db)
	if err != nil {
		logrus.Errorln(err)
		return renderManagerPageForRequest(request, managerPageData{Error: "Unable to load data."}, http.StatusInternalServerError, nil)
	}
	return renderManagerPageForRequest(request, page, http.StatusOK, nil)
}

func postManagerForm(ctx context.Context, request events.APIGatewayProxyRequest, db *mongo.Database) (events.APIGatewayProxyResponse, error) {
	form, err := managementForm(request)
	if err != nil {
		return renderManagerPageForRequest(request, managerPageData{Error: "Invalid request."}, http.StatusBadRequest, nil)
	}

	switch form.Get("action") {
	case "login":
		// Login validates the credentials and issues a signed session cookie.
		return loginManager(ctx, request, form, db)
	case "logout":
		// Logout clears the session cookie and returns to the sign-in page.
		return redirectManagerLogout(request)
	default:
		return runManagerAction(ctx, request, form, db)
	}
}

// loginManager turns successful credentials into an HTTP-only session cookie.
func loginManager(ctx context.Context, request events.APIGatewayProxyRequest, form url.Values, db *mongo.Database) (events.APIGatewayProxyResponse, error) {
	token, err := authenticateManager(ctx, form.Get("username"), form.Get("password"), db)
	if err != nil {
		if errors.Is(err, errInvalidCredentials) {
			return renderManagerPageForRequest(request, managerPageData{Error: "Incorrect username or password."}, http.StatusUnauthorized, nil)
		}
		logrus.Errorln(err)
		return renderManagerPageForRequest(request, managerPageData{Error: "An error occurred. Please try again later."}, http.StatusInternalServerError, nil)
	}
	return redirectManager(token, request)
}

// runManagerAction authenticates a form submission before changing user-owned data.
func runManagerAction(ctx context.Context, request events.APIGatewayProxyRequest, form url.Values, db *mongo.Database) (events.APIGatewayProxyResponse, error) {
	user, err := managerUserFromRequest(ctx, request, db)
	if err != nil {
		if managerStatus(err) != http.StatusUnauthorized {
			logrus.Errorln(err)
			return renderManagerPageForRequest(request, managerPageData{Error: "An error occurred. Please try again later."}, http.StatusInternalServerError, nil)
		}
		return renderManagerPageForRequest(request, managerPageData{Error: "Please sign in again."}, http.StatusUnauthorized, map[string]string{"Set-Cookie": expiredManagerCookie()})
	}

	if err := manageData(ctx, form.Get("action"), form, user, db); err != nil {
		status := managerStatus(err)
		if status >= http.StatusInternalServerError {
			logrus.Errorln(err)
		}
		page, loadErr := loadManagerPage(ctx, user, db)
		if loadErr != nil {
			logrus.Errorln(loadErr)
			return renderManagerPageForRequest(request, managerPageData{Error: "Unable to load data."}, http.StatusInternalServerError, nil)
		}
		if status >= http.StatusInternalServerError {
			page.Error = "An error occurred. Please try again later."
		} else {
			page.Error = err.Error()
		}
		return renderManagerPageForRequest(request, page, status, nil)
	}
	return redirectManager("", request)
}

// managementForm decodes the API Gateway body into standard URL-encoded form values.
func managementForm(request events.APIGatewayProxyRequest) (url.Values, error) {
	body := request.Body
	if request.IsBase64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			return nil, err
		}
		body = string(decoded)
	}
	return url.ParseQuery(body)
}

// managerUserFromRequest extracts the session cookie and validates its associated user.
func managerUserFromRequest(ctx context.Context, request events.APIGatewayProxyRequest, db *mongo.Database) (models.User, error) {
	var cookieHeader string
	for key, value := range request.Headers {
		if strings.EqualFold(key, "Cookie") {
			cookieHeader = value
			break
		}
	}
	if cookieHeader == "" {
		for key, values := range request.MultiValueHeaders {
			if strings.EqualFold(key, "Cookie") && len(values) > 0 {
				cookieHeader = values[0]
				break
			}
		}
	}
	for _, part := range strings.Split(cookieHeader, ";") {
		name, value, found := strings.Cut(strings.TrimSpace(part), "=")
		if found && name == "wishy_manager" {
			return managerFromSession(ctx, value, db)
		}
	}
	return models.User{}, fmt.Errorf("%d", http.StatusUnauthorized)
}

func managerStatus(err error) int {
	var managerErr *managementHTTPError
	if errors.As(err, &managerErr) {
		return managerErr.StatusCode
	}
	status, parseErr := strconv.Atoi(err.Error())
	if parseErr == nil && status >= 400 && status <= 599 {
		return status
	}
	return http.StatusInternalServerError
}

// renderManagerPage executes the management template and applies shared response headers.
func renderManagerPage(page managerPageData, status int, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	var body strings.Builder
	if err := templates.HtmlTpls[templates.ManagementHtmlTemplateType].Execute(&body, page); err != nil {
		logrus.Errorln(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Headers: map[string]string{"Content-Type": "text/plain"}}, fmt.Errorf("render management page: %w", err)
	}
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "text/html; charset=utf-8"
	headers["Cache-Control"] = "no-store"
	return events.APIGatewayProxyResponse{StatusCode: status, Headers: headers, Body: body.String()}, nil
}

func renderManagerPageForRequest(request events.APIGatewayProxyRequest, page managerPageData, status int, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	page.Path = managementPath(request)
	return renderManagerPage(page, status, headers)
}

func redirectManager(token string, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Location":      managementPath(request),
		"Cache-Control": "no-store",
	}
	if token != "" {
		headers["Set-Cookie"] = managerCookie(token)
	}
	return events.APIGatewayProxyResponse{StatusCode: http.StatusSeeOther, Headers: headers}, nil
}

func redirectManagerLogout(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusSeeOther,
		Headers: map[string]string{
			"Location":      managementPath(request),
			"Set-Cookie":    expiredManagerCookie(),
			"Cache-Control": "no-store",
		},
	}, nil
}

func managerCookie(token string) string {
	return "wishy_manager=" + token + "; Path=/; Max-Age=" + strconv.Itoa(int(managerSessionDuration.Seconds())) + "; HttpOnly; Secure; SameSite=Strict"
}

func expiredManagerCookie() string {
	return "wishy_manager=; Path=/; Max-Age=0; HttpOnly; Secure; SameSite=Strict"
}

// managementPath builds the public path including an API Gateway stage prefix.
func managementPath(request events.APIGatewayProxyRequest) string {
	stage := request.RequestContext.Stage
	if stage == "" || stage == "$default" {
		if strings.HasSuffix(request.Path, "/manage") {
			return request.Path
		}
		return "/manage"
	}
	prefix := "/" + stage
	if strings.HasPrefix(request.Path, prefix+"/") {
		return request.Path
	}
	return prefix + "/manage"
}

// NormalizedRequestPath removes an API Gateway stage prefix before endpoint routing.
func NormalizedRequestPath(request events.APIGatewayProxyRequest) string {
	path := request.Path
	stage := request.RequestContext.Stage
	if stage != "" && stage != "$default" && strings.HasPrefix(path, "/"+stage+"/") {
		path = strings.TrimPrefix(path, "/"+stage)
	}
	return path
}
