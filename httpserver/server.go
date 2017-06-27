package httpserver

import (
  "github.com/labstack/echo"
  "strconv"
  "net/http"
  "fmt"
  "Barracks/rank"
  "Barracks/data"
)

type ProblemStatusElem struct {
  ProblemId uint `json:"problemId"`
  Accepted  bool `json:"accepted"`
}

type StandingRow struct {
  UserId uint `json:"userId"`
  AcceptedCnt uint `json:"acceptedCnt"`
  Rank uint `json:"rank"`
  ProblemStatus []ProblemStatusElem `json:"problemStatus"`
}

func StartServer(contest *data.Contest, port uint) {
  e := echo.New()

  e.GET("/api/acceptedCnts/:userId", func(ctx echo.Context) error {
    userId, err := strconv.Atoi(ctx.Param("userId"))

    if err != nil {
      return ctx.NoContent(http.StatusNotFound)
    }

    val, ok := rank.MyRankData.UserMap[uint(userId)]
    if !ok {
      return ctx.NoContent(http.StatusNotFound)
    }

    userRowRef := &rank.MyRankData.UserRows[rank.MyRankData.UserMap[val]]
    if userRowRef == nil {
      return ctx.NoContent(http.StatusNotFound)
    }

    r := StandingRow{
      AcceptedCnt: userRowRef.AcceptedCnt,
      Rank:        userRowRef.Rank,
    }

    return ctx.JSON(http.StatusOK, r)
  })

  e.GET("/api/problemStatuses/:userId", func(ctx echo.Context) error {
    userId, err := strconv.Atoi(ctx.Param("userId"))
    if err != nil {
      return ctx.NoContent(http.StatusNotFound)
    }

    val, ok := rank.MyRankData.UserMap[uint(userId)]
    if !ok {
      return ctx.NoContent(http.StatusNotFound)
    }

    userRowRef := &rank.MyRankData.UserRows[rank.MyRankData.UserMap[val]]
    if userRowRef == nil {
      return ctx.NoContent(http.StatusNotFound)
    }

    r := StandingRow{}
    for key, val := range rank.MyRankData.ProblemMap {
      r.ProblemStatus = append(r.ProblemStatus,
        ProblemStatusElem{key, userRowRef.ProblemStatuses[val].Accepted})
    }

    return ctx.JSON(http.StatusOK, r)
  })

  e.POST(fmt.Sprintf("/api/%s/submissions/checked", contest.Name), func(ctx echo.Context) error {
    var r []StandingRow
    if err := ctx.Bind(r); err != nil {
      return err
    }

    fmt.Println("POST")
    for idx, val := range r {
      fmt.Println(idx)
      fmt.Println(val.Rank)
      fmt.Println(val.AcceptedCnt)
      fmt.Println(val.UserId)
    }

    return ctx.NoContent(http.StatusOK)
  })

  e.Logger.Fatal(e.Start(":"+strconv.Itoa(int(port))))
}
