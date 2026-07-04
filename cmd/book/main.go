package main

import (
	"BookingTickets/internal/services"
	"log"
	"sync"
)

func main() {

	const totalTickets = 100
	const concurrentUsers = 2000000

	log.Printf("[System] Initializing Ticket Booking System. Inventory Capacity: %d\n", totalTickets)
	ticketManager := services.NewChannelTicketManager(totalTickets, 500)
	defer ticketManager.Close()

	var waitGroup sync.WaitGroup

	var successfulBookings = 0
	var failedBookings = 0
	var mutexLock sync.Mutex

	log.Printf("[System] Simulating %d concurrent booking attempts...\n", concurrentUsers)

	for i := 1; i <= concurrentUsers; i++ {
		waitGroup.Add(1)
		go func(userId int) {
			defer waitGroup.Done()
			success, msg := ticketManager.BookTickets(userId)

			mutexLock.Lock()
			if success {
				successfulBookings++
				log.Println(msg)
			} else {
				failedBookings++
			}
			mutexLock.Unlock()

		}(i)
	}
	waitGroup.Wait()

	log.Println("==================================================")
	log.Printf("[METRICS REPORT]")
	log.Printf("-> Total Intended Inventory: %d", totalTickets)
	log.Printf("-> Actual Successful Bookings: %d", successfulBookings)
	log.Printf("-> Total Deflected Oversell Attempts: %d", failedBookings)
	log.Printf("-> Final System Inventory State: %d", ticketManager.AvaibableTickets)
	log.Println("==================================================")

	if int(successfulBookings) > totalTickets {
		log.Fatalf("[CRITICAL ERROR] Race Condition Detected! System Oversold Tickets!")
	} else {
		log.Println("[STATUS] Concurrency Control Verification: SUCCESS. Data Integrity Preserved.")
	}

}
