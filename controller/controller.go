package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mySchool/config"
	"github.com/mySchool/database"
	"github.com/mySchool/model"
	"github.com/mySchool/password"
	"github.com/mySchool/token"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GenCodeID(code int) string {
	codePrefix := "LAW"
	return fmt.Sprintf("%v%d", codePrefix, code)
}

func Register(c *gin.Context) {
	var student model.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		fmt.Println("An error occured while binding json data", err)
	}
	validate := validator.New()
	if err := validate.Struct(student); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("An error occured. Error(%v)", err),
		})
		c.Abort()
		return
	}
	student.Code = GetMaxCode() + 1
	student.CodeID = GenCodeID(student.Code)
	student.Token, student.RefreshToken, _ = token.GenerateToken(student.Name, student.Email, student.CodeID)
	student.RegisteredAt = time.Now()
	student.UpdatedAt = time.Now()
	student.Password = password.HashPassword(student.Password)
	student.IpAddress = GetStudentIpAddress()
	if err := config.ValidateEmail(student.Email); err != "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		c.Abort()
		return
	}
	message, err := config.InsertStudentIntoDB(student)
	if err != "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err,
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": message,
		"data":    student,
	})
}

func GetMaxCode() int {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	filter := bson.M{}
	findOptions := options.Find().SetSort(bson.M{"code": -1}).SetLimit(1)
	cursor, _ := database.Collection.Find(ctx, filter, findOptions)
	defer cursor.Close(ctx)
	var students []model.Student
	for cursor.Next(ctx) {
		var student model.Student
		cursor.Decode(&student)
		students = append(students, student)
	}
	var maxCode int
	for _, student := range students {
		maxCode = student.Code
	}
	return maxCode
}

func GetStudent(c *gin.Context) {
	codeId := c.Param("code_id")
	student, err := config.GetSingleStudentFromDB(codeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Student with Code ID(%v) is not found.", codeId),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, student)
}

func GetAllStudent(c *gin.Context) {
	students := config.GetAllStudentsFromDB()
	c.JSON(http.StatusOK, students)
}

func DeleteStudent(c *gin.Context) {
	codeId := c.Param("code_id")
	_, err := config.DeleteStudentFromDB(codeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("An error occured while deleteing student with Code ID(%v). Error(%v)", codeId, err),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Students with code id(%v) has been deleted succesfully.", codeId),
	})
}

func UpdateStudent(c *gin.Context) {
	codeID := c.Param("code_id")
	var student model.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}
	var validate = validator.New()
	if err := validate.Struct(student); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		c.Abort()
		return
	}
	_, err := config.GetSingleStudentFromDB(codeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Student with codeID(%v) is not found.", codeID),
		})
		c.Abort()
		return
	}
	result, updateErr := config.UpdateStudentFromDB(student, codeID)
	if updateErr != "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   fmt.Sprintf("An error occured while trying to update student data. Error(%v)", updateErr),
			"message": result,
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":        "Update Successful",
		"student":        student,
		"Modified Count": result,
	})
}

func DeleteAllStudents(c *gin.Context) {
	no, err := config.DeleteAllFromDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"No of students deleted": no,
	})
}

func Logout(c *gin.Context) {
	signedToken := c.GetHeader("token")

	if err := token.SetTokenToExpired(signedToken); err != "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful. Token expired",
	})
}

func GetStudentIpAddress() string {
	type IpAddress struct {
		IP string `json:"ip"`
	}
	var userip IpAddress
	url := "https://api64.ipify.org?format=json"
	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := client.Do(req)
	respByte, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(respByte, &userip)
	return userip.IP
}
