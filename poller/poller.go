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

func StartPoll(
	db *gorm.DB,
	rankInfo *rank.RankInfo,
	rankInfoFreeze *rank.RankInfo,
	contest *data.Contest,
	doneChan *chan bool,
	pushHost *string,
) {

	lastId := uint(0)
	var submissions []data.Submission

	// pushHost와 contest로 결과를 전달할 url 생성
	pushUrl := *pushHost + "/api/" + contest.Name + "/submissions/checked"
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	go func() {
		for {
			// lastId를 통해서 체크되지 않은 모든 제출들을 모아옴
			submissions = service.SelectNotCheckedSubmissions(db, contest, lastId)

			submissionsLen := len(submissions)

			// set last not pending submissions
			if submissionsLen > 0 {
				lastId = submissions[submissionsLen-1].ID
				rankInfo.AddSubmissions(submissions)

				{
					freezeAt := service.SelectContestFreezeById(db, contest.ID)

					if !freezeAt.IsZero() {
						var filteredSubmissions []data.Submission

						for _, submission := range submissions {
							if submission.CreatedAt.Before(freezeAt) {
								filteredSubmissions = append(filteredSubmissions, submission)
							}
						}

						rankInfoFreeze.AddSubmissions(filteredSubmissions)
					} else {
						rankInfoFreeze.AddSubmissions(submissions)
					}
				}

				changes := make(map[uint]bool)
				var ret struct {
					Results []rank.UserRankSummary `json:"results"`
				}
				for _, sub := range submissions {
					if _, ok := rankInfo.RankData.UserMap[sub.UserID]; ok {
						if _, present := changes[sub.UserID]; !present {
							changes[sub.UserID] = true

							sum := rankInfoFreeze.GetUserSummary(sub.UserID, sub.ID)

							realRankInfo := rankInfo.GetUserSummary(sub.UserID, 0)
							sum.AcceptedCnt = realRankInfo.AcceptedCnt
							sum.TotalScore = realRankInfo.TotalScore

							ret.Results = append(ret.Results, *sum)
						}
					}
				}
				if len(changes) > 0 {
					jsonValue, err := json.Marshal(ret)
					if err != nil {
						panic(err)
					}
					_, err = client.Post(pushUrl, "application/json", bytes.NewReader(jsonValue))
					if err != nil {
						fmt.Println(err)
					}
				}
			}

			time.Sleep(time.Second * 5)
		}

	}()

}
