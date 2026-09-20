package main

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrNotEnoughBalance     = errors.New("balance is not enough")
	ErrAmountMustBePositive = errors.New("amount must be positive")
)

type User struct {
	ID      string
	Balance float64
	Name    string
	mu      sync.Mutex
}

func (u *User) Deposit(amount float64) error {
	if amount < 0 {
		return ErrAmountMustBePositive
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	u.Balance += amount

	return nil
}

func (u *User) Withdraw(amount float64) error {
	if amount < 0 {
		return ErrAmountMustBePositive
	}

	u.mu.Lock()
	defer u.mu.Unlock()

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

func (p *PaymentService) worker(tq <-chan Transaction, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range tq {
		if err := p.ProcessingTransactions(t); err != nil {
			fmt.Printf("Transaction error: %s\n", err)
		}
	}
}

func main() {
	u1 := User{
		ID:      "1",
		Name:    "user1",
		Balance: 1500.00,
	}

	u2 := User{
		ID:      "2",
		Name:    "user2",
		Balance: 2000.00,
	}
	ps := NewPaymentService()

	ps.AddUser(&u1)
	ps.AddUser(&u2)

	t1 := Transaction{
		FromID: u1.ID,
		ToID:   u2.ID,
		Amount: 200,
	}
	t2 := Transaction{
		FromID: u1.ID,
		ToID:   u2.ID,
		Amount: 300,
	}

	t3 := Transaction{
		FromID: u1.ID,
		ToID:   u2.ID,
		Amount: 500,
	}

	t4 := Transaction{
		FromID: u2.ID,
		ToID:   u1.ID,
		Amount: 5000,
	}

	ps.AddTransaction(t1)
	ps.AddTransaction(t2)
	ps.AddTransaction(t3)
	ps.AddTransaction(t4)

	var wg sync.WaitGroup
	ch := make(chan Transaction, len(ps.Transactions))

	for _, t := range ps.Transactions {
		ch <- t
	}

	close(ch)

	for i := 0; i <= 4; i++ {
		wg.Add(1)
		go ps.worker(ch, &wg)
	}

	wg.Wait()

	u1.PrintUserBalance()
	u2.PrintUserBalance()
}
