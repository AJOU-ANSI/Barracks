package rank

import (
  "time"
  "Barracks/data"
)

type problemStatus struct {
  WrongCount      uint
  Status          string
  Accepted        bool
}

type userRow struct {
  Rank            uint
  StrId           string
  AcceptedCnt     uint
  Penalty         time.Duration
  ProblemStatuses []problemStatus
}

type rankData struct {
  CalcAt      time.Time
  UserRows    []userRow
  ContestInfo *data.Contest
  UserMap     map[uint]uint
  ProblemMap  map[uint]uint
}

type rankNode struct {
  UserIndex     uint
  Penalty       time.Duration
  AcceptedCnt   uint
}

type rankHeap []rankNode

func (h rankHeap) Len() int {
  return len(h)
}

func (h rankHeap) Less(i, j int) bool {
  if h[i].AcceptedCnt == h[j].AcceptedCnt {
    return h[i].Penalty < h[j].Penalty
  }

  return h[i].AcceptedCnt > h[j].AcceptedCnt
}

func (h rankHeap) Swap(i, j int) {
  h[i], h[j] = h[j], h[i]
}

func (h *rankHeap) Push(element interface{}) {
  *h = append(*h, element.(rankNode))
}

func (h *rankHeap) Pop() interface{} {
  old := *h
  n := len(old)
  element := old[n-1]
  *h = old[0 : n-1]
  return element
}