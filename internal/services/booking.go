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
