package config

import (
	"context"
	"fmt"
	"time"

	"github.com/badoux/checkmail"
	"github.com/mySchool/database"
	"github.com/mySchool/model"
	"github.com/mySchool/password"
	"go.mongodb.org/mongo-driver/bson"
)

var collection = database.Collection

func CheckEmailFormat(email string) string {
	if err := checkmail.ValidateFormat(email); err != nil {
		return fmt.Sprintf("Email format is invalid. Error('%v')", err)
	}
	return ""
}

func ValidateEmail(email string) string {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	filter := bson.M{"email": email}
	emailCount, _ := collection.CountDocuments(ctx, filter)
	if emailCount > 0 {
		return "A user with email already exist."
	}
	if err := CheckEmailFormat(email); err != "" {
		return err
	}
	return ""
}

func InsertStudentIntoDB(student model.Student) (string, string) {
	var ctx, cancel = context.WithTimeout(context.TODO(), 10*time.Second)
	result, err := collection.InsertOne(ctx, student)
	defer cancel()
	if err != nil {
		return "", fmt.Sprintf("An error occured while inserting the user into the database. Error(%v)", err)
	}
	return fmt.Sprintf("User has been inserted successfully with ID(%v)", result.InsertedID), ""
}

func GetSingleStudentFromDB(codeID string) (model.Student, error) {
	var ctx, cancel = context.WithTimeout(context.TODO(), 10*time.Second)
	var student model.Student
	defer cancel()
	filter := bson.M{"code_id": codeID}
	err := collection.FindOne(ctx, filter).Decode(&student)
	return student, err
}

func GetAllStudentsFromDB() []model.Student {
	var ctx, cancel = context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	var students []model.Student
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var student model.Student
		cursor.Decode(&student)
		students = append(students, student)
	}
	return students
}

func DeleteStudentFromDB(codeId string) (interface{}, error) {
	var ctx, cancel = context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	result, err := collection.DeleteOne(ctx, bson.M{"code_id": codeId})
	return result.DeletedCount, err
}

func UpdateStudentFromDB(student model.Student, codeID string) (string, string) {
	var ctx, cancel = context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	var updateDetails = bson.M{}
	if student.Name != "" {
		updateDetails["name"] = student.Name
	}
	if student.Age != 0 {
		updateDetails["age"] = student.Age
	}
	if student.Email != "" {
		if err := CheckEmailFormat(student.Email); err != "" {
			return "", err
		}
		updateDetails["email"] = student.Email
	}
	if student.Password != "" {
		updateDetails["password"] = password.HashPassword(student.Password)
	}
	updateTime := time.Now()
	updateDetails["updated_at"] = updateTime
	filter := bson.M{"code_id": codeID}
	updateData := bson.M{"$set": updateDetails}
	result, _ := collection.UpdateOne(ctx, filter, updateData)
	fmt.Println(updateDetails)
	return fmt.Sprintf("Modified Count: %v", result.ModifiedCount), ""
}

func DeleteAllFromDB() (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	result, err := collection.DeleteMany(ctx, bson.M{})
	return result.DeletedCount, err
}
