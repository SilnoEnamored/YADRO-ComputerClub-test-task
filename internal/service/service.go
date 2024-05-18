package service

import (
	"bufio"
	"club/internal/domain"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Process читает входной файл, разбирает данные и обрабатывает события
func Process(inputFile string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var club domain.Club

	// Чтение количества столов
	if scanner.Scan() {
		club.Tables, err = strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("invalid table count format: %s", scanner.Text())
		}
	}

	// Чтение времени работы клуба
	if scanner.Scan() {
		times := strings.Split(scanner.Text(), " ")
		if len(times) != 2 {
			return fmt.Errorf("invalid working hours format: %s", scanner.Text())
		}
		club.StartTime, err = time.Parse("15:04", times[0])
		if err != nil {
			return fmt.Errorf("invalid start time format: %s", times[0])
		}
		club.EndTime, err = time.Parse("15:04", times[1])
		if err != nil {
			return fmt.Errorf("invalid end time format: %s", times[1])
		}
		if club.StartTime.After(club.EndTime) {
			return fmt.Errorf("start time is after end time: %s", scanner.Text())
		}
	}

	// Чтение стоимости часа
	if scanner.Scan() {
		club.Price, err = strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("invalid price format: %s", scanner.Text())
		}
	}

	// Чтение событий
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		if len(parts) < 3 {
			return fmt.Errorf("invalid event format: %s", scanner.Text())
		}

		eventTime, err := time.Parse("15:04", parts[0])
		if err != nil {
			return fmt.Errorf("invalid event time format: %s", parts[0])
		}
		eventID, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("invalid event ID format: %s", parts[1])
		}
		event := domain.Event{
			Time:   eventTime,
			ID:     eventID,
			Client: parts[2],
		}

		if eventID == 2 || eventID == 12 {
			if len(parts) != 4 {
				return fmt.Errorf("invalid event format for table ID: %s", scanner.Text())
			}
			event.TableID, err = strconv.Atoi(parts[3])
			if err != nil {
				return fmt.Errorf("invalid table ID format: %s", parts[3])
			}
		} else {
			if len(parts) != 3 {
				return fmt.Errorf("invalid event format: %s", scanner.Text())
			}
		}

		club.Events = append(club.Events, event)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return processEvents(club)
}

// processEvents обрабатывает события и выводит результат в консоль
func processEvents(club domain.Club) error {
	// Вывод времени начала работы клуба
	fmt.Println(club.StartTime.Format("15:04"))

	// Инициализация данных для отслеживания клиентов, столов и ожидания
	clientsInClub := make(map[string]time.Time)
	clientsAtTable := make(map[int]string)
	clientsWaiting := []string{}
	tableRevenue := make(map[int]int)
	tableUsage := make(map[int]time.Duration)

	// Обработка каждого события
	for _, event := range club.Events {
		switch event.ID {
		case 1:
			// Клиент пришел
			if event.Time.Before(club.StartTime) {
				printError(event, "NotOpenYet")
			} else if _, exists := clientsInClub[event.Client]; exists {
				printError(event, "YouShallNotPass")
			} else {
				clientsInClub[event.Client] = event.Time
				printEvent(event)
			}
		case 2:
			// Клиент сел за стол
			if _, exists := clientsInClub[event.Client]; !exists {
				printError(event, "ClientUnknown")
			} else if occupant, occupied := clientsAtTable[event.TableID]; occupied && occupant != event.Client {
				printError(event, "PlaceIsBusy")
			} else {
				if previousTableID, ok := findClientAtTable(clientsAtTable, event.Client); ok {
					duration := event.Time.Sub(clientsInClub[event.Client])
					tableUsage[previousTableID] += duration
					tableRevenue[previousTableID] += int(duration.Hours()+0.9999) * club.Price
					delete(clientsAtTable, previousTableID)
				}
				clientsAtTable[event.TableID] = event.Client
				clientsInClub[event.Client] = event.Time
				printEvent(event)
			}
		case 3:
			// Клиент ожидает
			if len(clientsAtTable) < club.Tables {
				printError(event, "ICanWaitNoLonger!")
			} else {
				clientsWaiting = append(clientsWaiting, event.Client)
				printEvent(event)
			}
		case 4:
			// Клиент ушел
			if _, exists := clientsInClub[event.Client]; !exists {
				printError(event, "ClientUnknown")
			} else {
				if previousTableID, ok := findClientAtTable(clientsAtTable, event.Client); ok {
					duration := event.Time.Sub(clientsInClub[event.Client])
					tableUsage[previousTableID] += duration
					tableRevenue[previousTableID] += int(duration.Hours()+0.9999) * club.Price
					delete(clientsAtTable, previousTableID)
					printEvent(event)
					if len(clientsWaiting) > 0 {
						nextClient := clientsWaiting[0]
						clientsWaiting = clientsWaiting[1:]
						clientsAtTable[previousTableID] = nextClient
						clientsInClub[nextClient] = event.Time
						printEvent(domain.Event{Time: event.Time, ID: 12, Client: nextClient, TableID: previousTableID})
					}
				} else {
					printEvent(event)
				}
				delete(clientsInClub, event.Client)
			}
		}
	}

	// Обработка оставшихся клиентов в конце рабочего дня
	remainingClients := []string{}
	for client := range clientsInClub {
		remainingClients = append(remainingClients, client)
	}
	sort.Strings(remainingClients)

	for _, client := range remainingClients {
		if previousTableID, ok := findClientAtTable(clientsAtTable, client); ok {
			duration := club.EndTime.Sub(clientsInClub[client])
			tableUsage[previousTableID] += duration
			tableRevenue[previousTableID] += int(duration.Hours()+0.9999) * club.Price
			delete(clientsAtTable, previousTableID)
		}
		printEvent(domain.Event{Time: club.EndTime, ID: 11, Client: client})
	}

	// Вывод времени окончания работы клуба
	fmt.Println(club.EndTime.Format("15:04"))

	// Вывод статистики по столам
	for i := 1; i <= club.Tables; i++ {
		hours := int(tableUsage[i].Hours())
		minutes := int(tableUsage[i].Minutes()) % 60
		fmt.Printf("%d %d %02d:%02d\n", i, tableRevenue[i], hours, minutes)
	}

	return nil
}

// printEvent выводит событие в консоль
func printEvent(event domain.Event) {
	fmt.Printf("%s %d %s", event.Time.Format("15:04"), event.ID, event.Client)
	if event.ID == 2 || event.ID == 12 {
		fmt.Printf(" %d", event.TableID)
	}
	fmt.Println()
}

// printError выводит событие и соответствующую ошибку в консоль
func printError(event domain.Event, errorMsg string) {
	fmt.Printf("%s %d %s", event.Time.Format("15:04"), event.ID, event.Client)
	if event.ID == 2 || event.ID == 12 {
		fmt.Printf(" %d", event.TableID)
	}
	fmt.Println()
	fmt.Printf("%s 13 %s\n", event.Time.Format("15:04"), errorMsg)
}

// findClientAtTable находит стол, за которым сидит клиент
func findClientAtTable(clientsAtTable map[int]string, client string) (int, bool) {
	for tableID, occupant := range clientsAtTable {
		if occupant == client {
			return tableID, true
		}
	}
	return 0, false
}
