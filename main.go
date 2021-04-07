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
	book, _, _ := client.LastBookService.GetLastBook()
	bookId := eudic.GetTheFirstNumberFromString(book.Meta.Bookid)
	_, _ = client.SyncReciteService.SyncRecite(bookId, book.BookName)
	reciteStarted, _ := client.StartReciteService.StartRecite(bookId, book.BookName)
	// 标记背诵任务已经完成
	done = reciteStarted.TaskFinished

	if !done {
		fmt.Println("===Start===")
		fmt.Println("输入熟悉程度: 0 不认识  2 模糊  5 熟悉")
		_, _ = recite(client, bookId, book.BookName, reciteStarted)
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
	// 打印释义
	fmt.Printf("%s\n----- ", preCard.Card.Answer)
	progress := preCard.TodayProgress
	fmt.Printf("%d / %d -----\n\n", progress.TodayFinishedCount,
		progress.TodayFinishedCount+progress.PendingDueCardCount+progress.PendingNewCardCount)
	nextCard, _ = client.AnswerCardService.AnswerCard(bookId, bookName, preCard.Card.CardID, inputEase)
	if !nextCard.TaskFinished {
		return recite(client, bookId, bookName, nextCard)
	} else {
		fmt.Println("Already done...")
		done = true
		return nil, err
	}
}
