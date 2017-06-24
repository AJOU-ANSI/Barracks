package rank

import (
  "time"
  "Barracks/data"
  "container/heap"
)

func newUserRow (user *data.User, problems *[]data.Problem) (userRow userRow) {
  userRow.Rank = 1
  userRow.StrId = (*user).StrId
  userRow.ProblemStatuses = make([]problemStatus, len(*problems))

  return
}

func newRankData (contest *data.Contest, users *[]data.User, problems *[]data.Problem) (r *rankData) {

  r = &rankData{
    CalcAt: time.Now(),
    ContestInfo: contest,
    UserRows: make([]userRow, len(*users)),
    UserMap: make(map[uint]uint),
    ProblemMap: make(map[uint]uint),
  }

  for index, problem := range *problems {
    r.ProblemMap[problem.ID] = uint(index)
  }

  for index, user := range *users {
    r.UserMap[user.ID] = uint(index)
    r.UserRows[index] = newUserRow(&user, problems)
  }

  return
}

var MyRankData *rankData
var RankHeap *rankHeap

func InitData (contest *data.Contest, users *[]data.User, problems *[]data.Problem) {
  MyRankData = newRankData(contest, users, problems)

  RankHeap = &rankHeap{}
  heap.Init(RankHeap)
}

func calcRanks() {
  for index, userRow := range MyRankData.UserRows {
    heap.Push(RankHeap, rankNode{Penalty: userRow.Penalty, AcceptedCnt: userRow.AcceptedCnt, UserIndex: uint(index)})
  }

  rankValue := uint(1)
  var beforeRankNode *rankNode

  for RankHeap.Len() > 0 {
    popRankNode := rankNode(heap.Pop(RankHeap).(rankNode))

    if beforeRankNode != nil && (beforeRankNode.AcceptedCnt != popRankNode.AcceptedCnt || beforeRankNode.Penalty != popRankNode.Penalty) {
      rankValue++
    }

    MyRankData.UserRows[popRankNode.UserIndex].Rank = rankValue
    beforeRankNode = &popRankNode
  }
}

func analyzeSubmissions(submissions *[]data.Submission) {
  contestInfo := MyRankData.ContestInfo

  // 각각의 제출에 대하여
  for _, submission := range *submissions {
    // userRow와 problemStatus를 구한다.

    if _, ok := MyRankData.UserMap[submission.UserID]; !ok {
      continue
    }

    userRow := &MyRankData.UserRows[MyRankData.UserMap[submission.UserID]]
    problemStatus := &userRow.ProblemStatuses[MyRankData.ProblemMap[submission.ProblemID]]

    // 만약 제출이 정답 소스코드라면
    if data.IsAccepted(submission.Result) {

      // 문제가 맞지 않은 상황이라면
      if !problemStatus.Accepted {
        penalty := submission.CreatedAt.Sub(contestInfo.Start) + time.Duration(problemStatus.WrongCount) * 20 * time.Minute

        (*userRow).Penalty += penalty
        (*problemStatus).Accepted = true
        (*userRow).AcceptedCnt++
      }

    } else { // 제출이 틀린 소스코드라면
      (*problemStatus).WrongCount++
    }
  }
}

func AddSubmissions (submissions *[]data.Submission) {
  analyzeSubmissions(submissions)
  calcRanks()
}

