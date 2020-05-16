package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gavinc95/go-blog/db/models"
	"github.com/stretchr/testify/require"
)

var (
	app           *App
	uuidGenerator = &stubUUIDGenerator{}
)

var (
	samplePostID  = "85b02cdf-0021-4c82-a80a-9e8788503734"
	samplePostID2 = "85b02cdf-0021-4c82-a80a-9e87885037aa"

	sampleUserID = "553e5015-ce17-4c10-abf3-e7329f063dc9"
)

// this is used to prevent random UUIDs from being created for testing
type stubUUIDGenerator struct {
	shouldGenPostID bool
	shouldGenUserID bool
}

func (g *stubUUIDGenerator) UUID() string {
	if g.shouldGenUserID {
		return sampleUserID
	} else if g.shouldGenPostID {
		return samplePostID
	}

	return ""
}

func TestMain(m *testing.M) {
	app = NewApp(":8010", uuidGenerator)
	app.ensureTablesExists()

	code := m.Run()
	clearTable()
	os.Exit(code)
}

func clearTable() {
	if _, err := app.BlogStore.GetDB().Exec("DELETE FROM users"); err != nil {
		log.Fatal(err)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected response code: %d, but got: %d\n", expected, actual)
	}
}

func TestGetUser_Empty(t *testing.T) {
	clearTable()

	reqBytes, err := json.Marshal(&GetUserRequest{ID: "8440fc74-16f3-47b1-8b27-eb2851d2afaa"})
	require.NoError(t, err)
	req, err := http.NewRequest("GET", "/users", bytes.NewBuffer(reqBytes))
	require.NoError(t, err)
	resp := executeRequest(req)

	var user models.User
	err = json.Unmarshal(resp.Body.Bytes(), &user)
	require.NoError(t, err)
	checkResponseCode(t, http.StatusOK, resp.Code)
	require.Equal(t, "", user.ID)
	require.Equal(t, "", user.Name)
	require.Equal(t, "", user.Email)
}

func TestGetNonExistentUser(t *testing.T) {
	clearTable()

	resp := getTestUser(t, sampleUserID)
	var user models.User
	err := json.Unmarshal(resp.Body.Bytes(), &user)
	require.NoError(t, err)
	checkResponseCode(t, http.StatusOK, resp.Code)
	require.Equal(t, "", user.ID)
	require.Equal(t, "", user.Name)
	require.Equal(t, "", user.Email)
}

func TestCreateUser(t *testing.T) {
	clearTable()

	uuidGenerator.shouldGenUserID = true
	resp := createTestUser(t, "tiny cat", "tiny@cat.com")
	checkResponseCode(t, http.StatusOK, resp.Code)
	var res CreateUserResponse
	err := json.Unmarshal(resp.Body.Bytes(), &res)
	require.NoError(t, err)
	require.Equal(t, sampleUserID, res.ID)

	resp = getTestUser(t, sampleUserID)
	checkResponseCode(t, http.StatusOK, resp.Code)
	var getResp GetUserResponse
	err = json.Unmarshal(resp.Body.Bytes(), &getResp)
	require.NoError(t, err)
	require.NotNil(t, getResp.User)
	require.Equal(t, sampleUserID, getResp.User.ID)
	require.Equal(t, "tiny cat", getResp.User.Name)
	require.Equal(t, "tiny@cat.com", getResp.User.Email)
}

func TestCreateExistingUser(t *testing.T) {
	clearTable()

	uuidGenerator.shouldGenUserID = true
	resp := createTestUser(t, "tiny cat", "tiny@cat.com")
	checkResponseCode(t, http.StatusOK, resp.Code)
	var res CreateUserResponse
	err := json.Unmarshal(resp.Body.Bytes(), &res)
	require.NoError(t, err)
	require.Equal(t, sampleUserID, res.ID)

	// create the same user again, and check for an error
	resp = createTestUser(t, "tiny cat", "tiny@cat.com")
	require.Equal(t, http.StatusInternalServerError, resp.Result().StatusCode)
}

func TestCreateAndUpdateUser(t *testing.T) {
	clearTable()

	// create a new user
	uuidGenerator.shouldGenUserID = true
	resp := createTestUser(t, "tiny cat", "tiny@cat.com")
	checkResponseCode(t, http.StatusOK, resp.Code)
	var res CreateUserResponse
	err := json.Unmarshal(resp.Body.Bytes(), &res)
	require.NoError(t, err)
	require.Equal(t, sampleUserID, res.ID)

	// update an existing user's email
	resp = updateTestUser(t,
		sampleUserID, "", "tiny@enterprisecatz.com")
	checkResponseCode(t, http.StatusOK, resp.Code)
	err = json.Unmarshal(resp.Body.Bytes(), &res)
	require.NoError(t, err)
	require.Equal(t, sampleUserID, res.ID)

	// get that user and verfiy the update
	resp = getTestUser(t, sampleUserID)
	checkResponseCode(t, http.StatusOK, resp.Code)
	var getResp GetUserResponse
	err = json.Unmarshal(resp.Body.Bytes(), &getResp)
	require.NoError(t, err)
	require.NotNil(t, getResp.User)
	require.Equal(t, sampleUserID, getResp.User.ID)
	require.Equal(t, "tiny cat", getResp.User.Name)
	require.Equal(t, "tiny@enterprisecatz.com", getResp.User.Email)
}

func TestDeleteUser(t *testing.T) {
	clearTable()

	// create a user
	uuidGenerator.shouldGenUserID = true
	resp := createTestUser(t, "tiny cat", "tiny@cat.com")
	checkResponseCode(t, http.StatusOK, resp.Code)
	var res CreateUserResponse
	err := json.Unmarshal(resp.Body.Bytes(), &res)
	require.NoError(t, err)
	require.Equal(t, sampleUserID, res.ID)

	// delete the user
	resp = deleteTestUser(t, sampleUserID)
	checkResponseCode(t, http.StatusOK, resp.Code)
	var deleteRes DeleteUserResponse
	err = json.Unmarshal(resp.Body.Bytes(), &deleteRes)
	require.NoError(t, err)
	require.Equal(t, sampleUserID, deleteRes.ID)

	// try to delete a non-existant user and verify there is an error
	resp = deleteTestUser(t, sampleUserID)
	require.Equal(t, http.StatusInternalServerError, resp.Result().StatusCode)
	checkResponseCode(t, http.StatusInternalServerError, resp.Result().StatusCode)
}

func TestGetPost_EmptyTable(t *testing.T) {
	clearTable()

	// try to get a post for a non-existant user
	resp := getTestPost(t, samplePostID)
	checkResponseCode(t, http.StatusOK, resp.Code)
	var getPostResponse GetPostResponse
	err := json.Unmarshal(resp.Body.Bytes(), &getPostResponse)
	require.NoError(t, err)
	require.Nil(t, getPostResponse.Post)

	// create the user that has no posts saved yet
	uuidGenerator.shouldGenUserID = true
	resp = createTestUser(t, "tiny cat", "tiny@cat.com")
	checkResponseCode(t, http.StatusOK, resp.Code)
	var res CreateUserResponse
	err = json.Unmarshal(resp.Body.Bytes(), &res)
	require.NoError(t, err)
	require.Equal(t, sampleUserID, res.ID)

	// try to get a non-existant post for the user
	resp = getTestPost(t, samplePostID)
	checkResponseCode(t, http.StatusOK, resp.Result().StatusCode)
	var getRes GetPostResponse
	err = json.Unmarshal(resp.Body.Bytes(), &getRes)
	require.NoError(t, err)
	require.Nil(t, getRes.Post)
}

func TestCreateOrUpdatePost(t *testing.T) {
	clearTable()

	// create a new user
	uuidGenerator.shouldGenUserID = true
	resp := createTestUser(t, "tiny cat", "tiny@cat.com")
	checkResponseCode(t, http.StatusOK, resp.Code)
	var userRes CreateUserResponse
	err := json.Unmarshal(resp.Body.Bytes(), &userRes)
	require.NoError(t, err)
	require.Equal(t, sampleUserID, userRes.ID)

	// create a new post for that user
	uuidGenerator.shouldGenUserID = false
	uuidGenerator.shouldGenPostID = true
	resp = createTestPost(t, sampleUserID, "title", "content")
	checkResponseCode(t, http.StatusOK, resp.Code)
	var res CreatePostResponse
	err = json.Unmarshal(resp.Body.Bytes(), &res)
	require.NoError(t, err)

	// get the post and verify the creation
	resp = getTestPost(t, samplePostID)
	checkResponseCode(t, http.StatusOK, resp.Code)
	var getRes GetPostResponse
	err = json.Unmarshal(resp.Body.Bytes(), &getRes)
	require.NoError(t, err)
	require.Equal(t, "title", getRes.Post.Title)
	require.Equal(t, "content", getRes.Post.Content)

	// update the existing post
	resp = updateTestPost(t, samplePostID, "updated title", "updated content")
	checkResponseCode(t, http.StatusOK, resp.Code)
	err = json.Unmarshal(resp.Body.Bytes(), &res)
	require.NoError(t, err)

	// get the post to verify the updates
	resp = getTestPost(t, samplePostID)
	checkResponseCode(t, http.StatusOK, resp.Code)
	err = json.Unmarshal(resp.Body.Bytes(), &getRes)
	require.NoError(t, err)
	require.Equal(t, samplePostID, getRes.Post.ID)
	require.Equal(t, sampleUserID, getRes.Post.UserID)
	require.Equal(t, "updated title", getRes.Post.Title)
	require.Equal(t, "updated content", getRes.Post.Content)
}

func TestDeletePost(t *testing.T) {
	clearTable()

	// create a new user
	uuidGenerator.shouldGenUserID = true
	resp := createTestUser(t, "tiny cat", "tiny@cat.com")
	checkResponseCode(t, http.StatusOK, resp.Code)
	var userRes CreateUserResponse
	err := json.Unmarshal(resp.Body.Bytes(), &userRes)
	require.NoError(t, err)
	require.Equal(t, sampleUserID, userRes.ID)

	// create a post for that user
	uuidGenerator.shouldGenUserID = false
	uuidGenerator.shouldGenPostID = true
	resp = createTestPost(t, sampleUserID, "title", "content")
	checkResponseCode(t, http.StatusOK, resp.Code)
	var res CreatePostResponse
	err = json.Unmarshal(resp.Body.Bytes(), &res)
	require.NoError(t, err)

	// get the post and verify the creation
	resp = getTestPost(t, samplePostID)
	checkResponseCode(t, http.StatusOK, resp.Code)
	var getRes GetPostResponse
	err = json.Unmarshal(resp.Body.Bytes(), &getRes)
	require.NoError(t, err)
	require.Equal(t, "title", getRes.Post.Title)
	require.Equal(t, "content", getRes.Post.Content)

	// delete the post
	resp = deleteTestPost(t, samplePostID)
	checkResponseCode(t, http.StatusOK, resp.Code)
	var deleteRes DeletePostResponse
	err = json.Unmarshal(resp.Body.Bytes(), &deleteRes)
	require.NoError(t, err)
	require.Equal(t, samplePostID, deleteRes.ID)

	// try and get the deleted post
	resp = getTestPost(t, samplePostID)
	checkResponseCode(t, http.StatusOK, resp.Result().StatusCode)
	err = json.Unmarshal(resp.Body.Bytes(), &getRes)
	require.NoError(t, err)
	require.Nil(t, getRes.Post)

	// try to delete a post that doesn't exist
	resp = deleteTestPost(t, samplePostID2)
	checkResponseCode(t, http.StatusInternalServerError, resp.Result().StatusCode)

	// try to delete a post for a user that doesn't exist
	resp = deleteTestPost(t, samplePostID)
	checkResponseCode(t, http.StatusInternalServerError, resp.Result().StatusCode)
}

func deleteTestUser(t *testing.T, id string) *httptest.ResponseRecorder {
	reqBytes, err := json.Marshal(&DeleteUserRequest{
		ID: id,
	})
	require.NoError(t, err)
	req, err := http.NewRequest("DELETE", "/users", bytes.NewBuffer(reqBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return executeRequest(req)
}

func updateTestUser(t *testing.T, id, name, email string) *httptest.ResponseRecorder {
	reqBytes, err := json.Marshal(&UpdateUserRequest{
		ID:    id,
		Name:  name,
		Email: email,
	})
	require.NoError(t, err)
	req, err := http.NewRequest("PUT", "/users", bytes.NewBuffer(reqBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return executeRequest(req)
}

func createTestUser(t *testing.T, name, email string) *httptest.ResponseRecorder {
	reqBytes, err := json.Marshal(&CreateUserRequest{
		Name:  name,
		Email: email,
	})
	require.NoError(t, err)
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return executeRequest(req)
}

func getTestUser(t *testing.T, id string) *httptest.ResponseRecorder {
	reqBytes, err := json.Marshal(&GetUserRequest{ID: id})
	require.NoError(t, err)
	req, err := http.NewRequest("GET", "/users", bytes.NewBuffer(reqBytes))
	require.NoError(t, err)
	return executeRequest(req)
}

func createTestPost(t *testing.T, userID, title, content string) *httptest.ResponseRecorder {
	reqBytes, err := json.Marshal(&CreatePostRequest{
		UserID:  userID,
		Title:   title,
		Content: content,
	})
	require.NoError(t, err)
	req, err := http.NewRequest("POST", "/posts", bytes.NewBuffer(reqBytes))
	require.NoError(t, err)
	return executeRequest(req)
}

func updateTestPost(t *testing.T, id, title, content string) *httptest.ResponseRecorder {
	reqBytes, err := json.Marshal(&UpdatePostRequest{
		ID:      id,
		Title:   title,
		Content: content,
	})
	require.NoError(t, err)
	req, err := http.NewRequest("PUT", "/posts", bytes.NewBuffer(reqBytes))
	require.NoError(t, err)
	return executeRequest(req)
}

func getTestPost(t *testing.T, postID string) *httptest.ResponseRecorder {
	reqBytes, err := json.Marshal(&GetPostRequest{
		ID: postID,
	})
	require.NoError(t, err)
	req, err := http.NewRequest("GET", "/posts", bytes.NewBuffer(reqBytes))
	require.NoError(t, err)
	return executeRequest(req)
}

func deleteTestPost(t *testing.T, id string) *httptest.ResponseRecorder {
	reqBytes, err := json.Marshal(&DeletePostRequest{
		ID: id,
	})
	require.NoError(t, err)
	req, err := http.NewRequest("DELETE", "/posts", bytes.NewBuffer(reqBytes))
	require.NoError(t, err)
	return executeRequest(req)
}
