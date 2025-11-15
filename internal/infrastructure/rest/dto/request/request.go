package request

type CreateUserRequest struct {
	Name string `json:"name"`
}

type CreateQuestionRequest struct {
	Text string `json:"text"`
}

type CreateAnswerToQuestionRequest struct {
	UserId string `json:"user_id"`
	Text   string `json:"text"`
}
