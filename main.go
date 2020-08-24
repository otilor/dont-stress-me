package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

var ctx = context.Background()


func main() {
	r := gin.Default()

	type Details struct {
		Username string
		Public_repositories int
	}

	r.GET("/repositories/:username", func(c *gin.Context) {
		c.JSON(200, &Details{Username: c.Param("username"), Public_repositories: len(getRepositoriesFromUsername(c.Param("username")))})
	})
	r.Run("127.0.0.1:7000")
}

func ExampleClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		Password: "",
		DB: 0,
	})
	err := rdb.Set("key", "value", 10000)
	if err != nil {
		logrus.Fatalln(err)
	}

	fmt.Println("Successfully set key!")
}

func cacheDetails(details string) {

}

func getRepositoriesFromUsername(username string) []*github.Repository{

	client := github.NewClient(nil)


	// get the list of organizations for a specific user.

	opt := &github.RepositoryListOptions{Type: "owner", Sort: "updated", Direction: "desc"}
	var allRepos []*github.Repository

	for {
		repos, resp, err := client.Repositories.List(ctx, username, opt)
		if err != nil {
			logrus.Fatalln(err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allRepos
}
