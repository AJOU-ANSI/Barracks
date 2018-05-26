package main

import (
	"Barracks/data"
	"Barracks/httpserver"
	"Barracks/poller"
	"Barracks/rank"
	"Barracks/service"
	"flag"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

func main() {
	// 컨테스트 이름, 포트, 채점 완료시 푸시 서버 입력
	contestNamePtr := flag.String("contest", "", "contest name")
	portPtr := flag.Int("port", 8080, "port number")
	pushHostPtr := flag.String("pushHost", "", "host domain to push submission info")

	flag.Parse()

	// 값이 비어있다면 채워달라고 전달하고 종료
	if *contestNamePtr == "" || *pushHostPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// db 연결
	db := data.InitDB()
	defer db.Close()

	contestName := *contestNamePtr

	// 컨테스트 이름으로 db에서 컨테스트 획득
	contest := service.SelectContestByName(db, &contestName)
	// db에서 컨테스트 참가 인원 획득
	users := service.SelectNormalUsersByContest(db, &contest)
	// db에서 컨테스트 문제 획득
	problems := service.SelectProblemsByContest(db, &contest)

	doneChan := make(chan bool)

	// 콘테스트 정보 초기화
	rankInfo := rank.NewRankInfo(&contest, &users, &problems)
	// 프리징되었을때 콘테스트 정보 초기화
	rankInfoFreezed := rank.NewRankInfo(&contest, &users, &problems)

	// db로부터 제출 항목들 polling 시작
	poller.StartPoll(db, rankInfo, rankInfoFreezed, &contest, &doneChan, pushHostPtr)
	// http 서버 가동
	httpserver.StartServer(rankInfo, rankInfoFreezed, uint(*portPtr))
}
