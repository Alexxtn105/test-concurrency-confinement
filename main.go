package main

import (
	"fmt"
	"sync"
)

/*
manageTicket Функция, которая отвечает за управление распределением заявок
и реагирование на запросы пользователей.
Она прослушивает входящие запросы по каналу (ticketChan) и подает сигналы
по другому каналу (doneChan), когда пора остановиться.
*/
func manageTicket(ticketChan chan int, doneChan chan bool, tickets *int) {
	// крутим бесконечный цикл
	for {
		select {
		case user := <-ticketChan: // принимаем запрос на покупку билета
			if *tickets > 0 {
				*tickets--
				fmt.Printf("User %d purchased a ticket. Tickets remaining: %d\n",
					user, *tickets)
			} else {
				fmt.Printf("User %d found no tickets.\n", user)
			}
		case <-doneChan: // принимаем сигнал о том, что все билеты проданы
			fmt.Printf("Tickets remaining: %d\n", *tickets)
		}

	}
}

/*
buyTicket Функция, которая имитирует попытку пользователя купить билет.
Она отправляет запрос в функцию manageTicket через ticketChan.
*/
func buyTicket(wg *sync.WaitGroup, ticketChan chan int, userId int) {
	defer wg.Done()
	ticketChan <- userId
}

func main() {
	var wg sync.WaitGroup        // WaitGroup для ожидания завершения работы всех горутин
	tickets := 500               // Общее количество доступных билетов
	ticketChan := make(chan int) // Канал для отправки запросов на покупку билета
	doneChan := make(chan bool)  // Канал для сигнализации остановки

	go manageTicket(ticketChan, doneChan, &tickets)

	for userId := 0; userId < 1000; userId++ {
		wg.Add(1)
		go buyTicket(&wg, ticketChan, userId)
	}

	wg.Wait()
	// посылаем в канал сигнал остановки
	doneChan <- true
}
