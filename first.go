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
	adress   string
	halls    []*Hall
	vipHalls []*VipHall
}

type Hall struct {
	// key - movieMap, value - countOfBoughtTickets
	number   int
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
	// return hall, value of added cost coefficient, hallNumber and isHallVip
	getCurrentHall() (*Hall, int, int, bool)
	getAddedCostCoefficient() int
}

type Movie struct {
	title       string
	startTime   time.Time
	duration    time.Duration
	ticketPrice int
}

type Tickets struct {
	isVipHall      bool
	hallNumber     int
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
	printTicketsOfCurrentHall(hallNumber int)
	buyTicket(hall HallI, movie Movie, countOfTickets int) error
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
		adress: "ул. Пушкина",
		halls: []*Hall{
			{number: 1, movieMap: map[Movie]int{movie1: 0, movie2: 0, movie3: 0}, capacity: 40},
			{number: 2, movieMap: map[Movie]int{movie4: 0, movie5: 0}, capacity: 100},
			{number: 3, movieMap: map[Movie]int{movie1: 0, movie3: 0, movie6: 0}, capacity: 20},
		},
		vipHalls: []*VipHall{
			{
				addedCostCoefficient: 2,
				Hall:                 &Hall{number: 1, movieMap: map[Movie]int{movie2: 0, movie3: 0}, capacity: 10},
			},
			{
				addedCostCoefficient: 3,
				Hall:                 &Hall{number: 2, movieMap: map[Movie]int{movie1: 0, movie6: 0}, capacity: 5},
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

	fmt.Printf("В кинотеатре на улице %s были оформлены следующие билеты на фильм\n", cinema.adress)
	fmt.Println("\n...Обычные залы...")
	for _, hall := range cinema.halls {
		customer.printTicketsOfCurrentHall(hall.number)
		beneficiary.printTicketsOfCurrentHall(hall.number)
	}
	fmt.Println("\n...VIP залы...")
	for _, hall := range cinema.halls {
		customer.printTicketsOfCurrentHall(hall.number)
		beneficiary.printTicketsOfCurrentHall(hall.number)
	}
}

func (c *BasicCustomer) printTicketsOfCurrentHall(hallNumber int) {
	fmt.Printf("Зал №%d:\n", hallNumber)
	for movie, tickets := range c.tickets {
		if tickets.hallNumber == hallNumber {
			fmt.Printf("--%s приобрел %d билетов на фильм \"%s\" на время %s",
				c.name, c.tickets[movie].countOfTickets, movie.title, movie.startTime.UTC())
		}
	}
	fmt.Println()
}

func (b *Beneficiary) printTicketsOfCurrentHall(hallNumber int) {
	for movie, tickets := range b.tickets {
		if tickets.hallNumber == hallNumber {
			fmt.Printf("Зал №%d:\n", hallNumber)
			fmt.Printf("--%s приобрел %d билетов на фильм \"%s\" на время %s",
				b.name, b.tickets[movie].countOfTickets, movie.title, movie.startTime.UTC())
		}
	}
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

func (hall *Hall) getCurrentHall() (*Hall, int, int, bool) {
	return hall, getAddedCostCoefficient(hall), hall.number, false
}

func (vipHall *VipHall) getCurrentHall() (*Hall, int, int, bool) {
	return vipHall.Hall, getAddedCostCoefficient(vipHall), vipHall.number, true
}

func buyTicket(customerI CustomerI, hall HallI, movie Movie, countOfTickets int) {
	err := customerI.buyTicket(hall, movie, countOfTickets)
	if err != nil {
		log.Println(err)
	}
}

// Обычный покупатель покупает за полную стоимость
func (c *BasicCustomer) buyTicket(hallI HallI, movie Movie, countOfTickets int) error {
	hall, addedCostCoefficient, hallNumber, isVipHall := hallI.getCurrentHall()

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
		hallNumber:     hallNumber,
		isVipHall:      isVipHall,
		movie:          movie,
		price:          movie.ticketPrice * addedCostCoefficient,
		countOfTickets: countOfTickets,
	}
	return nil
}

// Льготники покупают за 50% стоимости
func (b *Beneficiary) buyTicket(hallI HallI, movie Movie, countOfTickets int) error {
	hall, addedCostCoefficient, hallNumber, isVipHall := hallI.getCurrentHall()

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
		hallNumber:     hallNumber,
		isVipHall:      isVipHall,
		movie:          movie,
		price:          (movie.ticketPrice / 2) * addedCostCoefficient,
		countOfTickets: countOfTickets,
	}
	return nil

}
