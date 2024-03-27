package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Transaction struct {
	HumanId uint `json:"humanId"`
	Amount  uint `json:"amount"`
}

type UserAccount struct {
	HumanId uint `json:"humanId"`
	Amount  uint `json:"amount"`
}

func main() {
	println("Do you want to enter transactions or calculate the result of present data? (INPUT / RESULT): ")
	var mode string

	_, err := fmt.Scanln(&mode)

	if err != nil {
		println("Error Reading Input: ", err)
		return
	}

	mode = strings.ToUpper(mode)

	switch mode {
	case "INPUT":
		inputUser()
	case "RESULT":
		calc()
	default:
		fmt.Println("Invalid mode.")
	}
}

type JsonData struct {
	Humans       []uint        `json:"humans"`
	Transactions []Transaction `json:"transactions"`
}

func inputUser() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter user ids, separated by commas: ")
	humanInput, _ := reader.ReadString('\n')

	humanIdsStr := strings.Split(strings.TrimSpace(humanInput), ",")
	var humanIds []uint

	for _, idStr := range humanIdsStr {
		id, err := strconv.ParseUint(idStr, 10, 64)

		if err != nil {
			panic("Invalid user ID")
		}

		humanIds = append(humanIds, uint(id))
	}

	var transactions []Transaction

	println("---------------------")
	println("Enter transactions as 'humanId,amount' or type 'exit' to finish: ")
	for {
		transactionInput, _ := reader.ReadString('\n')
		transactionInput = strings.TrimSpace(transactionInput)
		if transactionInput == "" {
			break
		}

		parts := strings.Split(transactionInput, ",")
		if len(parts) != 2 {
			log.Fatal("Incorrect input, expected 'humanId,amount'")
		}

		humanId, err := strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			panic("Invalid use id")
		}

		amount, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			panic("Invalid Amount")
		}

		transactions = append(transactions, Transaction{HumanId: uint(humanId), Amount: uint(amount)})
	}

	data := JsonData{
		Humans:       humanIds,
		Transactions: transactions,
	}

	file, _ := json.MarshalIndent(data, "", " ")

	_ = os.WriteFile("data.json", file, 0644)
}

func loadData() ([]uint, []Transaction) {
	data, err := os.ReadFile("data.json")
	if err != nil {
		log.Fatal(err)
	}

	var loadedData JsonData
	err = json.Unmarshal(data, &loadedData)

	if err != nil {
		panic("Invalid JSON data")
	}

	return loadedData.Humans, loadedData.Transactions
}

func calc() {
	humans, transactions := loadData()
	var totalPrice uint = 0

	overview := map[uint]uint{}

	// Pre-generate overview
	for _, h := range humans {
		overview[h] = 0
	}

	// Check for impostors that are not part of the humans list
	for _, t := range transactions {
		var found bool
		for _, h := range humans {
			if t.HumanId == h {
				found = true
			}
		}

		if !found {
			panic("There is an impostor")
		}

		// Since one human can have many transactions, we have to sum them app
		overview[t.HumanId] = overview[t.HumanId] + t.Amount

		// And also increase the total price, so we can calculate the average later on
		totalPrice += t.Amount
	}

	var average = totalPrice / uint(len(humans))
	println("---------------------")
	println("Average amount per user: " + strconv.Itoa(int(average)))
	println("---------------------")

	var overpaid []UserAccount
	var underpaid []UserAccount

	// No we'll split all the users into two groups, those who have to pay their delta to the average
	// And those who paid above the average
	for humanId, amount := range overview {
		if amount > average {
			overpaid = append(overpaid, UserAccount{
				HumanId: humanId,
				Amount:  amount,
			})
		} else if amount < average {
			underpaid = append(underpaid, UserAccount{
				HumanId: humanId,
				Amount:  amount,
			})
		}
	}

	sort.Slice(overpaid, func(i, j int) bool {
		return overpaid[i].Amount > overpaid[j].Amount
	})

	sort.Slice(underpaid, func(i, j int) bool {
		return underpaid[i].Amount > underpaid[j].Amount
	})

	// Loop through users, which paid less than average (or nothing)
	for _, human := range underpaid {
		// Recursively split the amount to users which have overpaid
		overpaid = payUser(human, average, overpaid)
	}
}

func payUser(human UserAccount, average uint, overpaid []UserAccount) []UserAccount {
	delta := average - human.Amount
	humanIndex := findFirstHumanIndex(&overpaid, average)
	foundHuman := overpaid[humanIndex]
	availableToPay := foundHuman.Amount - average

	if delta > availableToPay {
		overpaid[humanIndex] = UserAccount{
			HumanId: foundHuman.HumanId,
			Amount:  foundHuman.Amount - availableToPay,
		}

		println("User: " + strconv.Itoa(int(human.HumanId)) + " will pay user: " + strconv.Itoa(int(foundHuman.HumanId)) + " amount: " + strconv.Itoa(int(availableToPay)))

		delta = delta - availableToPay
		return payUser(UserAccount{HumanId: human.HumanId, Amount: average - delta}, average, overpaid)
	} else {
		overpaid[humanIndex] = UserAccount{
			HumanId: foundHuman.HumanId,
			Amount:  foundHuman.Amount - delta,
		}

		println("User: " + strconv.Itoa(int(human.HumanId)) + " will pay user: " + strconv.Itoa(int(foundHuman.HumanId)) + " amount: " + strconv.Itoa(int(delta)))
	}

	return overpaid
}

func findFirstHumanIndex(humans *[]UserAccount, average uint) int {
	for i, h := range *humans {
		if h.Amount != average {
			return i
		}
	}

	panic("No user found | If this happens, that means I f*cked something app and the total count doesn't add up :D")
}
