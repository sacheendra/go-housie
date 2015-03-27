package main

import (
	"fmt"
	"math/rand"
)

type PersonID int

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Person(called_number chan int, claim_prize chan PersonID, received chan bool, coupon []int, person_id PersonID) {
	numFound := 0
	for i := 0; i < len(coupon); i++ {
		current_number := <-called_number
		found := contains(coupon, current_number)
		if found {
			numFound++
		}
		if numFound == len(coupon) {
			claim_prize <- person_id
		} else {
			received <- true
		}
	}
}

func main() {
	var called_number chan int
	var claim_prize chan PersonID
	var received chan bool

	tokens := make([][]int, 3)
	for i := 0; i < 3; i++ {
		tokens[i] = make([]int, 12)
		for j := 0; j < 12; j++ {
			num := rand.Intn(100) + 1
			found := contains(tokens[i], num)
			for found {
				num = rand.Intn(100) + 1
				found = contains(tokens[i], num)
			}
			tokens[i][j] = num
		}
	}

	go Person(called_number, claim_prize, received, tokens[0], 0)
	go Person(called_number, claim_prize, received, tokens[1], 1)
	go Person(called_number, claim_prize, received, tokens[2], 2)

	claimants := make([]PersonID, 0)
	prev_called := make(map[int]bool)
	for i := 0; i < 100; i++ {
		if len(claimants) == 3 {
			break
		}
		num := rand.Intn(100) + 1
		_, ok := prev_called[num]
		for ok {
			num = rand.Intn(100) + 1
			_, ok = prev_called[num]
		}
		prev_called[num] = true
		called_number <- num
		for j := 0; j < 3; j++ {
			select {
			case _ = <-received:
				continue
			case pid := <-claim_prize:
				claimants = append(claimants, pid)
			}
		}
	}

	fmt.Println(claimants)
}
