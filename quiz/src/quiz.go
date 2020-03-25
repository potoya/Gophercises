package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	csvFilePath := flag.String("csvPath", "../in/problems.csv", "csv file with questions, answers for the quiz.")
	timeLimit := flag.Int("timeLimit", 30, "time limit to answer the quiz in secs")
	flag.Parse()

	fmt.Println("\n	   Quiz")
	fmt.Println("---------------------------")
	rows := readCsvRows(*csvFilePath)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correctAnswers := 0
	for i := 0; i < len(rows); i++ {
		currentProblem := problem{
			question: rows[i][0],
			answer:   rows[i][1],
		}
		displayQuestion(currentProblem.question, i+1)

		userAnswerCh := make(chan string)
		go func() {
			var givenAnswer string
			fmt.Scanf("%s\n", &givenAnswer)
			userAnswerCh <- givenAnswer
		}()

		select {
		case <-timer.C:
			fmt.Println("\nYour time is up!")
			fmt.Printf("\nYou got %d/%d correct answers!\n", correctAnswers, len(rows))
			return
		case ans := <-userAnswerCh:
			currentProblem.checkAnswerAndInc(ans, &correctAnswers)
		}
	}
	fmt.Printf("\nYou got %d/%d correct answers!\n", correctAnswers, len(rows))
}

func displayQuestion(question string, questionNumber int) {
	fmt.Printf("Problem #%d.	%s = ", questionNumber, question)
}

func readCsvRows(path string) [][]string {
	csvFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	rows, _ := csv.NewReader(csvFile).ReadAll()
	return rows
}

type problem struct {
	question string
	answer   string
}

// Bad practice or not? I can view it as a private extension function so as long as no other has this method may be good.
// The bad part is allowing a type to share state of "correctAnswers"
func (p *problem) checkAnswerAndInc(userAnswer string, correctAnswers *int) {
	if p.answer == userAnswer {
		*correctAnswers++
	}
}
