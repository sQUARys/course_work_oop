package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// Вариант 14
// Имитационное моделирование расписания автобусов

// TODO: пауза между выездом автобус из парка, сделать кусочно линейной пассажиров, простаивание автобусов в парке
// если такое не нужно

// Контроллер управления движением автобусов
type BusController struct {
	buses []*Bus // Список автобусов на маршруте
}

// Маршрут автобуса, который проходит через несколько остановок
type BusRoute struct {
	stops   []*Stop       // Остановки на маршруте
	wayTime map[*Stop]int // Расстояние между данной и следующей остановкой
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
	positionIndex int       // Индекс остановки относительно маршрута
	passengers    int       // Количество пассажиров в автобусе
	capacity      int       // Вместимость пассажиров автобусом
}

func main() {
	busController := createBusController()
	busController.Simulate()
}

func createBusController() BusController {
	firstStop := &Stop{"First line", func(t time.Time) int {
		return 10 * (t.Second() % 10) // количество пассажиров
	}}
	secondStop := &Stop{"Second line", func(t time.Time) int {
		return 3 * (t.Second() % 10) // количество пассажиров
	}}
	thirdStop := &Stop{"Third line", func(t time.Time) int {
		return 2 * (t.Second() % 10) // количество пассажиров
	}}
	fourthStop := &Stop{"Fourth line", func(t time.Time) int {
		return 3 * (t.Second() % 10) // количество пассажиров
	}}
	fifthStop := &Stop{"Fifth line", func(t time.Time) int {
		return 5 * (t.Second() % 10) // количество пассажиров
	}}
	sixthStop := &Stop{"Sixth line", func(t time.Time) int {
		return 7 * (t.Second() % 10) // количество пассажиров
	}}
	seventhStop := &Stop{"Seventh line", func(t time.Time) int {
		return 15 * (t.Second() % 10) // количество пассажиров
	}}
	eightStop := &Stop{"Eight line", func(t time.Time) int {
		return 11 * (t.Second() % 10) // количество пассажиров
	}}
	ninthStop := &Stop{"Nine line", func(t time.Time) int {
		return 2 * (t.Second() % 10) // количество пассажиров
	}}

	busRoute1 := &BusRoute{
		stops: []*Stop{firstStop, secondStop, thirdStop},
		wayTime: map[*Stop]int{
			firstStop:  10,
			secondStop: 5,
			thirdStop:  2,
		},
	}

	busRoute2 := &BusRoute{
		stops: []*Stop{secondStop, firstStop, thirdStop},
		wayTime: map[*Stop]int{
			secondStop: 40,
			firstStop:  20,
			thirdStop:  2,
		},
	}

	busRoute3 := &BusRoute{
		stops: []*Stop{fifthStop, sixthStop, seventhStop},
		wayTime: map[*Stop]int{
			fifthStop:   10,
			sixthStop:   50,
			seventhStop: 10,
		},
	}

	busRoute4 := &BusRoute{
		stops: []*Stop{fourthStop, ninthStop, thirdStop},
		wayTime: map[*Stop]int{
			fourthStop: 100,
			ninthStop:  20,
			thirdStop:  54,
		},
	}

	busRoute5 := &BusRoute{
		stops: []*Stop{eightStop, fifthStop, secondStop},
		wayTime: map[*Stop]int{
			secondStop: 2,
			eightStop:  60,
			fifthStop:  20,
		},
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
		go bus.startTwoSideRoute(&wg)
	}

	wg.Wait()
}

func (b *Bus) startTwoSideRoute(wg *sync.WaitGroup) {
	defer wg.Done()

	f, err := os.OpenFile(fmt.Sprintf("route%d.txt", b.routeNumber), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	for len(b.route.stops)-1 > b.positionIndex {
		b.goToNextStop()
		logStop(f, fmt.Sprintf("Автобус %d на остановке %s\n", b.number, b.position.name))

		logStop(f, b.takePassengers())

		time.Sleep(time.Duration(b.route.wayTime[b.position]) * time.Millisecond)
	}

	logStop(f, fmt.Sprintf("Автобус %d завершил маршрут, и идет в обратную сторону\n", b.number))

	for b.positionIndex-1 >= 0 {
		b.goToPrevStop()
		logStop(f, fmt.Sprintf("Автобус %d на остановке %s\n", b.number, b.position.name))

		logStop(f, b.takePassengers())
		time.Sleep(time.Duration(b.route.wayTime[b.position]) * time.Millisecond)
	}

	b.capacity += b.passengers
	b.passengers = 0

	logStop(f, fmt.Sprintf("Автобус %d завершил маршрут, высадил всех пассажиров, и идет в изначальную сторону.\n\n\n\n", b.number))

}

func logStop(f *os.File, log string) {
	t := time.Now().Format(time.RFC3339Nano)
	_, err := f.WriteString(t + ": " + log)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (b *Bus) takePassengers() string {
	// Если на остановке есть пассажиры, то добавляем их в автобус
	passengers := b.position.passengerFunc(time.Now())
	if passengers > 0 && b.capacity-passengers >= 0 {
		b.capacity -= passengers
		b.passengers += passengers
		logs := fmt.Sprintf("Автобус %d взял на остановке %d пассажиров, осталось вместительности %d\n\n", b.number, passengers, b.capacity)
		return logs
	}
	return ""
}

func (b *Bus) goToNextStop() {
	b.positionIndex++
	b.position = b.route.stops[b.positionIndex]
}

func (b *Bus) goToPrevStop() {
	b.positionIndex--
	b.position = b.route.stops[b.positionIndex]
}
