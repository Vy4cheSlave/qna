package rest

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Vy4cheSlave/qna/internal/domain"
	"github.com/Vy4cheSlave/qna/internal/infrastructure/rest/dto/request"
	"github.com/Vy4cheSlave/qna/internal/infrastructure/rest/dto/response"
	"github.com/Vy4cheSlave/qna/internal/infrastructure/rest/middleware"
	"github.com/google/uuid"

	"github.com/pkg/errors"
)

type QNADispatcher interface {
	CreateUser(ctx context.Context, userName *string) (userId *string, err error)
	GetUsers(ctx context.Context) (*[]domain.User, error)
	DeleteUser(ctx context.Context, userId *string) error
	GetQuestions(ctx context.Context) (*[]domain.Question, error)
	CreateQuestion(ctx context.Context, question *string) (questionId int, err error)
	GetQuestionAndAnswers(ctx context.Context, questionId int) (*domain.Question, *[]domain.Answer, error)
	DeleteQuestionAndAnswers(ctx context.Context, questionId int) error
	CreateAnswerToQuestion(ctx context.Context, answer *domain.Answer) (answerId int, err error)
	GetAnswer(ctx context.Context, answerId int) (*domain.Answer, error)
	DeleteAnswer(ctx context.Context, answerId int) error
}

type Server struct {
	log        *slog.Logger
	service    QNADispatcher
	restServer *http.Server
	addr       *string
}

type serverAPI struct {
	addr    *string
	log     *slog.Logger
	service QNADispatcher
}

func NewServer(log *slog.Logger, service QNADispatcher, addr *string) *Server {
	restServer := NewRestServer(&serverAPI{addr: addr, log: log, service: service})
	return &Server{
		log:        log,
		restServer: restServer,
		service:    service,
		addr:       addr,
	}
}

func NewRestServer(api *serverAPI) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /users/", api.CreateUser)
	mux.HandleFunc("GET /users/", api.GetUsers)
	mux.HandleFunc("DELETE /users/{id}", api.DeleteUser)
	mux.HandleFunc("GET /questions/{id}", api.GetQuestionAndAnswers)
	mux.HandleFunc("DELETE /questions/{id}", api.DeleteQuestionAndAnswers)
	mux.HandleFunc("POST /questions/{id}/answers/", api.CreateAnswerToQuestion)
	mux.HandleFunc("GET /questions/", api.GetQuestions)
	mux.HandleFunc("POST /questions/", api.CreateQuestion)
	mux.HandleFunc("GET /answers/{id}", api.GetAnswer)
	mux.HandleFunc("DELETE /answers/{id}", api.DeleteAnswer)

	var handler http.Handler = mux
	handler = middleware.CORSMiddleware(handler)
	handler = middleware.LoggMiddleware(api.log, handler)

	server := &http.Server{
		Addr:         *api.addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server
}

func (t *Server) Run() error {
	const op = "internal/infrastructure/rest/handler.Server.Run"
	log := t.log.With(slog.String("operation", op), slog.String("addr", *t.addr))
	log.Info("server is running")

	if err := t.restServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, strings.Join([]string{op, "failed to serve rest server"}, ": "))
	}

	return nil
}

func (t *serverAPI) CreateUser(w http.ResponseWriter, r *http.Request) {
	var errorList []error
	ctx := r.Context()

	var req request.CreateUserRequest

	// Десериализация JSON-запроса
	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()
	if err != nil {
		errorList = append(errorList, err)
		err := response.ReturnResponse(
			w,
			http.StatusBadRequest,
			response.WithError(response.ErrCodeJsonParsingFailed, "Invalid request body"),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusBadRequest, &errorList)
		return
	}

	// Валидация входных данных
	if len(req.Name) == 0 {
		errorList = append(errorList, errors.New("field \"name\" must not be empty"))
		err := response.ReturnResponse(
			w,
			http.StatusBadRequest,
			response.WithError(response.ErrCodeValidationFailed, "field \"name\" must not be empty"),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusBadRequest, &errorList)
		return
	}

	// Вызов метода сервиса
	userId, err := t.service.CreateUser(ctx, &req.Name)
	if err != nil {
		errorList = append(errorList, err)
		err := response.ReturnResponse(
			w,
			http.StatusInternalServerError,
			response.WithError(response.ErrCodeInternalServerError, "internal server error"),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusInternalServerError, &errorList)
		return
	}

	// Формирование ответа
	err = response.ReturnResponse(
		w,
		http.StatusOK,
		response.WithData(response.CreateUserResponse{UserId: *userId}),
	)
	if err != nil {
		errorList = append(errorList, err)
	}
	middleware.UpdateContext(ctx, r, http.StatusOK, &errorList)
}

func (t *serverAPI) GetUsers(w http.ResponseWriter, r *http.Request) {
	var errorList []error
	ctx := r.Context()

	// Вызов метода сервиса
	users, err := t.service.GetUsers(ctx)
	if err != nil {
		errorList = append(errorList, err)
		err := response.ReturnResponse(
			w,
			http.StatusInternalServerError,
			response.WithError(response.ErrCodeInternalServerError, "internal server error"),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusInternalServerError, &errorList)
		return
	}

	// Формирование ответа
	err = response.ReturnResponse(
		w,
		http.StatusOK,
		response.WithData(*users),
	)
	if err != nil {
		errorList = append(errorList, err)
	}
	middleware.UpdateContext(ctx, r, http.StatusOK, &errorList)
}

func (t *serverAPI) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var errorList []error
	ctx := r.Context()

	userId := r.PathValue("id")

	// Валидация входных данных
	_, err := uuid.Parse(userId)
	if err != nil {
		errorList = append(errorList, err)
		err := response.ReturnResponse(
			w,
			http.StatusBadRequest,
			response.WithError(response.ErrCodeValidationFailed, "invalid UUID format for \"id\""),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusBadRequest, &errorList)
		return
	}

	// Вызов метода сервиса
	err = t.service.DeleteUser(ctx, &userId)
	if err != nil {
		errorList = append(errorList, err)
		err := response.ReturnResponse(
			w,
			http.StatusInternalServerError,
			response.WithError(response.ErrCodeInternalServerError, "internal server error"),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusInternalServerError, &errorList)
		return
	}

	// Формирование ответа
	err = response.ReturnResponse(
		w,
		http.StatusOK,
	)
	if err != nil {
		errorList = append(errorList, err)
	}
	middleware.UpdateContext(ctx, r, http.StatusOK, &errorList)
}

func (t *serverAPI) GetQuestions(w http.ResponseWriter, r *http.Request) {
	var errorList []error
	ctx := r.Context()

	// Вызов метода сервиса
	questions, err := t.service.GetQuestions(ctx)
	if err != nil {
		errorList = append(errorList, err)
		err := response.ReturnResponse(
			w,
			http.StatusInternalServerError,
			response.WithError(response.ErrCodeInternalServerError, "internal server error"),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusInternalServerError, &errorList)
		return
	}

	// Формирование ответа
	err = response.ReturnResponse(
		w,
		http.StatusOK,
		response.WithData(*questions),
	)
	if err != nil {
		errorList = append(errorList, err)
	}
	middleware.UpdateContext(ctx, r, http.StatusOK, &errorList)
}

func (t *serverAPI) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	var errorList []error
	ctx := r.Context()

	var req request.CreateQuestionRequest

	// Десериализация JSON-запроса
	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()
	if err != nil {
		errorList = append(errorList, err)
		err := response.ReturnResponse(
			w,
			http.StatusBadRequest,
			response.WithError(response.ErrCodeJsonParsingFailed, "Invalid request body"),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusBadRequest, &errorList)
		return
	}

	// Валидация входных данных
	if len(req.Text) == 0 {
		errorList = append(errorList, errors.New("field \"text\" must not be empty"))
		err := response.ReturnResponse(
			w,
			http.StatusBadRequest,
			response.WithError(response.ErrCodeValidationFailed, "field \"text\" must not be empty"),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusBadRequest, &errorList)
		return
	}

	// Вызов метода сервиса
	questionId, err := t.service.CreateQuestion(ctx, &req.Text)
	if err != nil {
		errorList = append(errorList, err)
		err := response.ReturnResponse(
			w,
			http.StatusInternalServerError,
			response.WithError(response.ErrCodeInternalServerError, "internal server error"),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusInternalServerError, &errorList)
		return
	}

	// Формирование ответа
	err = response.ReturnResponse(
		w,
		http.StatusOK,
		response.WithData(response.CreateQuestionResponse{QuestionId: questionId}),
	)
	if err != nil {
		errorList = append(errorList, err)
	}
	middleware.UpdateContext(ctx, r, http.StatusOK, &errorList)
}

func (t *serverAPI) GetQuestionAndAnswers(w http.ResponseWriter, r *http.Request) {
	var errorList []error
	ctx := r.Context()

	questionId := r.PathValue("id")

	// Валидация входных данных
	questionIdInt, err := strconv.Atoi(questionId)
	if err != nil {
		errorList = append(errorList, err)
		err := response.ReturnResponse(
			w,
			http.StatusBadRequest,
			response.WithError(response.ErrCodeValidationFailed, "invalid ID format for \"id\""),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusBadRequest, &errorList)
		return
	}
	if questionIdInt < 1 {
		errorList = append(errorList, errors.New("ID must be a positive integer"))
		err := response.ReturnResponse(
			w,
			http.StatusBadRequest,
			response.WithError(response.ErrCodeValidationFailed, "ID must be a positive integer"),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusBadRequest, &errorList)
		return
	}

	// Вызов метода сервиса
	question, answers, err := t.service.GetQuestionAndAnswers(ctx, questionIdInt)
	if err != nil {
		errorList = append(errorList, err)
		err := response.ReturnResponse(
			w,
			http.StatusInternalServerError,
			response.WithError(response.ErrCodeInternalServerError, "internal server error"),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusInternalServerError, &errorList)
		return
	}

	// Формирование ответа
	err = response.ReturnResponse(
		w,
		http.StatusOK,
		response.WithData(response.GetQuestionAndAnswersResponse{Question: *question, Answers: *answers}),
	)
	if err != nil {
		errorList = append(errorList, err)
	}
	middleware.UpdateContext(ctx, r, http.StatusOK, &errorList)
}

func (t *serverAPI) DeleteQuestionAndAnswers(w http.ResponseWriter, r *http.Request) {
	var errorList []error
	ctx := r.Context()

	questionId := r.PathValue("id")

	// Валидация входных данных
	questionIdInt, err := strconv.Atoi(questionId)
	if err != nil {
		errorList = append(errorList, err)
		err := response.ReturnResponse(
			w,
			http.StatusBadRequest,
			response.WithError(response.ErrCodeValidationFailed, "invalid ID format for \"id\""),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusBadRequest, &errorList)
		return
	}
	if questionIdInt < 1 {
		errorList = append(errorList, errors.New("ID must be a positive integer"))
		err := response.ReturnResponse(
			w,
			http.StatusBadRequest,
			response.WithError(response.ErrCodeValidationFailed, "ID must be a positive integer"),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusBadRequest, &errorList)
		return
	}

	// Вызов метода сервиса
	err = t.service.DeleteQuestionAndAnswers(ctx, questionIdInt)
	if err != nil {
		errorList = append(errorList, err)
		err := response.ReturnResponse(
			w,
			http.StatusInternalServerError,
			response.WithError(response.ErrCodeInternalServerError, "internal server error"),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusInternalServerError, &errorList)
		return
	}

	// Формирование ответа
	err = response.ReturnResponse(
		w,
		http.StatusOK,
	)
	if err != nil {
		errorList = append(errorList, err)
	}
	middleware.UpdateContext(ctx, r, http.StatusOK, &errorList)
}

func (t *serverAPI) CreateAnswerToQuestion(w http.ResponseWriter, r *http.Request) {
	var errorList []error
	ctx := r.Context()

	var req request.CreateAnswerToQuestionRequest
	questionId := r.PathValue("id")

	// Десериализация JSON-запроса
	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()
	if err != nil {
		errorList = append(errorList, err)
		err := response.ReturnResponse(
			w,
			http.StatusBadRequest,
			response.WithError(response.ErrCodeJsonParsingFailed, "Invalid request body"),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusBadRequest, &errorList)
		return
	}

	// Валидация входных данных
	questionIdInt, err := strconv.Atoi(questionId)
	if err != nil {
		errorList = append(errorList, err)
		err := response.ReturnResponse(
			w,
			http.StatusBadRequest,
			response.WithError(response.ErrCodeValidationFailed, "invalid ID format for \"id\""),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusBadRequest, &errorList)
		return
	}
	if questionIdInt < 1 {
		errorList = append(errorList, errors.New("ID must be a positive integer"))
		err := response.ReturnResponse(
			w,
			http.StatusBadRequest,
			response.WithError(response.ErrCodeValidationFailed, "ID must be a positive integer"),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusBadRequest, &errorList)
		return
	}
	if len(req.Text) == 0 {
		errorList = append(errorList, errors.New("field \"text\" must not be empty"))
		err := response.ReturnResponse(
			w,
			http.StatusBadRequest,
			response.WithError(response.ErrCodeValidationFailed, "field \"text\" must not be empty"),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusBadRequest, &errorList)
		return
	}
	_, err = uuid.Parse(req.UserId)
	if err != nil {
		errorList = append(errorList, err)
		err := response.ReturnResponse(
			w,
			http.StatusBadRequest,
			response.WithError(response.ErrCodeValidationFailed, "invalid UUID format for \"user_id\""),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusBadRequest, &errorList)
		return
	}

	// Вызов метода сервиса
	answer := domain.Answer{
		QuestionId: questionIdInt,
		UserId:     req.UserId,
		Text:       req.Text,
	}
	answerId, err := t.service.CreateAnswerToQuestion(ctx, &answer)
	if err != nil {
		errorList = append(errorList, err)
		err := response.ReturnResponse(
			w,
			http.StatusInternalServerError,
			response.WithError(response.ErrCodeInternalServerError, "internal server error"),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusInternalServerError, &errorList)
		return
	}

	// Формирование ответа
	err = response.ReturnResponse(
		w,
		http.StatusOK,
		response.WithData(response.CreateAnswerToQuestionResponse{AnswerId: answerId}),
	)
	if err != nil {
		errorList = append(errorList, err)
	}
	middleware.UpdateContext(ctx, r, http.StatusOK, &errorList)
}

func (t *serverAPI) GetAnswer(w http.ResponseWriter, r *http.Request) {
	var errorList []error
	ctx := r.Context()

	answerId := r.PathValue("id")

	// Валидация входных данных
	answerIdInt, err := strconv.Atoi(answerId)
	if err != nil {
		errorList = append(errorList, err)
		err := response.ReturnResponse(
			w,
			http.StatusBadRequest,
			response.WithError(response.ErrCodeValidationFailed, "invalid ID format for \"id\""),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusBadRequest, &errorList)
		return
	}
	if answerIdInt < 1 {
		errorList = append(errorList, errors.New("ID must be a positive integer"))
		err := response.ReturnResponse(
			w,
			http.StatusBadRequest,
			response.WithError(response.ErrCodeValidationFailed, "ID must be a positive integer"),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusBadRequest, &errorList)
		return
	}

	// Вызов метода сервиса
	answer, err := t.service.GetAnswer(ctx, answerIdInt)
	if err != nil {
		errorList = append(errorList, err)
		err := response.ReturnResponse(
			w,
			http.StatusInternalServerError,
			response.WithError(response.ErrCodeInternalServerError, "internal server error"),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusInternalServerError, &errorList)
		return
	}

	// Формирование ответа
	err = response.ReturnResponse(
		w,
		http.StatusOK,
		response.WithData(answer),
	)
	if err != nil {
		errorList = append(errorList, err)
	}
	middleware.UpdateContext(ctx, r, http.StatusOK, &errorList)
}

func (t *serverAPI) DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	var errorList []error
	ctx := r.Context()

	answerId := r.PathValue("id")

	// Валидация входных данных
	answerIdInt, err := strconv.Atoi(answerId)
	if err != nil {
		errorList = append(errorList, err)
		err := response.ReturnResponse(
			w,
			http.StatusBadRequest,
			response.WithError(response.ErrCodeValidationFailed, "invalid ID format for \"id\""),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusBadRequest, &errorList)
		return
	}
	if answerIdInt < 1 {
		errorList = append(errorList, errors.New("ID must be a positive integer"))
		err := response.ReturnResponse(
			w,
			http.StatusBadRequest,
			response.WithError(response.ErrCodeValidationFailed, "ID must be a positive integer"),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusBadRequest, &errorList)
		return
	}

	// Вызов метода сервиса
	err = t.service.DeleteAnswer(ctx, answerIdInt)
	if err != nil {
		errorList = append(errorList, err)
		err := response.ReturnResponse(
			w,
			http.StatusInternalServerError,
			response.WithError(response.ErrCodeInternalServerError, "internal server error"),
		)
		if err != nil {
			errorList = append(errorList, err)
		}
		middleware.UpdateContext(ctx, r, http.StatusInternalServerError, &errorList)
		return
	}

	// Формирование ответа
	err = response.ReturnResponse(
		w,
		http.StatusOK,
	)
	if err != nil {
		errorList = append(errorList, err)
	}
	middleware.UpdateContext(ctx, r, http.StatusOK, &errorList)
}
