package main

// Вариант 38
// Разработать программу для обеспечения продажи билетов в кинотеатр.

import (
	"fmt"
	"log"
	"time"
)

// DONE: добавить кинотеатр с залами и с капасити

type Cinema struct {
	halls    []*Hall
	vipHalls []*VipHall
}

type Hall struct {
	// key - movieMap, value - countOfBoughtTickets
	movieMap map[Movie]int
	capacity int
	HallI
}

type VipHall struct {
	*Hall
	addedCostCoefficient int
	HallI
}

type HallI interface {
	// return hall and value of added cost coefficient
	getCurrentHall() (*Hall, int)
	getAddedCostCoefficient() int
}

func getAddedCostCoefficient(hallI HallI) int {
	return hallI.getAddedCostCoefficient()
}

func (hall *Hall) getAddedCostCoefficient() int {
	// hall don't have an added cost
	return 1
}

func (vipHall *VipHall) getAddedCostCoefficient() int {
	return vipHall.addedCostCoefficient
}

func (hall *Hall) getCurrentHall() (*Hall, int) {
	return hall, getAddedCostCoefficient(hall)
}

func (vipHall *VipHall) getCurrentHall() (*Hall, int) {
	return vipHall.Hall, getAddedCostCoefficient(vipHall)
}

type Movie struct {
	title       string
	startTime   time.Time
	duration    time.Duration
	ticketPrice int
}

type Tickets struct {
	movie          Movie
	price          int
	countOfTickets int
}

type BasicCustomer struct {
	CustomerI
	name    string
	balance int
	tickets map[Movie]Tickets
}

type Beneficiary struct {
	CustomerI
	name    string
	balance int
	tickets map[Movie]Tickets
}

type CustomerI interface {
	buyTicket(hall HallI, movie Movie, countOfTickets int) error
}

func buyTicket(customerI CustomerI, hall HallI, movie Movie, countOfTickets int) {
	err := customerI.buyTicket(hall, movie, countOfTickets)
	if err != nil {
		log.Println(err)
	}
}

// Обычный покупатель покупает за полную стоимость
func (c *BasicCustomer) buyTicket(hallI HallI, movie Movie, countOfTickets int) error {
	hall, addedCostCoefficient := hallI.getCurrentHall()

	if hall.movieMap[movie]+countOfTickets > hall.capacity {
		return fmt.Errorf("недостаточно мест в зале")
	}

	costOfTickets := countOfTickets * (movie.ticketPrice * addedCostCoefficient)
	if c.balance-costOfTickets <= 0 {
		return fmt.Errorf("на вашем балансе недостаточно средств")
	}

	hall.movieMap[movie] += countOfTickets
	c.balance -= costOfTickets
	c.tickets[movie] = Tickets{
		movie:          movie,
		price:          movie.ticketPrice * addedCostCoefficient,
		countOfTickets: countOfTickets,
	}
	return nil
}

// Льготники покупают за 50% стоимости
func (b *Beneficiary) buyTicket(hallI HallI, movie Movie, countOfTickets int) error {
	hall, addedCostCoefficient := hallI.getCurrentHall()

	if hall.movieMap[movie]+countOfTickets > hall.capacity {
		return fmt.Errorf("недостаточно мест в зале")
	}

	costOfTickets := countOfTickets * (movie.ticketPrice / 2) * addedCostCoefficient
	if b.balance-costOfTickets <= 0 {
		return fmt.Errorf("на вашем балансе недостаточно средств")
	}

	hall.movieMap[movie] += countOfTickets
	b.balance -= costOfTickets
	b.tickets[movie] = Tickets{
		movie:          movie,
		price:          (movie.ticketPrice / 2) * addedCostCoefficient,
		countOfTickets: countOfTickets,
	}
	return nil

}

func main() {

	// создание фильмов
	movie1 := Movie{
		title:       "Игра в имитацию",
		startTime:   time.Now(),
		duration:    90,
		ticketPrice: 10,
	}

	movie2 := Movie{
		title: "Игра Престолов",
		// идет сразу за первым фильмом
		startTime:   movie1.startTime.Add(movie1.duration * time.Minute),
		duration:    120,
		ticketPrice: 20,
	}

	movie3 := Movie{
		title:       "Игра",
		startTime:   movie2.startTime.Add(movie2.duration * time.Minute),
		ticketPrice: 50,
	}

	movie4 := Movie{
		title:       "Поле чудес",
		startTime:   time.Now(),
		duration:    200,
		ticketPrice: 10,
	}

	movie5 := Movie{
		title:       "Кингзмэн",
		startTime:   movie4.startTime.Add(movie4.duration * time.Second),
		duration:    120,
		ticketPrice: 14,
	}

	movie6 := Movie{
		title:       "Соник в кино",
		startTime:   movie3.startTime.Add(movie3.duration * time.Second),
		duration:    150,
		ticketPrice: 15,
	}

	cinema := &Cinema{
		halls: []*Hall{
			{movieMap: map[Movie]int{movie1: 0, movie2: 0, movie3: 0}, capacity: 40},
			{movieMap: map[Movie]int{movie4: 0, movie5: 0}, capacity: 100},
			{movieMap: map[Movie]int{movie1: 0, movie3: 0, movie6: 0}, capacity: 20},
		},
		vipHalls: []*VipHall{
			{
				addedCostCoefficient: 2,
				Hall:                 &Hall{movieMap: map[Movie]int{movie2: 0, movie3: 0}, capacity: 10},
			},
			{
				addedCostCoefficient: 3,
				Hall:                 &Hall{movieMap: map[Movie]int{movie1: 0, movie6: 0}, capacity: 5},
			},
		},
	}

	// создание обычного клиента
	customer := &BasicCustomer{
		name:    "Иван Иванов",
		tickets: make(map[Movie]Tickets),
		balance: 3500,
	}

	// создание ветерана
	beneficiary := &Beneficiary{
		name:    "Игорь Петрович",
		tickets: make(map[Movie]Tickets),
		balance: 2000,
	}

	// покупка билета
	buyTicket(customer, cinema.halls[0], movie1, 10)
	buyTicket(customer, cinema.halls[1], movie2, 20)
	buyTicket(customer, cinema.halls[2], movie3, 15)
	buyTicket(customer, cinema.vipHalls[1], movie1, 10)

	buyTicket(beneficiary, cinema.halls[0], movie1, 25)
	buyTicket(beneficiary, cinema.vipHalls[0], movie1, 25)

	// вывод информации о билете и фильме
	fmt.Printf("Клиент: %s\n", customer.name)
	fmt.Printf("Баланс: %d\n", customer.balance)
	for movie, ticket := range customer.tickets {
		fmt.Printf("Билеты на фильм %s, цена : %d , количество купленных билетов : %d\n",
			movie.title, ticket.price, ticket.countOfTickets)
	}

	// вывод информации о билете и фильме
	fmt.Printf("Клиент: %s\n", beneficiary.name)
	fmt.Printf("Баланс: %d\n", beneficiary.balance)
	for movie, ticket := range customer.tickets {
		fmt.Printf("Билеты на фильм %s, цена : %d , количество купленных билетов : %d\n",
			movie.title, ticket.price, ticket.countOfTickets)
	}
}
