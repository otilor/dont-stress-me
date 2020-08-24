package main

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

var ctx = context.Background()


type Details struct {
	Public_repositories int
	Username string
}

var rdb = redis.NewClient(&redis.Options{
Addr: "127.0.0.1:6379",
Password: "",
DB: 0,
})

func main() {
	r := gin.Default()



	r.GET("/repositories/:username", func(c *gin.Context) {
		responseDetails := &Details{Username: c.Param("username"), Public_repositories: len(getRepositoriesFromUsername(c.Param("username")))}
		c.JSON(200, responseDetails)
		if cacheDetails(responseDetails)  {
			logrus.Println("Successfully cached!")
		} else {
			logrus.Println("Something else happened!")
		}
	})
	r.Run("127.0.0.1:7000")
}



func cacheDetails(details *Details) bool {
	username := details.Username

	marshalledDetails, err := json.Marshal(details)
	if err != nil {
		logrus.Fatalln(err)
	}

	err = rdb.Set(string(username), marshalledDetails, 0).Err()
	if err != nil {
		logrus.Fatalln(err)
	}
	return true
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
