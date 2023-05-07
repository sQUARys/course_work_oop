package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// Вариант 14
// Имитационное моделирование расписания автобусов

// DONE 1. пауза между выездом автобус из парка
// DONE  2. сделать кусочно линейной пассажиров
// DONE 3. простаивание автобусов в парке

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
	name          string
	passengerFunc func(time.Time) int // Функция зависимости количества пассажиров от времени
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

func main() {
	busController := createBusController()
	busController.Simulate()
}

func createBusController() BusController {
	// test case(для того, чтобы не ждать часами работы программы)
	// rand.Seed(time.Now().UnixNano())
	// min := 8
	// max := 22

	passengersFunc := func(t time.Time) int {
		switch hour := t.Hour(); { //  test case : hour := rand.Intn(max-min+1) + min;
		case 8 <= hour && 11 >= hour:
			return 20 * hour
		case 18 <= hour && 20 >= hour:
			return 40 * hour
		default:
			return 5 * hour
		}
	}
	firstStop := &Stop{"First line", passengersFunc}
	secondStop := &Stop{"Second line", passengersFunc}
	thirdStop := &Stop{"Third line", passengersFunc}
	fourthStop := &Stop{"Fourth line", passengersFunc}
	fifthStop := &Stop{"Fifth line", passengersFunc}
	sixthStop := &Stop{"Sixth line", passengersFunc}
	seventhStop := &Stop{"Seventh line", passengersFunc}
	eightStop := &Stop{"Eight line", passengersFunc}
	ninthStop := &Stop{"Nine line", passengersFunc}

	countOfBussesNeededInTheRouteFunc := func(t time.Time) int {
		switch hour := t.Second(); { // change to hour
		case 8 <= hour && 11 >= hour:
			return 5 * hour
		case 18 <= hour && 20 >= hour:
			return 10 * hour
		default:
			return 2 * hour
		}
	}
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
		stops: []*Stop{secondStop, firstStop, thirdStop},
		wayTime: map[*Stop]int{
			secondStop: 40,
			firstStop:  20,
			thirdStop:  2,
		},
		countOfCurrentBusesInTheRoute: 0,
		timeGapBetweenBuses:           timeGapBetweenBuses,
		countOfBussesNeededInTheRoute: countOfBussesNeededInTheRouteFunc,
	}
	busRoute3 := &BusRoute{
		stops: []*Stop{fifthStop, sixthStop, seventhStop},
		wayTime: map[*Stop]int{
			fifthStop:   10,
			sixthStop:   50,
			seventhStop: 10,
		},
		countOfCurrentBusesInTheRoute: 0,
		timeGapBetweenBuses:           timeGapBetweenBuses,
		countOfBussesNeededInTheRoute: countOfBussesNeededInTheRouteFunc,
	}
	busRoute4 := &BusRoute{
		stops: []*Stop{fourthStop, ninthStop, thirdStop},
		wayTime: map[*Stop]int{
			fourthStop: 100,
			ninthStop:  20,
			thirdStop:  54,
		},
		countOfCurrentBusesInTheRoute: 0,
		timeGapBetweenBuses:           timeGapBetweenBuses,
		countOfBussesNeededInTheRoute: countOfBussesNeededInTheRouteFunc,
	}
	busRoute5 := &BusRoute{
		stops: []*Stop{eightStop, fifthStop, secondStop},
		wayTime: map[*Stop]int{
			secondStop: 2,
			eightStop:  60,
			fifthStop:  20,
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

	b.route.countOfCurrentBusesInTheRoute++
	logStop(f, time.Now(), fmt.Sprintf("Автобус %d вышел на маршрут №%d\n", b.number, b.routeNumber))

	for len(b.route.stops)-1 > b.positionIndex {
		b.goToNextStop()
		logStop(f, time.Now(), fmt.Sprintf("Автобус %d на остановке %s\n", b.number, b.position.name))

		passengers := b.position.passengerFunc(time.Now())
		if passengers > 0 {
			tackedPassengersCount := 0
			for i := 0; i < passengers; i++ {
				if b.capacity >= b.passengers+1 {
					b.passengers++
					tackedPassengersCount++
				}
			}
			if tackedPassengersCount != 0 {
				logStop(f, time.Now(), fmt.Sprintf("Автобус %d взял на остановке %d пассажиров, осталось вместительности %d\n\n", b.number, tackedPassengersCount, b.capacity-b.passengers))
			}
		}

		time.Sleep(time.Duration(b.route.wayTime[b.position]) * multipliedTimeForWaysBetweenStops)
	}

	logStop(f, time.Now(), fmt.Sprintf("Автобус %d завершил маршрут, и идет в обратную сторону\n", b.number))

	for b.positionIndex-1 >= 0 {
		b.goToPrevStop()
		logStop(f, time.Now(), fmt.Sprintf("Автобус %d на остановке %s\n", b.number, b.position.name))

		passengers := b.position.passengerFunc(time.Now())
		if passengers > 0 {
			tackedPassengersCount := 0
			for i := 0; i < passengers; i++ {
				if b.capacity >= b.passengers+1 {
					b.passengers++
					tackedPassengersCount++
				}
			}
			if tackedPassengersCount != 0 {
				logStop(f, time.Now(), fmt.Sprintf("Автобус %d взял на остановке %d пассажиров, осталось вместительности %d\n\n", b.number, tackedPassengersCount, b.capacity-b.passengers))
			}
		}
		time.Sleep(time.Duration(b.route.wayTime[b.position]) * multipliedTimeForWaysBetweenStops)
	}

	b.passengers = 0

	logStop(f, time.Now(), fmt.Sprintf("Автобус %d завершил маршрут, высадил всех пассажиров, и идет в изначальную сторону.\n\n\n\n", b.number))
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
