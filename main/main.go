package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	csvPath = flag.String("csv", "problems.csv", "csv file question/answer'")
	timeSeconds = flag.Int("time", 10, "number of seconds for test duration")
	finish = make(chan bool)
	questions int
	correctAnswers int
)

func main() {
	flag.Parse()
	settings := Config{*csvPath, *timeSeconds }
	timer := time.NewTicker(time.Second * time.Duration(settings.GetTimeDuration()))
	quizes := readCsvFile(settings)
	go processQuizes(quizes)
	select {
		case <- finish:
			fmt.Println("Test has been finished")
			fmt.Printf("Total questions %d \n", questions)
			fmt.Printf("Total correct answers %d \n", correctAnswers)
		case <- timer.C :
			fmt.Println("Time is over")
	}
}
func processQuizes(quizes []Quiz){
	scanner :=  bufio.NewScanner(os.Stdin)
	for _, quiz := range quizes {
		questions++
		fmt.Println("What is the answer for "+ quiz.GetQuestion())
		scanner.Scan()
		answer := scanner.Text()
		if quiz.IsCorrect(answer) {
			correctAnswers++
		}
	}
	finish <- true
}

func readCsvFile(settings Config) ([] Quiz) {
	file, err := os.Open(settings.GetFilename())
	if err != nil {
		fmt.Println("Error during csv file reading.")
		panic(err)

	}
	defer file.Close()
	csvReader := csv.NewReader(file)
	csvData, err := csvReader.ReadAll()
	quizes := make([]Quiz, 0, len(csvData))
	for _, value := range csvData {
		quiz := Quiz{value[0],value[1] }
		quizes = append(quizes,quiz)
	}
	return quizes
}

type Config struct {
	Filename string
	SecondsLimit int
}

func (conf *Config) GetFilename() string{
	return conf.Filename
}

func (conf *Config) GetTimeDuration() int {
	return conf.SecondsLimit
}

type Quiz struct {
	Question string
	Answer string
}

func (q *Quiz) GetQuestion() string  {
	return q.Question
}

func (q *Quiz) GetAnswer() string  {
	return q.Answer
}

func (q *Quiz) IsCorrect(answer string) bool  {
	answer = q.formatData(answer)
	correctAnswer := q.formatData(q.GetAnswer())
	if answer == correctAnswer {
		return true
	}
	return false
}

func (q *Quiz) formatData(data string) string {
	data = strings.ToLower(data)
	data = strings.TrimSpace(data)
	return data
}