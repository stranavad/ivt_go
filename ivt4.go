package main

import "strconv"

type Transaction struct {
	HumanId uint
	Amount  uint
}

type UserAccount struct {
	HumanId uint
	Amount  uint
}

func main() {
	humans := []uint{1, 2, 3, 4, 5, 6, 7, 8}
	transactions := []Transaction{
		{HumanId: 1, Amount: 69},
		{HumanId: 2, Amount: 250},
		{HumanId: 1, Amount: 1},
		{HumanId: 1, Amount: 120},
		{HumanId: 4, Amount: 250},
		{HumanId: 4, Amount: 250},
		{HumanId: 7, Amount: 85},
	}
	var totalPrice uint = 0

	overview := map[uint]uint{}

	// Pre-generate overview
	for _, h := range humans {
		overview[h] = 0
	}

	// Check for impostors
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

		overview[t.HumanId] = overview[t.HumanId] + t.Amount
		totalPrice += t.Amount
	}

	var average uint = totalPrice / uint(len(humans))
	println(average)

	var overpayed []UserAccount
	var underpayed []UserAccount

	// Since we've already pre-generated the map with data, we don't have to check for 0 or nil values
	// When accessing the map
	for humanId, amount := range overview {
		if amount > average {
			overpayed = append(overpayed, UserAccount{
				HumanId: humanId,
				Amount:  amount,
			})
		} else if amount < average {
			underpayed = append(underpayed, UserAccount{
				HumanId: humanId,
				Amount:  amount,
			})
		}
	}

	for _, human := range underpayed {
		overpayed = payUser(human, average, overpayed)
	}
}

func payUser(human UserAccount, average uint, overpayed []UserAccount) []UserAccount {
	delta := average - human.Amount
	humanIndex := findFirstHumanIndex(&overpayed, average)
	foundHuman := overpayed[humanIndex]
	availableToPay := foundHuman.Amount - average

	if delta > availableToPay {
		overpayed[humanIndex] = UserAccount{
			HumanId: foundHuman.HumanId,
			Amount:  foundHuman.Amount - availableToPay,
		}
		println("User: " + strconv.Itoa(int(human.HumanId)) + " will pay user" + strconv.Itoa(int(foundHuman.HumanId)) + " Amount: " + strconv.Itoa(int(availableToPay)))
		delta = delta - availableToPay
		return payUser(UserAccount{HumanId: human.HumanId, Amount: average - delta}, average, overpayed)
	} else {
		overpayed[humanIndex] = UserAccount{
			HumanId: foundHuman.HumanId,
			Amount:  foundHuman.Amount - delta,
		}
		println("User: " + strconv.Itoa(int(human.HumanId)) + " will pay user" + strconv.Itoa(int(foundHuman.HumanId)) + " Amount: " + strconv.Itoa(int(delta)))
	}

	return overpayed
}

func findFirstHumanIndex(humans *[]UserAccount, average uint) int {
	for i, h := range *humans {
		if h.Amount != average {
			return i
		}
	}

	panic("No user found")
}
