package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

func main() {
	r := gin.Default()
	r.GET("/repositories/:username", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":"pong",
			"username":c.Param("username"),
			"public_repositories": len (getRepositoriesFromUsername(c.Param("username"))),
		})
	})
	r.Run("127.0.0.1:7000")
}




func getRepositoriesFromUsername(username string) []*github.Repository{

	client := github.NewClient(nil)
	ctx := context.Background()

	// get the list of organizations for a specific user.

	orgs, _, err := client.Repositories.List(ctx, username, nil)
	if err != nil {
		logrus.Fatalln(err)
	}

	return orgs
}
