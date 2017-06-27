package poller

import (
    "time"
    "Barracks/rank"
    "encoding/json"
    "net/http"
    "Barracks/service"
    "Barracks/data"
    "github.com/jinzhu/gorm"
    "bytes"
)


func StartPoll(db *gorm.DB, rankInfo *rank.RankInfo, tickDuraion *time.Duration, contest *data.Contest, doneChan *chan bool) {
    tickerChan := time.NewTicker(*tickDuraion).C

    lastId := uint(0)
    var submissions []data.Submission
    go func() {
        loop:
        for {
            select {
            case <-tickerChan:
                submissions = service.SelectNotCheckedSubmissions(db, contest, lastId)

                submissionsLen := len(submissions)

                // set last not pending submissions
                if submissionsLen > 0 {
                    lastId = submissions[submissionsLen-1].ID
                    rankInfo.AddSubmissions(&submissions)

                    changes := make(map[uint]bool)
                    var ret struct {
                        Results []rank.UserRankSummary `json:"results"`
                    }
                    for _, sub := range submissions {
                        if _, ok := rankInfo.RankData.UserMap[sub.UserID]; ok {
                            if _, present := changes[sub.UserID]; !present {
                                changes[sub.UserID] = true
                                sum := rankInfo.GetUserSummary(sub.UserID)
                                ret.Results = append(ret.Results, *sum)
                            }
                        }
                    }
                    if len(changes) > 0 {
                        jsonValue, err := json.Marshal(ret)
                        if err != nil {
                            panic(err)
                        }
                        http.Post("http://127.0.0.1:8080/api/"+contest.Name+"/submissions/checked", "application/json",
                        bytes.NewReader(jsonValue))
                    }
                }
            case <-*doneChan:
                *doneChan <- true
                break loop
            }
        }
    }()
}
