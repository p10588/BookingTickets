// CONCURRENCY PATTERN 2: BUFFERED CHANNEL PIPELINE (CSP MODEL WITHOUT FAST-FAIL)
//
// [Analogy]: The ticket hall sets up velvet ropes (Channel Queue) — no more pushing, just an
// orderly line. But there is only one ticket counter (a single Event Loop goroutine), and the
// queue stretches 2 million users long, so each person waits 6 seconds just to reach the
// counter, only to hear "sold out in the first millisecond." Worse, once the queue itself is
// full, newcomers freeze at the entrance and cannot even join — they block the doorway until
// a slot opens up.
//
// [Technical]: A single event loop goroutine serializes all requests through the channel,
// eliminating shared-memory race conditions entirely. However, without a fast-fail check
// before enqueue, rejected requests must traverse the full pipeline before being turned away,
// causing backpressure latency. A full buffer also blocks the caller's goroutine at the send.
//
// [Trade-off]: Race conditions are eliminated, but throughput is limited to a single worker
// and latency grows linearly with queue depth — every request pays the full pipeline cost
// even when tickets are already sold out.

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
