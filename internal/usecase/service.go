package usecase

import (
	"context"
	"github.com/Vy4cheSlave/qna/internal/domain"
	"github.com/pkg/errors"
)

type QNAManager interface {
	ReadQuestions(ctx context.Context) (*[]domain.Question, error)
	CreateQuestion(ctx context.Context, question *string) (questionId int, err error)
	ReadQuestionAndAnswers(ctx context.Context, questionId int) (*domain.Question, *[]domain.Answer, error)
	DeleteQuestionAndAnswers(ctx context.Context, questionId int) error
	CreateAnswerToQuestion(ctx context.Context, answer *domain.Answer) (answerId int, err error)
	ReadAnswer(ctx context.Context, answerId int) (*domain.Answer, error)
	DeleteAnswer(ctx context.Context, answerId int) error
}

type UserManager interface {
	CreateUser(ctx context.Context, userName *string) (userId *string, err error)
	ReadUsers(ctx context.Context) (*[]domain.User, error)
	DeleteUser(ctx context.Context, userId *string) error
}

type QNACrud struct {
	qnaManager  QNAManager
	userManager UserManager
}

func NewQNAManagerService(qnaManager QNAManager, userManager UserManager) *QNACrud {
	return &QNACrud{
		qnaManager:  qnaManager,
		userManager: userManager,
	}
}

func (t *QNACrud) CreateUser(ctx context.Context, userName *string) (userId *string, err error) {
	const op = "internal/usecase/service.QNACrud.CreateUser"

	userId, err = t.userManager.CreateUser(ctx, userName)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	return userId, nil
}

func (t *QNACrud) GetUsers(ctx context.Context) (*[]domain.User, error) {
	const op = "internal/usecase/service.QNACrud.GetUsers"

	users, err := t.userManager.ReadUsers(ctx)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	return users, nil
}

func (t *QNACrud) DeleteUser(ctx context.Context, userId *string) error {
	const op = "internal/usecase/service.QNACrud.GetUsers"

	err := t.userManager.DeleteUser(ctx, userId)
	if err != nil {
		return errors.Wrap(err, op)
	}
	return nil
}

func (t *QNACrud) GetQuestions(ctx context.Context) (*[]domain.Question, error) {
	const op = "internal/usecase/service.QNACrud.GetQuestions"

	questions, err := t.qnaManager.ReadQuestions(ctx)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	return questions, nil
}

func (t *QNACrud) CreateQuestion(ctx context.Context, question *string) (questionId int, err error) {
	const op = "internal/usecase/service.QNACrud.CreateQuestion"

	questionId, err = t.qnaManager.CreateQuestion(ctx, question)
	if err != nil {
		return 0, errors.Wrap(err, op)
	}
	return questionId, nil
}

func (t *QNACrud) GetQuestionAndAnswers(ctx context.Context, questionId int) (*domain.Question, *[]domain.Answer, error) {
	const op = "internal/usecase/service.QNACrud.GetQuestionAndAnswers"

	question, answers, err := t.qnaManager.ReadQuestionAndAnswers(ctx, questionId)
	if err != nil {
		return nil, nil, errors.Wrap(err, op)
	}
	return question, answers, nil
}

func (t *QNACrud) DeleteQuestionAndAnswers(ctx context.Context, questionId int) error {
	const op = "internal/usecase/service.QNACrud.DeleteQuestionAndAnswers"

	err := t.qnaManager.DeleteQuestionAndAnswers(ctx, questionId)
	if err != nil {
		return errors.Wrap(err, op)
	}
	return nil
}

func (t *QNACrud) CreateAnswerToQuestion(ctx context.Context, answer *domain.Answer) (answerId int, err error) {
	const op = "internal/usecase/service.QNACrud.CreateAnswerToQuestion"

	answerId, err = t.qnaManager.CreateAnswerToQuestion(ctx, answer)
	if err != nil {
		return 0, errors.Wrap(err, op)
	}
	return answerId, nil
}

func (t *QNACrud) GetAnswer(ctx context.Context, answerId int) (*domain.Answer, error) {
	const op = "internal/usecase/service.QNACrud.GetAnswer"

	answers, err := t.qnaManager.ReadAnswer(ctx, answerId)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	return answers, nil
}

func (t *QNACrud) DeleteAnswer(ctx context.Context, answerId int) error {
	const op = "internal/usecase/service.QNACrud.DeleteAnswer"

	err := t.qnaManager.DeleteAnswer(ctx, answerId)
	if err != nil {
		return errors.Wrap(err, op)
	}
	return nil
}
