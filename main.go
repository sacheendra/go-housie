package main

import (
	"fmt"
	"math/rand"
)

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Person(called_number chan int, claim_prize chan int, received chan bool, coupon []int, person_id int) {
	numFound := 0
	for current_number := range called_number {
		if contains(coupon, current_number) {
			numFound++
		}
		if numFound == len(coupon) {
			fmt.Println(person_id)
			claim_prize <- person_id
		} else {
			received <- true
		}
	}
}

func main() {
	var (
		called_number = make([]chan int, 3)
		claim_prize   = make(chan int, 1)
		received      = make(chan bool, 1)
	)

	called_number[0] = make(chan int, 1)
	called_number[1] = make(chan int, 1)
	called_number[2] = make(chan int, 1)

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

	for i := 0; i < 3; i++ {
		go Person(called_number[i], claim_prize, received, tokens[i], i)
	}

	claimants := make([]int, 0)
	prev_called := make(map[int]bool)
	for i := 0; i < 100; i++ {
		if len(claimants) == 3 {
			close(received)
			close(claim_prize)
			break
		}
		num := rand.Intn(100) + 1
		_, ok := prev_called[num]
		for ok {
			num = rand.Intn(100) + 1
			_, ok = prev_called[num]
		}
		prev_called[num] = true

		for k := 0; k < 3; k++ {
			if !contains(claimants, k) {
				called_number[k] <- num
			}
		}
		for j := 0; j < 3-len(claimants); j++ {
			select {
			case _ = <-received:
				continue
			case pid := <-claim_prize:
				if pid == 0 {
					close(called_number[0])
				}
				if pid == 1 {
					close(called_number[1])
				}
				if pid == 2 {
					close(called_number[2])
				}
				claimants = append(claimants, pid)
			}
		}
	}
	fmt.Println(claimants)
}
