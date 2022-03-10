package feedbackModel

import "todoList/src/models"

type FeedbackModel struct {
	Content string `json:"content"`
	Imgs    string `json:"imgs"`
	UserId  string `json:"user_id"`
	models.Model
}

func (FeedbackModel) TableName() string {
	return "feed_back"
}
