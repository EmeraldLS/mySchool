package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mySchool/controller"
	"github.com/mySchool/middleware"
)

func RunSever() {
	r := gin.Default()
	api := r.Group("/api")
	api.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Homepage",
		})
	})
	api.POST("/register", controller.Register)
	{
		secured := api.Group("/secured").Use(middleware.Auth())
		{
			secured.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Authenticated",
				})
			})
			secured.GET("/students", controller.GetAllStudent)
			secured.GET("/students/:code_id", controller.GetStudent)
			secured.DELETE("/students/:code_id", controller.DeleteStudent)
			secured.PUT("/students/:code_id", controller.UpdateStudent)
			secured.DELETE("/students", controller.DeleteAllStudents)
		}
	}
	r.Run(":8080")
}
