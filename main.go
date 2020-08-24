package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"net/http"
)

var ctx = context.Background()

var client = github.NewClient(nil)

type Details struct {
	PublicRepositories int
	Username           string
}

// Start redis client.
var rdb = redis.NewClient(&redis.Options{
Addr: "127.0.0.1:6379",
Password: "",
DB: 0,
})

func main() {
	r := gin.Default()
	r.Handle("GET", "greet", greetMe)
	r.Handle("GET","repositories/:username", repositoryHandler)

	// Start the server
	r.Run("127.0.0.1:7000")
}

func repositoryHandler(c *gin.Context) {
	username := c.Param("username")

	if isCached(username) {
		cachedDetails,err  := getUserDetailsIfCached(username)
		if err != nil {
			logrus.Fatalln(err)
		}

		c.JSON(http.StatusOK, cachedDetails)
	} else {
		logrus.Println("result is not cached!")
		responseDetails := &Details{
			Username:           username,
			PublicRepositories: len(getRepositoriesFromUsername(username)),
		}

		c.JSON(http.StatusOK, responseDetails)
		cacheDetails(responseDetails)
		logrus.Println("Successfully cached!")
	}
}

func greetMe(c *gin.Context) {
	fmt.Println("Welcome, Soldier")
}

func getUserDetailsIfCached(username string) (interface{}, error) {
	return rdb.Get(username).Result()
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

func isCached(username string) bool {
	err := rdb.Get(username).Err()

	if err != nil {
		return false
	}

	return true
}