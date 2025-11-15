package db

import (
	"context"

	"github.com/Vy4cheSlave/qna/internal/domain"
	"github.com/Vy4cheSlave/qna/internal/infrastructure/db/dto"

	"github.com/pkg/errors"
)

var (
	ErrNotFound = errors.New("not found")
)

func (r *Repository) CreateUser(ctx context.Context, userName *string) (userId *string, err error) {
	const op = "internal/infrastructure/db/repository.Repository.CreateUser"

	newUser := dto.User{
		Name: *userName,
	}

	result := r.db.WithContext(ctx).Create(&newUser)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, op)
	}

	if result.RowsAffected == 0 {
		return nil, errors.Wrap(ErrNotFound, op)
	}

	return &newUser.Id, nil
}

func (r *Repository) ReadUsers(ctx context.Context) (*[]domain.User, error) {
	const op = "internal/infrastructure/db/repository.Repository.ReadUsers"

	var usersDb []dto.User

	result := r.db.WithContext(ctx).Find(&usersDb)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, op)
	}

	users := make([]domain.User, 0, len(usersDb))
	for _, user := range usersDb {
		users = append(users, domain.User{
			Id:   user.Id,
			Name: user.Name,
		})
	}

	return &users, nil
}

func (r *Repository) DeleteUser(ctx context.Context, userId *string) error {
	const op = "internal/infrastructure/db/repository.Repository.DeleteUser"

	result := r.db.WithContext(ctx).Delete(&domain.User{}, "id = ?", *userId)

	if result.Error != nil {
		return errors.Wrap(result.Error, op)
	}

	if result.RowsAffected == 0 {
		return errors.Wrap(ErrNotFound, op)
	}

	return nil
}

func (r *Repository) ReadQuestions(ctx context.Context) (*[]domain.Question, error) {
	const op = "internal/infrastructure/db/repository.Repository.ReadQuestion"

	var questionsDb []dto.Question

	result := r.db.WithContext(ctx).Find(&questionsDb)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, op)
	}

	questions := make([]domain.Question, 0, len(questionsDb))
	for _, q := range questionsDb {
		questions = append(questions, domain.Question{
			Id:   q.Id,
			Text: q.Text,
		})
	}

	return &questions, nil
}

func (r *Repository) CreateQuestion(ctx context.Context, question *string) (questionId int, err error) {
	const op = "internal/infrastructure/db/repository.Repository.CreateQuestion"

	newQuestion := dto.Question{
		Text: *question,
	}

	result := r.db.WithContext(ctx).Create(&newQuestion)

	if result.Error != nil {
		return 0, errors.Wrap(result.Error, op)
	}

	if result.RowsAffected == 0 {
		return 0, errors.Wrap(ErrNotFound, op)
	}

	return newQuestion.Id, nil
}

func (r *Repository) ReadQuestionAndAnswers(ctx context.Context, questionId int) (*domain.Question, *[]domain.Answer, error) {
	const op = "internal/infrastructure/db/repository.Repository.ReadQuestionAndAnswers"

	var questionDb dto.Question
	var answersDb []dto.Answer

	result := r.db.WithContext(ctx).First(&questionDb, questionId)
	if result.Error != nil {
		return nil, nil, errors.Wrap(result.Error, op)
	}

	result = r.db.WithContext(ctx).Where("question_id = ?", questionId).Find(&answersDb)
	if result.Error != nil {
		return nil, nil, errors.Wrap(result.Error, op)
	}

	question := domain.Question{
		Id:   questionDb.Id,
		Text: questionDb.Text,
	}

	answers := make([]domain.Answer, 0, len(answersDb))
	for _, a := range answersDb {
		answers = append(answers, domain.Answer{
			Id:         a.Id,
			QuestionId: a.QuestionId,
			UserId:     a.UserId,
			Text:       a.Text,
		})
	}

	return &question, &answers, nil
}

func (r *Repository) DeleteQuestionAndAnswers(ctx context.Context, questionId int) error {
	const op = "internal/infrastructure/db/repository.Repository.ReadQuestionAndAnswers"

	result := r.db.WithContext(ctx).Delete(&dto.Question{}, questionId)

	if result.Error != nil {
		return errors.Wrap(result.Error, op)
	}

	if result.RowsAffected == 0 {
		return errors.Wrap(ErrNotFound, op)
	}

	return nil
}

func (r *Repository) CreateAnswerToQuestion(ctx context.Context, answer *domain.Answer) (answerId int, err error) {
	const op = "internal/infrastructure/db/repository.Repository.CreateAnswerToQuestion"

	answerDb := dto.Answer{
		UserId:     answer.UserId,
		QuestionId: answer.QuestionId,
		Text:       answer.Text,
	}

	result := r.db.WithContext(ctx).Create(&answerDb)

	if result.Error != nil {
		return 0, errors.Wrap(result.Error, op)
	}

	if result.RowsAffected == 0 {
		return 0, errors.Wrap(ErrNotFound, op)
	}

	return answerDb.Id, nil
}

func (r *Repository) ReadAnswer(ctx context.Context, answerId int) (*domain.Answer, error) {
	const op = "internal/infrastructure/db/repository.Repository.ReadAnswer"

	var answerDb dto.Answer

	result := r.db.WithContext(ctx).First(&answerDb, answerId)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, op)
	}

	answer := domain.Answer{
		Id:         answerDb.Id,
		QuestionId: answerDb.QuestionId,
		UserId:     answerDb.UserId,
		Text:       answerDb.Text,
	}

	return &answer, nil
}

func (r *Repository) DeleteAnswer(ctx context.Context, answerId int) error {
	const op = "internal/infrastructure/db/repository.Repository.DeleteAnswer"

	result := r.db.WithContext(ctx).Delete(&dto.Answer{}, answerId)

	if result.Error != nil {
		return errors.Wrap(result.Error, op)
	}

	if result.RowsAffected == 0 {
		return errors.Wrap(ErrNotFound, op)
	}

	return nil
}
