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

type Transaction struct {
	FromID string
	ToID   string
	Amount float64
}

type PaymentService struct {
	Users        map[string]*User
	Transactions []Transaction
}

func NewPaymentService() *PaymentService {
	return &PaymentService{
		Users:        make(map[string]*User),
		Transactions: make([]Transaction, 0),
	}
}

func (p *PaymentService) AddUser(user *User) {
	p.Users[user.ID] = user
}

func (p *PaymentService) AddTransaction(t Transaction) {
	p.Transactions = append(p.Transactions, t)
}

func (p *PaymentService) ProcessingTransactions(t Transaction) error {
	recipient, ok := p.Users[t.ToID]

	if !ok {
		return fmt.Errorf("Recipient not found: ID %s\n", t.ToID)
	}

	sender, ok := p.Users[t.FromID]

	if !ok {
		return fmt.Errorf("Sender not found: ID %s\n", t.FromID)
	}

	if err := sender.Withdraw(t.Amount); err != nil {
		return err
	}

	if err := recipient.Deposit(t.Amount); err != nil {
		sender.Deposit(t.Amount)
		return err
	}

	return nil
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

	p := NewPaymentService()
	p.AddUser(&u1)
	p.AddUser(&u2)

	t1 := Transaction{
		u1.ID,
		u2.ID,
		4200,
	}

	t2 := Transaction{
		u2.ID,
		u1.ID,
		1000,
	}

	t3 := Transaction{
		u2.ID,
		u1.ID,
		300,
	}

	p.AddTransaction(t1)
	p.AddTransaction(t2)
	p.AddTransaction(t3)

	for _, t := range p.Transactions {
		if err := p.ProcessingTransactions(t); err != nil {
			fmt.Printf("Transaction error: %s\n ", err)
		}
	}

	u1.PrintUserBalance()
	u2.PrintUserBalance()

}
