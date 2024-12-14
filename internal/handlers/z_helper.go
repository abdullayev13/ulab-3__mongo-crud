package handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func fail(c *gin.Context, err error) {

}

func failMsg(c *gin.Context, msg string) {

}

func objIfFromParam(c *gin.Context, name string) (primitive.ObjectID, error) {

	param := c.Param(name)

	return primitive.ObjectIDFromHex(param)
}

type Obj = map[string]any
