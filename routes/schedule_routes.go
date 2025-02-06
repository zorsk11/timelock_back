package routes

import (
	"access-control-system/controllers"

	"github.com/gin-gonic/gin"
)

func ScheduleRoutes(r *gin.Engine) {
	r.POST("/schedule", controllers.CreateSchedule)
	r.GET("/schedule", controllers.GetSchedules)
	r.PUT("/schedules/:id", controllers.UpdateSchedule)
	r.DELETE("/schedules/:id", controllers.DeleteSchedule)

}
