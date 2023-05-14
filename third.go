package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"sync"
	"time"
)

// Вариант 14
// Имитационное моделирование расписания автобусов

// Добавить коэффициент к остановкам что каждая остановка

// DONE : добавить пассажирам время ожидания сортировать
//  по времени и брать их с остановки и логировать их

// Контроллер управления движением автобусов
type BusController struct {
	buses []*Bus // Список автобусов на маршруте
}

// Маршрут автобуса, который проходит через несколько остановок
type BusRoute struct {
	stops                         []*Stop                  // Остановки на маршруте
	wayTime                       map[*Stop]int            // Время между данной и следующей остановкой
	timeGapBetweenBuses           time.Duration            // Промежуток времени между двумя автобусами при старте
	countOfBussesNeededInTheRoute func(time time.Time) int // Нужное количество автобусов на маршруте в данное время
	countOfCurrentBusesInTheRoute int
}

// Остановка на маршруте с заданными координатами и функцией зависимости количества пассажиров от времени
type Stop struct {
	*sync.Mutex
	name              string
	passengerCount    int
	waitingPassengers []Passenger
}

type Passenger struct {
	startTimeOfWaiting time.Time
}

// Автобус на маршруте с заданным номером и маршрутом
type Bus struct {
	number        int       // Номер автобуса
	routeNumber   int       // Номер маршрута
	route         *BusRoute // Маршрут автобуса
	position      *Stop     // Текущая остановка
	positionIndex int       // Индекс текущей остановки относительно маршрута
	passengers    int       // Количество пассажиров в автобусе
	capacity      int       // Вместимость пассажиров автобусом
}

const (
	// Когда автобус сделал один маршрут туда обратно, у него есть перерыв
	breakAfterTwoSideRoute = time.Millisecond
	// Промежуток между выходом автобусов из парка
	timeGapBetweenBuses = 2 * time.Millisecond
	// Количество раз, сколько автобус сделает маршрутов туда обратно, перед тем как закончить
	countOfTwoSideWaysOfBus           = 5
	multipliedTimeForWaysBetweenStops = time.Millisecond
)

// functions
var (
	countOfBussesNeededInTheRouteFunc = func(t time.Time) int {
		switch hour := t.Second(); { // change to hour
		case 8 <= hour && 11 >= hour:
			return 5 * hour
		case 18 <= hour && 20 >= hour:
			return 10 * hour
		default:
			return 2 * hour
		}
	}

	// test case(для того, чтобы не ждать часами работы программы)
	passengersGenerateFunc = func(prevPassengers int, t time.Time) int {
		if prevPassengers == 0 {
			return rand.Intn(max-min+1) + min
		}

		var count int
		switch hour := rand.Intn(max-min+1) + min; { // t.Hour(); { //  test case : hour := rand.Intn(max-min+1) + min;
		case 8 <= hour && 11 >= hour:
			count = prevPassengers / 6
		case 18 <= hour && 20 >= hour:
			count = prevPassengers / 7
		default:
			count = prevPassengers / 5
		}
		return prevPassengers + count
	}

	passengersOutFunc = func(prevPassengers int, t time.Time) int {
		if prevPassengers == 0 {
			return rand.Intn(max-min+1) + min
		}
		var count int
		switch hour := rand.Intn(max-min+1) + min; { // t.Hour(); { //  test case : hour := rand.Intn(max-min+1) + min;
		case 8 <= hour && 11 >= hour:
			count = prevPassengers / 8 // 10% выходят из автобуса утром
		case 18 <= hour && 20 >= hour:
			count = prevPassengers / 7
		default:
			count = prevPassengers / 5
		}
		return prevPassengers - count
	}
)

var max int
var min int

func main() {
	rand.Seed(time.Now().UnixNano())
	min = 10
	max = 22

	busController := createBusController()
	busController.Simulate()
}

func createBusController() BusController {
	var mux sync.Mutex
	firstStop := &Stop{&mux, "First line", 0, []Passenger{}}
	secondStop := &Stop{&mux, "Second line", 0, []Passenger{}}
	thirdStop := &Stop{&mux, "Third line", 0, []Passenger{}}
	fourthStop := &Stop{&mux, "Fourth line", 0, []Passenger{}}
	fifthStop := &Stop{&mux, "Fifth line", 0, []Passenger{}}
	sixthStop := &Stop{&mux, "Sixth line", 0, []Passenger{}}
	seventhStop := &Stop{&mux, "Seventh line", 0, []Passenger{}}
	eightStop := &Stop{&mux, "Eight line", 0, []Passenger{}}
	ninthStop := &Stop{&mux, "Nine line", 0, []Passenger{}}
	tenStop := &Stop{&mux, "Ten line", 0, []Passenger{}}

	busRoute1 := &BusRoute{
		stops: []*Stop{firstStop, secondStop, thirdStop},
		wayTime: map[*Stop]int{
			firstStop:  10,
			secondStop: 5,
			thirdStop:  2,
		},
		timeGapBetweenBuses:           timeGapBetweenBuses,
		countOfCurrentBusesInTheRoute: 0,
		countOfBussesNeededInTheRoute: countOfBussesNeededInTheRouteFunc,
	}
	busRoute2 := &BusRoute{
		stops: []*Stop{firstStop, secondStop, firstStop, thirdStop},
		wayTime: map[*Stop]int{
			firstStop:  10,
			secondStop: 40,
			firstStop:  20,
			thirdStop:  2,
		},
		countOfCurrentBusesInTheRoute: 0,
		timeGapBetweenBuses:           timeGapBetweenBuses,
		countOfBussesNeededInTheRoute: countOfBussesNeededInTheRouteFunc,
	}
	busRoute3 := &BusRoute{
		stops: []*Stop{firstStop, fifthStop, sixthStop, seventhStop},
		wayTime: map[*Stop]int{
			firstStop:   10,
			fifthStop:   10,
			sixthStop:   50,
			seventhStop: 10,
		},
		countOfCurrentBusesInTheRoute: 0,
		timeGapBetweenBuses:           timeGapBetweenBuses,
		countOfBussesNeededInTheRoute: countOfBussesNeededInTheRouteFunc,
	}
	busRoute4 := &BusRoute{
		stops: []*Stop{firstStop, fourthStop, ninthStop, thirdStop},
		wayTime: map[*Stop]int{
			firstStop:  10,
			fourthStop: 100,
			ninthStop:  20,
			thirdStop:  54,
		},
		countOfCurrentBusesInTheRoute: 0,
		timeGapBetweenBuses:           timeGapBetweenBuses,
		countOfBussesNeededInTheRoute: countOfBussesNeededInTheRouteFunc,
	}
	busRoute5 := &BusRoute{
		stops: []*Stop{firstStop, eightStop, fifthStop, secondStop, tenStop},
		wayTime: map[*Stop]int{
			firstStop:  10,
			secondStop: 2,
			eightStop:  60,
			fifthStop:  20,
			tenStop:    10,
		},
		countOfCurrentBusesInTheRoute: 0,
		timeGapBetweenBuses:           timeGapBetweenBuses,
		countOfBussesNeededInTheRoute: countOfBussesNeededInTheRouteFunc,
	}

	busController := BusController{
		buses: []*Bus{
			{1, 1, busRoute1, firstStop, 0, 0, 100},
			{2, 1, busRoute1, secondStop, 0, 0, 200},
			{3, 2, busRoute2, secondStop, 0, 0, 300},
			{4, 2, busRoute2, secondStop, 0, 0, 500},
			{5, 2, busRoute2, secondStop, 0, 0, 0},
			{6, 3, busRoute3, fifthStop, 0, 0, 300},
			{7, 4, busRoute4, fourthStop, 0, 0, 200},
			{8, 5, busRoute5, eightStop, 0, 0, 1000},
		},
	}
	return busController
}

// Имитационное моделирование движения автобусов на маршруте
func (c *BusController) Simulate() {
	var wg sync.WaitGroup

	for _, bus := range c.buses {
		wg.Add(1)
		go bus.startWorking(&wg)
	}

	wg.Wait()
}

func (b *Bus) startWorking(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < countOfTwoSideWaysOfBus; i++ {
		currentTime := time.Now()
		if b.route.countOfCurrentBusesInTheRoute+1 <= b.route.countOfBussesNeededInTheRoute(currentTime) {
			b.route.countOfCurrentBusesInTheRoute++
			b.startTwoSideRoute()
			time.Sleep(breakAfterTwoSideRoute)
		}
	}
}

func (b *Bus) startTwoSideRoute() {

	f, err := os.OpenFile(fmt.Sprintf("route%d.txt", b.routeNumber), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	// автобусы выходят из парка на маршрут с периодичностью
	if b.route.countOfCurrentBusesInTheRoute != 0 {
		time.Sleep(b.route.timeGapBetweenBuses)
	}

	logStop(f, time.Now(), fmt.Sprintf("Автобус %d вышел на маршрут №%d\n", b.number, b.routeNumber))

	for len(b.route.stops)-1 > b.positionIndex {
		b.goToNextStop()
		logStop(f, time.Now(), fmt.Sprintf("Автобус %d на остановке %s\n", b.number, b.position.name))

		// часть вышла
		b.passengers = passengersOutFunc(b.passengers, time.Now())

		b.position.Lock()

		// генерируем еще пассажиров к уже стоящим
		passengers := passengersGenerateFunc(b.position.passengerCount, time.Now())
		for i := 0; i < passengers; i++ {
			b.position.waitingPassengers = append(b.position.waitingPassengers, Passenger{time.Now()})
		}

		// сортируем по времени ожидания в порядке: от самого дальнего времени к самому ближнему к настоящему
		sort.Slice(b.position.waitingPassengers, func(i, j int) bool {
			firstTime := b.position.waitingPassengers[i].startTimeOfWaiting
			secondTime := b.position.waitingPassengers[j].startTimeOfWaiting
			return firstTime.Before(secondTime)
		})

		b.position.passengerCount = passengers

		// для логирования создадим переменную которая накопит время ожидания пассажиров
		waitingTimeOfEachPassenger := ""

		if passengers > 0 {
			tackedPassengersCount := 0
			for i := 0; i < passengers; i++ {
				if b.capacity >= b.passengers+1 {
					b.passengers++
					tackedPassengersCount++
					b.position.passengerCount--

					waitingTimeOfCurrentPassenger := time.Now().Sub(b.position.waitingPassengers[0].startTimeOfWaiting)
					waitingTimeOfEachPassenger += fmt.Sprintf("Пассажир №%d : %s; ", i+1, waitingTimeOfCurrentPassenger)
					b.position.waitingPassengers = b.position.waitingPassengers[1:]
				}
			}
			if tackedPassengersCount != 0 {
				logStop(f, time.Now(), fmt.Sprintf("Автобус %d взял на остановке %d пассажиров, осталось вместительности %d\n", b.number, tackedPassengersCount, b.capacity-b.passengers))
				logStop(f, time.Now(), fmt.Sprintf("Время ожидания взятых автобусом %d пассажиров: %s \n\n", b.number, waitingTimeOfEachPassenger))
			}
		}

		b.position.Unlock()

		time.Sleep(time.Duration(b.route.wayTime[b.position]) * multipliedTimeForWaysBetweenStops)
	}

	logStop(f, time.Now(), fmt.Sprintf("Автобус %d завершил маршрут, и идет в обратную сторону\n", b.number))

	for b.positionIndex-1 >= 0 {
		b.goToPrevStop()
		logStop(f, time.Now(), fmt.Sprintf("Автобус %d на остановке %s\n", b.number, b.position.name))

		// часть вышла
		b.passengers = passengersOutFunc(b.passengers, time.Now())

		b.position.Lock()

		// генерируем еще пассажиров к уже стоящим
		passengers := passengersGenerateFunc(b.position.passengerCount, time.Now())
		for i := 0; i < passengers; i++ {
			b.position.waitingPassengers = append(b.position.waitingPassengers, Passenger{time.Now()})
		}

		// сортируем по времени ожидания в порядке: от самого дальнего времени к самому ближнему к настоящему
		sort.Slice(b.position.waitingPassengers, func(i, j int) bool {
			firstTime := b.position.waitingPassengers[i].startTimeOfWaiting
			secondTime := b.position.waitingPassengers[j].startTimeOfWaiting
			return firstTime.Before(secondTime)
		})

		b.position.passengerCount = passengers

		// для логирования создадим переменную которая накопит время ожидания пассажиров
		waitingTimeOfEachPassenger := ""

		if passengers > 0 {
			tackedPassengersCount := 0
			for i := 0; i < passengers; i++ {
				if b.capacity >= b.passengers+1 {
					b.passengers++
					tackedPassengersCount++
					b.position.passengerCount--

					waitingTimeOfCurrentPassenger := time.Now().Sub(b.position.waitingPassengers[0].startTimeOfWaiting)
					waitingTimeOfEachPassenger += fmt.Sprintf("Пассажир №%d : %s; ", i+1, waitingTimeOfCurrentPassenger)
					b.position.waitingPassengers = b.position.waitingPassengers[1:]
				}
			}
			if tackedPassengersCount != 0 {
				logStop(f, time.Now(), fmt.Sprintf("Автобус %d взял на остановке %d пассажиров, осталось вместительности %d\n", b.number, tackedPassengersCount, b.capacity-b.passengers))
				logStop(f, time.Now(), fmt.Sprintf("Время ожидания взятых автобусом %d пассажиров: %s \n\n", b.number, waitingTimeOfEachPassenger))
			}
		}

		b.position.Unlock()

		time.Sleep(time.Duration(b.route.wayTime[b.position]) * multipliedTimeForWaysBetweenStops)
	}

	b.passengers = 0

	logStop(f, time.Now(), fmt.Sprintf("Автобус %d завершил маршрут, высадил всех пассажиров, и ждет следующего рейса в парке.\n\n\n\n", b.number))
}

func logStop(f *os.File, time time.Time, log string) {
	if log == "" {
		return
	}

	_, err := f.WriteString(time.UTC().String() + ": " + log)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (b *Bus) goToNextStop() {
	b.positionIndex++
	b.position = b.route.stops[b.positionIndex]
}

func (b *Bus) goToPrevStop() {
	b.positionIndex--
	b.position = b.route.stops[b.positionIndex]
}
