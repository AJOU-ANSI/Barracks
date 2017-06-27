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
    "fmt"
)


func StartPoll(db *gorm.DB, rankInfo *rank.RankInfo, tickDuraion *time.Duration, contest *data.Contest, doneChan *chan bool) {
    tickerChan := time.NewTicker(*tickDuraion).C

    lastId := uint(0)
    var submissions []data.Submission
    go func() {
        pushUrl := "http://127.0.0.1/api/" + contest.Name + "/submissions/checked"
        client := &http.Client{
            Timeout: time.Second * 10,
        }
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
                        _, err = client.Post(pushUrl, "application/json",bytes.NewReader(jsonValue))
                        if err != nil {
                            fmt.Println(err)
                        }
                    }
                }
            case <-*doneChan:
                *doneChan <- true
                break loop
            }
        }
    }()
}
