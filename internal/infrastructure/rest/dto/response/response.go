package response

import (
	"encoding/json"
	"net/http"

	"github.com/Vy4cheSlave/qna/internal/domain"
	"github.com/pkg/errors"
)

// Error Codes
const (
	// ErrCodeInvalidToken        = "INVALID_TOKEN"
	// ErrCodeMissingAuthHeader   = "MISSING_AUTH_HEADER"
	ErrCodeValidationFailed    = "VALIDATION_FAILED"
	ErrCodeJsonParsingFailed   = "JSON_PARSING_FAILED"
	ErrCodeInternalServerError = "INTERNAL_SERVER_ERROR"
	// ErrCodeUnauthorized        = "UNAUTHORIZED"
	// ErrCodeNotFound            = "NOT_FOUND"
)

type Response struct {
	Status string `json:"status"`
	Error  *Error `json:"error,omitempty"`
	Data   any    `json:"data,omitempty"`
}

type Error struct {
	Code string `json:"code"`
	Desc string `json:"desc,omitempty"`
}

func ReturnResponse(w http.ResponseWriter, httpStatus int, opts ...Option) error {
	resp := &Response{Status: http.StatusText(httpStatus)}

	for _, opt := range opts {
		opt(resp)
	}

	w.WriteHeader(httpStatus)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return errors.Wrap(err, "failed to encode and write JSON response")
	}
	return nil
}

type Option func(*Response)

func WithError(code, desc string) Option {
	return func(r *Response) {
		r.Error = &Error{Code: code, Desc: desc}
	}
}

func WithData(data any) Option {
	return func(r *Response) {
		r.Data = data
	}
}

// схемы
type CreateUserResponse struct {
	UserId string `json:"user_id"`
}

type CreateQuestionResponse struct {
	QuestionId int `json:"question_id"`
}

type GetQuestionAndAnswersResponse struct {
	Question domain.Question `json:"question"`
	Answers  []domain.Answer `json:"answers"`
}

type CreateAnswerToQuestionResponse struct {
	AnswerId int `json:"answer_id"`
}
