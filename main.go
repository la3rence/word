package main

import (
	"fmt"
	eudic "github.com/Lonor/go-eudic"
	"os"
)

var done = false

func main() {
	client, err := eudic.NewEudicClientByPassword(
		os.Getenv("EUDIC_USERNAME"),
		os.Getenv("EUDIC_PASSWORD"),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("获取书本...")
	book, _, _ := client.LastBookService.GetLastBook()
	bookId := eudic.GetTheFirstNumberFromString(book.Meta.Bookid)
	// 同步进度
	fmt.Println("获取学习进度...")
	_, _ = client.SyncReciteService.SyncRecite(bookId, book.BookName)
	fmt.Println("获取单词...")
	reciteStarted, _ := client.StartReciteService.StartRecite(bookId, book.BookName)
	// 标记背诵任务已经完成
	done = reciteStarted.TaskFinished

	if !done {
		fmt.Println("===Start===")
		fmt.Println("输入熟悉程度: 0 不认识  2 模糊  5 熟悉")
		_, _ = recite(client, bookId, book.BookName, reciteStarted)
	} else {
		fmt.Println("今日进度已完成.")
	}
}

func recite(client *eudic.EudicClient, bookId, bookName string, preCard *eudic.ReciteResponse) (nextCard *eudic.ReciteResponse, err error) {
	question := preCard.Card.Question
	fmt.Println(question)
	fmt.Print("熟悉程度: ")
	var inputEase int64
	_, err = fmt.Scan(&inputEase)
	if err != nil {
		fmt.Println("请输入数字: 0 / 2 / 5")
		return nil, err
	}
	// 打印释义和当日进度
	progress := preCard.TodayProgress
	fmt.Printf("%s\n----- %d / %d -----\n\n", preCard.Card.Answer, progress.TodayFinishedCount,
		progress.TodayFinishedCount+progress.PendingDueCardCount+progress.PendingNewCardCount)
	nextCard, _ = client.AnswerCardService.AnswerCard(bookId, bookName, preCard.Card.CardID, inputEase)
	if !nextCard.TaskFinished {
		return recite(client, bookId, bookName, nextCard)
	} else {
		fmt.Println("今日进度已完成, 打卡中...")
		done = true
		// 打卡
		checkin, err := client.CheckInService.CheckIn()
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		if checkin.Ischeckin {
			fmt.Println("已自动打卡, 同步学习进度...")
		}
		// 同步最新进度
		syncRecite, err := client.SyncReciteService.SyncRecite(bookId, bookName)
		if syncRecite {
			fmt.Println("已同步最新进度.")
		}
		return nil, err
	}
}
