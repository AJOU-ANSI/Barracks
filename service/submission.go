package service

import (
	"github.com/jinzhu/gorm"
	"Barracks/data"
	"fmt"
)

func SelectNotCheckedSubmissions(db *gorm.DB, contest *data.Contest, lastId uint) []data.Submission {
	// 첫 pending 중인 제출
	var firstPendingSubmission data.Submission
	// 마지막 제출
	var lastSubmission data.Submission

	// 마지막 제출을 찾는다
	db.Last(&lastSubmission)
	// pending 중인 제출 중에서 첫 시작을 찾는다
	db.Where("ContestId = ? AND result < 4 AND id <= ?", contest.ID, lastSubmission.ID).First(&firstPendingSubmission)

	fmt.Println(firstPendingSubmission.ID)
	fmt.Println(lastSubmission.ID)

	var targetId uint

	// 만약 pending 중인 제출이 있다면 그 전까지 쭉
	if firstPendingSubmission.ID != 0 {
		targetId = firstPendingSubmission.ID - 1
	} else { // pending 중인 제출이 없다면 끝까지 쭉
		targetId = lastSubmission.ID
	}

	// lastId 보다는 크면서 targetId 보다는 작은 제출들
	// 즉 결과가 나온 확인되지 않은 제출들을 쭉 긁어 온다
	var submissions []data.Submission
	if lastId < targetId {
		db.Where("ContestId = ? AND id > ? AND id <= ?", contest.ID, lastId, targetId).Find(&submissions)
	}

	return submissions

	//return submissions

	//
	//freezeAt := SelectContestFreezeById(db, contest.ID)
	//
	//if freezeAt.IsZero() {
	//  return submissions
	//}
	//
	//var filteredSubmissions []data.Submission
	//
	//for _, submission := range submissions {
	//  if submission.CreatedAt.Before(freezeAt) {
	//    filteredSubmissions = append(filteredSubmissions, submission)
	//  }
	//}
	//
	//return filteredSubmissions
}
