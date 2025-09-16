package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/maulana1k/forum-app/internal/app/dto"
	shared_suite "github.com/maulana1k/forum-app/tests/suite"
)

// --------------------------
// PostHandler Test Suite
// --------------------------
type PostHandlerTestSuite struct {
	suite.Suite
	App   *fiber.App
	Token string
}

func (s *PostHandlerTestSuite) SetupSuite() {
	shared := shared_suite.GetSharedSuite()
	s.App = shared.App
	s.Token = shared.Token
}

func (s *PostHandlerTestSuite) TearDownSuite() {
	shared_suite.TeardownSharedSuite()
}

// --------------------------
// Helpers
// --------------------------
func createPostPayload(content, tags string) []byte {
	req := dto.CreatePostRequest{
		Content: content,
		Tags:    tags,
	}
	data, _ := json.Marshal(req)
	return data
}

func parsePostResponse(t *testing.T, resp *http.Response) dto.PostResponse {
	var post dto.PostResponse
	if err := json.NewDecoder(resp.Body).Decode(&post); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	return post
}

// --------------------------
// Tests
// --------------------------
func (s *PostHandlerTestSuite) TestCreatePost() {
	payload := createPostPayload("Hello singleton", "golang,testing")
	req := httptest.NewRequest(http.MethodPost, "/v1/posts/", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.Token)

	resp, err := s.App.Test(req, -1)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusCreated, resp.StatusCode)

	post := parsePostResponse(s.T(), resp)
	assert.Equal(s.T(), "Hello singleton", post.Content)
}

// --------------------------
// Entry point
// --------------------------
func TestPostHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(PostHandlerTestSuite))
}
