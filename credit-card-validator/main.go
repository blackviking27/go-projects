package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type card struct {
	Number string `json:"number"`
}

type error struct {
	Error string `json:"error"`
}

type response struct {
	IsValid bool `json:"isValid"`
}

func reverse(s string) string {
	rev_string := ""
	for _, c := range s {
		rev_string = string(c) + rev_string
	}
	return rev_string
}

func luhnAlgo(card_num string) bool {

	odd_sum := 0
	even_sum := 0

	card_num = reverse(card_num)

	for i, _ := range card_num {
		val, err := strconv.Atoi(string(card_num[i]))
		if err != nil {
			panic(err)
		}
		if (i+1)%2 == 0 {
			// for even
			// double each num and if doubled num is double digit then add the individual digit
			x := val * 2
			if x >= 10 {
				even_sum += (1 + (x % 10))
			} else {
				even_sum += x
			}

		} else {
			// for odd nums
			// add all the odd nums
			odd_sum += val
		}
	}

	return (even_sum+odd_sum)%10 == 0
}

func mainScreen(w http.ResponseWriter, r *http.Request) {

	welcom_string := `
	Welcome to credit card validator in golang

	GET /verify -> to verify your credit card send card number in json format
	json data format => {"number": "{your credit card number}"}
	`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(welcom_string))
}

func verifyCardNum(w http.ResponseWriter, r *http.Request) {
	var card card

	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		resErr := error{Error: err.Error()}
		json.NewEncoder(w).Encode(resErr)
	}

	response := response{}
	response.IsValid = luhnAlgo(card.Number)
	json.NewEncoder(w).Encode(response)
}

func main() {
	// routes for the server
	http.HandleFunc("/", mainScreen)
	http.HandleFunc("/verify", verifyCardNum)

	fmt.Println("Started the server at port: 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error encountered:", err)
	}
}
