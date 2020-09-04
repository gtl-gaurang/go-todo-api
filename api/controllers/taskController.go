package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"todo-api/api/auth"
	"todo-api/api/models"
	"todo-api/api/utils/formaterror"

	"github.com/gin-gonic/gin"
)

// AddTask ... New task add
func AddTask(c *gin.Context) {
	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	fmt.Println(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	uid, err := auth.ExtractTokenID(c.Request)

	task := models.Task{}
	task.UserID = uid
	err = json.Unmarshal(body, &task)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	task.Prepare()
	errorMessages := task.Validate("add")
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	taskCreated, err := task.AddTask(app.DB)
	fmt.Println("err", err)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		errList = formattedError
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"data":   taskCreated,
	})
}

// GetAllTask ... get the list of all task
func GetAllTask(c *gin.Context) {
	task := models.Task{}

	tasks, err := task.GetAllTask(s.DB)
	if err != nil {
		errList["No_task"] = "No task Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   tasks,
	})
}

// UpdateTask ... Update Task by Id
func UpdateTask(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	taskID := c.Param("id")

	// Check if the user id is valid
	tid, err := strconv.ParseInt(taskID, 10, 32)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	// Start processing the request
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	task := models.Task{}
	err = json.Unmarshal(body, &task)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	// Check for previous details
	formerTask := models.Task{}
	err = s.DB.Debug().Model(models.Task{}).Where("id = ?", tid).Take(&formerTask).Error
	if err != nil {
		errList["Task_invalid"] = "The task is does not exist"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	task.Prepare()
	errorMessages := task.Validate("")
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	updatedTask, err := task.UpdateTask(s.DB, uint32(tid))
	if err != nil {
		errList := formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   updatedTask,
	})

}

//DeleteTask ... Delete Task by Id
func DeleteTask(c *gin.Context) {

	taskID := c.Param("id")
	// Is a valid post id given to us?
	tid, err := strconv.ParseInt(taskID, 10, 32)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	// Check if the task exist
	task := models.Task{}
	err = s.DB.Debug().Model(models.Task{}).Where("id = ?", tid).Take(&task).Error
	if err != nil {
		errList["No_task"] = "No task Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}

	// If all the conditions are met, delete the post
	_, err = task.DeleteTask(s.DB, uint32(tid))
	if err != nil {
		errList["Other_error"] = "Please try again later"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Task deleted",
	})
}
