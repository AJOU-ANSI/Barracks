package poller

import (
  "time"
  "Barracks/rank"
  "encoding/json"
  "fmt"
  "net/http"
  "Barracks/service"
  "Barracks/data"
  "github.com/jinzhu/gorm"
  "Barracks/httpserver"
  "bytes"
)


func StartPoll(db *gorm.DB, tickDuraion *time.Duration, contest *data.Contest, doneChan *chan bool) {
  tickerChan := time.NewTicker(*tickDuraion).C

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
          var changes map[uint]*httpserver.StandingRow
          var ret []*httpserver.StandingRow

          for _, sub := range submissions {
            if _, present := changes[sub.UserID]; !present {
              userRowRef := &rank.MyRankData.UserRows[rank.MyRankData.UserMap[sub.UserID]]
              r := &httpserver.StandingRow{
                UserId:      sub.UserID,
                AcceptedCnt: userRowRef.AcceptedCnt,
                Rank:        userRowRef.Rank,
              }
              for key, val := range rank.MyRankData.ProblemMap {
                r.ProblemStatus = append(r.ProblemStatus,
                  httpserver.ProblemStatusElem{key, userRowRef.ProblemStatuses[val].Accepted})
              }
              changes[sub.UserID] = r
              ret = append(ret, r)
            }
            jsonValue, err := json.Marshal(ret)
            if err != nil {
              panic (err)
            }
            fmt.Println(jsonValue)
            http.Post(
              "/api/" + contest.Name + "/submissions/checked",
              "application/json",
              bytes.NewReader(jsonValue))

          }
        }
      case <- *doneChan:
        *doneChan <- true
        break loop
      }
    }
  }()
}