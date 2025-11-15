package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Vy4cheSlave/qna/internal/domain"
	"github.com/Vy4cheSlave/qna/internal/infrastructure/rest/dto/response"
	"github.com/Vy4cheSlave/qna/internal/infrastructure/rest/mocks"
)

type testCase struct {
	name           string
	requestBody    string
	requestPath    string
	setupMock      func(*mocks.MockQNADispatcher)
	expectedStatus int
	expectedResp   interface{}
}

func TestCreateAnswerToQuestion(t *testing.T) {
	testCases := []testCase{
		{
			name:        "Success",
			requestBody: `{"user_id": "f47ac10b-58cc-4372-a567-0e02b2c3de91", "text": "text"}`,
			requestPath: "1",
			setupMock: func(mockDispatcher *mocks.MockQNADispatcher) {
				mockDispatcher.On("CreateAnswerToQuestion",
					mock.Anything,
					&domain.Answer{
						UserId:     "f47ac10b-58cc-4372-a567-0e02b2c3de91",
						QuestionId: 1,
						Text:       "text",
					},
				).Return(1, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedResp: map[string]interface{}{
				"data": map[string]interface{}{
					"answer_id": float64(1),
				},
				"status": http.StatusText(http.StatusOK),
			},
		},
		{
			name:           "Invalid request body",
			requestBody:    `{"user_id": "f47ac10b-58cc-4372-a567-0e02b2c3de91", "text": "text"`,
			requestPath:    "1",
			setupMock:      func(mockDispatcher *mocks.MockQNADispatcher) {},
			expectedStatus: http.StatusBadRequest,
			expectedResp: map[string]interface{}{
				"error": map[string]interface{}{
					"code": response.ErrCodeJsonParsingFailed,
					"desc": "Invalid request body",
				},
				"status": http.StatusText(http.StatusBadRequest),
			},
		},
		{
			name:           "invalid ID format",
			requestBody:    `{"user_id": "f47ac10b-58cc-4372-a567-0e02b2c3de91", "text": "text"}`,
			requestPath:    "slovo",
			setupMock:      func(mockDispatcher *mocks.MockQNADispatcher) {},
			expectedStatus: http.StatusBadRequest,
			expectedResp: map[string]interface{}{
				"error": map[string]interface{}{
					"code": response.ErrCodeValidationFailed,
					"desc": "invalid ID format for \"id\"",
				},
				"status": http.StatusText(http.StatusBadRequest),
			},
		},
		{
			name:           "field \"text\" must not be empty",
			requestBody:    `{"user_id": "f47ac10b-58cc-4372-a567-0e02b2c3de91", "text": ""}`,
			requestPath:    "1",
			setupMock:      func(mockDispatcher *mocks.MockQNADispatcher) {},
			expectedStatus: http.StatusBadRequest,
			expectedResp: map[string]interface{}{
				"error": map[string]interface{}{
					"code": response.ErrCodeValidationFailed,
					"desc": "field \"text\" must not be empty",
				},
				"status": http.StatusText(http.StatusBadRequest),
			},
		},
		{
			name:           "invalid UUID format for \"user_id\"",
			requestBody:    `{"user_id": "not-uuid", "text": "text"}`,
			requestPath:    "1",
			setupMock:      func(mockDispatcher *mocks.MockQNADispatcher) {},
			expectedStatus: http.StatusBadRequest,
			expectedResp: map[string]interface{}{
				"error": map[string]interface{}{
					"code": response.ErrCodeValidationFailed,
					"desc": "invalid UUID format for \"user_id\"",
				},
				"status": http.StatusText(http.StatusBadRequest),
			},
		},
		{
			name:        "internal server error",
			requestBody: `{"user_id": "f47ac10b-58cc-4372-a567-0e02b2c3de91", "text": "text"}`,
			requestPath: "1",
			setupMock: func(mockDispatcher *mocks.MockQNADispatcher) {
				mockDispatcher.On("CreateAnswerToQuestion",
					mock.Anything,
					&domain.Answer{
						UserId:     "f47ac10b-58cc-4372-a567-0e02b2c3de91",
						QuestionId: 1,
						Text:       "text",
					},
				).Return(0, errors.New("some error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResp: map[string]interface{}{
				"error": map[string]interface{}{
					"code": response.ErrCodeInternalServerError,
					"desc": "internal server error",
				},
				"status": http.StatusText(http.StatusInternalServerError),
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockQNADispatcher := mocks.NewMockQNADispatcher(t)
			tt.setupMock(mockQNADispatcher)

			handler := &serverAPI{
				addr:    nil,
				service: mockQNADispatcher,
				log:     slog.Default(),
			}

			req := httptest.NewRequest(
				http.MethodPost,
				fmt.Sprintf("/questions/%s/answers/", tt.requestPath),
				bytes.NewBufferString(tt.requestBody),
			)
			req.SetPathValue("id", tt.requestPath)

			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handler.CreateAnswerToQuestion(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody map[string]interface{}
			err := json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedResp, responseBody)

			mockQNADispatcher.AssertExpectations(t)
		})
	}
}
