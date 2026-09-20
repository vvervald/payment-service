package main

import (
	"errors"
	"fmt"
)

var (
	ErrNotEnoughBalance     = errors.New("balance is not enough")
	ErrAmountMustBePositive = errors.New("amount must be positive")
)

type User struct {
	ID      string
	Balance float64
	Name    string
}

func (u *User) Deposit(amount float64) error {
	if amount < 0 {
		return ErrAmountMustBePositive
	}

	u.Balance += amount

	return nil
}

func (u *User) Withdraw(amount float64) error {
	if amount < 0 {
		return ErrAmountMustBePositive
	}

	if u.Balance-amount < 0 {
		return ErrNotEnoughBalance
	}

	u.Balance -= amount
	return nil
}

func (u *User) PrintUserBalance() {
	fmt.Printf("User: %s, Balance: %.2f\n", u.ID, u.Balance)
}

func main() {
	u1 := User{
		ID:      "1",
		Name:    "user1",
		Balance: 1500.00,
	}

	u2 := User{
		ID:      "2",
		Name:    "user1",
		Balance: 2000.00,
	}

	if err := u1.Deposit(599); err != nil {
		fmt.Println("Operation declined:", err)
	}

	u1.PrintUserBalance()

	if err := u1.Withdraw(100); err != nil {
		fmt.Println("Operation declined:", err)
	}

	u1.PrintUserBalance()

	if err := u2.Deposit(2000); err != nil {
		fmt.Println("Operation declined:", err)
	}

	u2.PrintUserBalance()

	if err := u2.Withdraw(5000); err != nil {
		fmt.Println("Operation declined:", err)
	}

	u2.PrintUserBalance()
}
