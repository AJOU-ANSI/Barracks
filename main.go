package main

import (
	"Barracks/data"
	"Barracks/httpserver"
	"Barracks/poller"
	"Barracks/rank"
	"Barracks/service"
	"bufio"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
	"strings"
	"time"
	"flag"
)

var db *gorm.DB

func init() {

}

func askContestName() (contestName string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter contest name: ")
	contestName, _ = reader.ReadString('\n')
	contestName = strings.Trim(contestName, "\n ")

	fmt.Printf("[%s]에 대한 랭킹 계산을 시작합니다.\n", contestName)

	return
}

func main() {
	contestNamePtr := flag.String("contest", "", "contest name")
	portPtr := flag.Int("port", 8080, "port number")
	flag.Parse()


	if *contestNamePtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	var err error

	db, err = gorm.Open("mysql", data.DbConfig)

	if err != nil {
		panic(err)
	}

	db.LogMode(true)

	defer db.Close()

	contestName := *contestNamePtr

	contest := service.SelectContestByName(db, &contestName)
	users := service.SelectNormalUsersByContest(db, &contest)
	problems := service.SelectProblemsByContest(db, &contest)

	tickDuration := 5 * time.Second
	doneChan := make(chan bool)

	rank.InitData(&contest, &users, &problems)
	poller.StartPoll(db, &tickDuration, &contest, &doneChan)
	httpserver.StartServer(&contest, uint(*portPtr))
}
