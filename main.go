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

          go func() {
            rank.AddSubmissions(&submissions)

            for index, userRow := range rank.MyRankData.UserRows {
              fmt.Println(index, userRow.Rank, userRow.AcceptedCnt, uint(userRow.Penalty))
            }
            fmt.Println("----------------------")
          }()

        }
      case <- doneChan:
        doneChan <- true
        break loop
      }
    }
  }()

  /*
  GET MyRankData 요청

   */
}

/*
  - gorm (db orm)
  - ginkgo, gomego (testing tool)
  - echo (웹프레임워크
 */