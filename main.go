package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	// "time"
	"bufio"
	"os"
	"strconv"
	"strings"
)

var scanner *bufio.Scanner

func readLine() string {
	if scanner.Scan() {
		line := scanner.Text()
		return line
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	log.Fatal("Ошибка ввода")
	return ""
}

func pickDifficulty() (size int, attempts int) {
	fmt.Println("Выберите сложность: Easy, Medium, Hard")
	for true {
		line := readLine()
		line = strings.ToLower(line)
		if line == "easy" {
			return 50, 15
		} else if line == "medium" {
			return 100, 10
		} else if line == "hard" {
			return 200, 5
		} else {
			fmt.Println("Такой сложности нет, попробуйте еще раз")
		}
	}
	return 0, 0
}

func getDistanceHint(target int, guess int) string {
	diff := target - guess
	if diff < 0 { diff = -diff } // get abs distance

	if diff <= 5 {
		return "Горячо 🔥"
	} else if diff <= 15 {
		return "Тепло 🙂"
	} else {
		return "Холодно ❄️" 
	}
}

type OnePlayResult struct {
	Date time.Time
	DidWin bool
	Attempts int
}

type ResultsJson struct {
	Plays []OnePlayResult
}

var results ResultsJson

func loadResults() {
	jsonFile, e := os.Open("results.json")
	if e != nil {
		fmt.Println("Не удалось загрузить results.json")
		log.Print(e)
		return
	}
	e = json.NewDecoder(jsonFile).Decode(&results)
	if e != nil {
		fmt.Println("Не удалось загрузить results.json")
		log.Print(e)
		results = ResultsJson{}
	}
	jsonFile.Close()
}

func saveResults() {
	jsonFile, e := os.OpenFile("results_new.json", os.O_CREATE | os.O_RDWR | os.O_TRUNC, 0644)
	if e != nil {
		log.Fatal(e)
	}
	e = json.NewEncoder(jsonFile).Encode(results)
	if e != nil {
		fmt.Printf("Не удалось сохранить results.json")
		log.Print(e)
	}
	jsonFile.Close()
	os.Rename("results_new.json", "results.json")
}

func main() {
	scanner = bufio.NewScanner(os.Stdin)
	
	loadResults()

	for true {
		targetSize, attempts := pickDifficulty()
		startingAttempts := attempts

		fmt.Printf("Игра 'Угадай число' - от 1 до %v началась!\n", targetSize)
		fmt.Printf("Угадайте число за %v попыток!\n", attempts)

		target := rand.Intn(targetSize) + 1
		attemptHistory := []int{}

		didWin := false

		for attempts > 0 && !didWin {
			line := readLine()
			guess, e := strconv.Atoi(line)
			if e != nil {
				fmt.Println("Введите число")
			} else {
				attemptHistory = append(attemptHistory, guess)
				if guess == target {
					didWin = true
				} else if guess < target {
					fmt.Println("Секретное число больше👆")
					fmt.Println(getDistanceHint(target, guess))
				} else {
					fmt.Println("Секретное число меньше👇")
					fmt.Println(getDistanceHint(target, guess))
				}
			}
			attempts -= 1
		}
		usedAttempts := startingAttempts - attempts
		if didWin {
			fmt.Printf("Вы выиграли за %v попыток\n", usedAttempts)
		} else {
			fmt.Println("Попытки закончились, вы проебали")
			fmt.Printf("Секретное число: %v\n", target)
		}
		fmt.Println("Игра закончена!")

		results.Plays = append(results.Plays, OnePlayResult{
			Date: time.Now(),
			Attempts: usedAttempts,
			DidWin: didWin,
		})
		saveResults()

		keepGoing := true
		for true {
			fmt.Println("Сыграть еще раз, y/n?")
			line := strings.ToLower(readLine())
			if line == "n" {
				keepGoing = false
				break
			} else if line == "y" {
				break
			}
		}

		if !keepGoing {
			break
		}
	}
}
