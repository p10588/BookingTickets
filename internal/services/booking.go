// CONCURRENCY PATTERN 1: STRUCTURAL MUTEX BLOCKING (SHARED MEMORY)
//
// [Analogy]: Millions of users rush the ticket hall at once. There is a single locked door
// (Mutex) — only one person enters at a time, everyone else crowds and shoves outside.
// The door prevents chaos inside, but the crowd fighting for the handle outside is pure waste.
//
// [Technical]: High lock contention causes the OS to continuously context-switch between
// blocked goroutines, burning CPU cycles on thread lifecycle management instead of executing
// actual business logic.
//
// [Trade-off]: Data safety is guaranteed, but throughput is bottlenecked — only one goroutine
// runs at a time while the rest queue up waiting for the lock.

package services

import (
	"fmt"
	"sync"
)

type TicketManager struct {
	mu               sync.Mutex
	AvailableTickets int
	BookedCount      int
}

func NewTicketManager(intialTickets int) *TicketManager {
	return &TicketManager{
		AvailableTickets: intialTickets,
	}
}

func (tm *TicketManager) BookTickets(userId int) (bool, string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if tm.AvailableTickets <= 0 {
		return false, "Oversold Blocked: No tickets remaining."
	}

	tm.AvailableTickets--
	tm.BookedCount++
	return true, fmt.Sprintf("Success: User %d secured ticket. Remaining: %d", userId, tm.AvailableTickets)
}
