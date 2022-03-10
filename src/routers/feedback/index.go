package feedback

import (
	"github.com/gin-gonic/gin"
	"todoList/src/controllers/feedbackController"
	"todoList/src/controllers/middleware/authorize"
)

type FeedbackRouter struct {}

func (*FeedbackRouter) InitRouter(router gin.IRoutes)  {
	auth := new(authorize.Authorize)
	router = router.Use(auth.Auth)

	controller := new(feedbackController.FeedbackController)
	router.POST("feedback", controller.Create)
}