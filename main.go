package main

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
  "fmt"
  "Barracks/service"
  "bufio"
  "os"
  "strings"
  "Barracks/rank"
  "time"
  "Barracks/data"
  "net/http"
  "github.com/labstack/echo"
  "strconv"
)

var db *gorm.DB

func init() {

}

func askContestName() (contestName string) {
  reader := bufio.NewReader(os.Stdin)
  fmt.Print("Enter contest name: ")
  contestName, _ = reader.ReadString('\n')
  contestName = strings.Trim(contestName, "\n ")

  fmt.Printf("[%s]에 대한 랭킹 계산을 시작합니다.\n", contestName)

  return
}

func main() {
  var err error

  db, err = gorm.Open("mysql", data.DbConfig)

  if err != nil {
    panic(err)
  }

  db.LogMode(true)

  defer db.Close()

  contestName := "shake16open" //askContestName()
  tickDuration := 5 * time.Second

  contest := service.SelectContestByName(db, &contestName)
  users := service.SelectNormalUsersByContest(db, &contest)
  problems := service.SelectProblemsByContest(db, &contest)

  rank.InitData(&contest, &users, &problems)

  tickerChan := time.NewTicker(tickDuration).C
  doneChan := make(chan bool)

  lastId := uint(0)
  var submissions []data.Submission
  go func() {
  loop:
    for {
      select {
      case <- tickerChan:
        submissions = service.SelectNotCheckedSubmissions(db, contest, lastId)

        submissionsLen := len(submissions)

        // set last not pending submissions
        if submissionsLen > 0 {
          lastId = submissions[submissionsLen-1].ID
          rank.AddSubmissions(&submissions)
        }
      case <- doneChan:
        doneChan <- true
        break loop
      }
    }
  }()

  e := echo.New()
  e.GET("/api/acceptedCnts/:userId", func(ctx echo.Context) error {
    userId, err := strconv.Atoi(ctx.Param("userId"))
    if err != nil {
      return ctx.NoContent(http.StatusNotFound)
    }
    userRowRef := &rank.MyRankData.UserRows[rank.MyRankData.UserMap[uint(userId)]]
    if userRowRef == nil {
      return ctx.NoContent(http.StatusNotFound)
    }
    r := &struct {
      AcceptedCnt uint `json:"acceptedCnt"`
      Rank uint `json:"rank"`
    }{
      userRowRef.AcceptedCnt,
      userRowRef.Rank,
    }
    return ctx.JSON(http.StatusOK, r)
  })
  e.GET("/api/problemStatuses/:userId", func(ctx echo.Context) error {
    userId, err := strconv.Atoi(ctx.Param("userId"))
    if err != nil {
      return ctx.NoContent(http.StatusNotFound)
    }
    userRowRef := &rank.MyRankData.UserRows[rank.MyRankData.UserMap[uint(userId)]]
    if userRowRef == nil {
      return ctx.NoContent(http.StatusNotFound)
    }
    type problemStatusElem struct {
      ProblemId uint
      Accepted  bool
    }
    r := &struct {
      ProblemStatus []problemStatusElem
    }{}
    for key, val := range rank.MyRankData.ProblemMap {
      r.ProblemStatus = append(r.ProblemStatus, problemStatusElem{key, userRowRef.ProblemStatuses[val].Accepted})
    }
    return ctx.JSON(http.StatusOK, r)
  })
  e.Logger.Fatal(e.Start(":8080"))
}
