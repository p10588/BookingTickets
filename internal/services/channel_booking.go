package services

import (
	"fmt"
)

type BookResult struct {
	Success bool
	Message string
}

type BookRequest struct {
	UserId     int
	ResultChan chan BookResult
}

type ChannelTicketManager struct {
	AvaibableTickets int
	bookedCount      int
	requestChan      chan BookRequest
}

func NewChannelTicketManager(initalTickets int, queueSize int) *ChannelTicketManager {
	cm := &ChannelTicketManager{
		AvaibableTickets: initalTickets,
		requestChan:      make(chan BookRequest, queueSize),
	}

	go cm.startEventLoop()

	return cm
}

func (cm *ChannelTicketManager) BookTickets(userId int) (bool, string) {

	resultChan := make(chan BookResult)

	cm.requestChan <- BookRequest{
		UserId:     userId,
		ResultChan: resultChan,
	}

	result := <-resultChan

	return result.Success, result.Message
}

func (cm *ChannelTicketManager) Close() {
	close(cm.requestChan)
}

func (cm *ChannelTicketManager) startEventLoop() {
	for req := range cm.requestChan {
		if cm.AvaibableTickets <= 0 {
			req.ResultChan <- BookResult{Success: false, Message: "No tickets available."}
			continue
		}
		cm.AvaibableTickets--
		cm.bookedCount++
		req.ResultChan <- BookResult{
			Success: true,
			Message: fmt.Sprintf("Ticket booked successfully for user %d", req.UserId),
		}
	}
}
