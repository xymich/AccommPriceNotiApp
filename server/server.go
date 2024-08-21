package server

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoclient *mongo.Client
)

func init() {
	server = gin.Default()
	ctx = context.TODO()
}

func Start() {
	mongoclient = MongoConnect()
}