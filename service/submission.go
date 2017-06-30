package service

import (
  "github.com/jinzhu/gorm"
  "Barracks/data"
  "time"
)

func SelectNotCheckedSubmissions(db *gorm.DB, contest *data.Contest, lastId uint) (submissions []data.Submission){
  // find first pending submission
  var firstPendingSubmission data.Submission
  var lastSubmission data.Submission

  db.Last(&lastSubmission)
  db.Where("ContestId = ? AND result = ? AND id <= ?", contest.ID, 0, lastSubmission.ID).First(&firstPendingSubmission)

  var targetId uint

  if firstPendingSubmission.ID != 0 {
    targetId = firstPendingSubmission.ID-1
  } else {
    targetId = lastSubmission.ID
  }

  if lastId < targetId {
    db.Where("ContestId = ? AND id > ? AND id <= ?", contest.ID, lastId, targetId).Find(&submissions)
  }

  freezeAt := SelectContestFreezeById(db, contest.ID)

  if !freezeAt.IsZero() {
    return submissions
  }

  var filteredSubmissions []data.Submission

  for _, submission := range submissions {
    if !submission.CreatedAt.After(freezeAt) {
      filteredSubmissions = append(filteredSubmissions, submission)
    }
  }

  return
}