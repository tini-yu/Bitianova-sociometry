package service

import (
	"cmp"
	"fmt"
	"log"
	"math"
	"slices"
	"strconv"
)

// Всего людей в тесте (длина списка людей)
func getPeopleTotal() int {
	if len(names) == 0 {
		log.Print("Список имён не должен быть пустым")
		return 1
	}
	return len(names)
}

// Всего будет 7 пунктов, которые необходимо отразить в отчёте. Для каждого численного значения свой подпункт.

type n1struct struct {
	n1_1 string // date
	n1_2 int    // кол-во участников диагностики
	// n1_3 int //кол-во человек в классе
}

var n1 n1struct

// В данный момент мы считаем результаты на основе, что все ученики прошли диагностику (Разделение n1_2/n1_3 - на перспективу)

func (n1 *n1struct) filln1() {
	n1.n1_1 = prettyTestDate
	n1.n1_2 = getPeopleTotal()
}

type n2struct struct {
	n2_1_1   float64 // Процент
	n2_1_1_p int     // Кол-во человек
	n2_1_2   float64
	n2_1_2_p int
	n2_1_3   float64
	n2_1_3_p int
	n2_1_4   float64
	n2_1_4_p int
	n2_2     float64
	n2_2_p   int
	n2_3     float64
	n2_3_p   int
	n2_4     float64
	n2_4_p   int
}

var n2 n2struct

func (n2 *n2struct) filln2(results [][]string) {
	acceptedPeople := 0
	notAcceptedPeople := 0
	stars := 0
	popular := 0
	undefined := 0
	rejected := 0
	isolated := 0
	focusGroupRejected := 0

	total := getPeopleTotal()

	for _, row := range results {
		switch row[2] {
		case "Звезда":
			acceptedPeople += 1
			stars += 1
		case "Популярные":
			acceptedPeople += 1
			popular += 1
		case "Принятые": // До 20.09.2026 был "Не определено"
			acceptedPeople += 1
			undefined += 1
		case "Отверженные":
			notAcceptedPeople += 1
			rejected += 1
			if row[8] != "-" {
				focusGroupRejected += 1
			}
		case "Изолированные":
			notAcceptedPeople += 1
			isolated += 1
			if row[8] != "-" {
				focusGroupRejected += 1
			}
		default:
			fmt.Println("Неизвестный статус: ", row[2])
		}
	}

	n2.n2_1_1 = math.RoundToEven(float64(acceptedPeople)*1000/float64(total)) / 10
	n2.n2_1_1_p = acceptedPeople
	n2.n2_1_4 = math.RoundToEven(float64(notAcceptedPeople)*1000/float64(total)) / 10
	n2.n2_1_4_p = notAcceptedPeople
	if acceptedPeople != 0 {
		n2.n2_1_2 = math.RoundToEven(float64(stars)*1000/float64(acceptedPeople)) / 10
		n2.n2_1_2_p = stars
		n2.n2_1_3 = math.RoundToEven(float64(popular)*1000/float64(acceptedPeople)) / 10
		n2.n2_1_3_p = popular
	} else {
		n2.n2_1_2 = 0
		n2.n2_1_2_p = 0
		n2.n2_1_3 = 0
		n2.n2_1_3_p = 0
	}
	if notAcceptedPeople != 0 {
		n2.n2_2 = math.RoundToEven(float64(rejected)*1000/float64(notAcceptedPeople)) / 10
		n2.n2_2_p = rejected
		n2.n2_3 = math.RoundToEven(float64(isolated)*1000/float64(notAcceptedPeople)) / 10
		n2.n2_3_p = isolated
		n2.n2_4 = math.RoundToEven(float64(focusGroupRejected)*1000/float64(notAcceptedPeople)) / 10
		n2.n2_4_p = focusGroupRejected
	} else {
		n2.n2_2 = 0
		n2.n2_2_p = 0
		n2.n2_3 = 0
		n2.n2_3_p = 0
		n2.n2_4 = 0
		n2.n2_4_p = 0
	}
}

type n3struct struct {
	n3_1   float64
	n3_1_p int
	n3_2   float64
	n3_2_p int
	n3_3   float64
	n3_3_p int
}

var n3 n3struct

func (n3 *n3struct) filln3(results [][]string) {
	total := getPeopleTotal()

	mutualPeopleTotal := 0
	mutualPositiveTotal := 0
	mutualNegativeTotal := 0

	for _, row := range results {
		positive, err := strconv.Atoi(row[3])
		if err != nil {
			fmt.Println(err)
		}
		negative, err := strconv.Atoi(row[4])
		if err != nil {
			fmt.Println(err)
		}

		if positive >= 3 || negative >= 3 {
			mutualPeopleTotal += 1
			if positive >= 3 {
				mutualPositiveTotal += 1
			}
			if negative >= 3 {
				mutualNegativeTotal += 1
			}
		}

		n3.n3_1 = math.RoundToEven(float64(mutualPeopleTotal)*1000/float64(total)) / 10
		n3.n3_1_p = mutualPeopleTotal
		if mutualPeopleTotal != 0 {
			n3.n3_2 = math.RoundToEven(float64(mutualPositiveTotal)*1000/float64(mutualPeopleTotal)) / 10
			n3.n3_2_p = mutualPositiveTotal
			n3.n3_3 = math.RoundToEven(float64(mutualNegativeTotal)*1000/float64(mutualPeopleTotal)) / 10
			n3.n3_3_p = mutualNegativeTotal
		} else {
			n3.n3_2 = 0
			n3.n3_2_p = 0
			n3.n3_3 = 0
			n3.n3_3_p = 0
		}

	}
}

type n4struct struct {
	n4_1   float64
	n4_1_p int
	n4_2   float64
	n4_2_p int
}

var n4 n4struct

func (n4 *n4struct) filln4(results [][]string, matrix1 [][]int) {
	total := getPeopleTotal()
	contradictory := 0
	contradictoryPositive := 0

	for row_id, row := range results {
		choice, err := strconv.Atoi(row[5])
		if err != nil {
			fmt.Println(err)
		}
		if choice > 0 {
			contradictory += 1
			pos := 0
			neg := 0
			for _, val := range matrix1[row_id] {
				if val == 1 {
					pos += 1
				} else if val == 2 {
					neg += 1
				}
			}
			if pos > neg {
				contradictoryPositive += 1
			}
		}
	}

	n4.n4_1 = math.RoundToEven(float64(contradictory)*1000/float64(total)) / 10
	n4.n4_1_p = contradictory
	if contradictory != 0 {
		n4.n4_2 = math.RoundToEven(float64(contradictoryPositive)*1000/float64(contradictory)) / 10
		n4.n4_2_p = contradictoryPositive
	} else {
		n4.n4_2 = 0
		n4.n4_2_p = 0
	}
}

type n56struct struct {
	n5_1   float64
	n5_1_p int
	n5_2   float64
	n5_2_p int
	n6_1   float64
	n6_1_p int
	n6_2   float64
	n6_2_p int
	n6_3   float64
	n6_3_p int
	n6_4   float64
	n6_4_p int
}

var n56 n56struct

func (n56 *n56struct) filln56(results [][]string) {
	total := getPeopleTotal()
	adeq := 0
	adeqRejected := 0
	unadeq := 0
	unadeqFocus := 0
	unadeqRejected := 0
	unadeqAccepted := 0

	for _, row := range results {
		if row[6] == "Адекватная" {
			adeq += 1
			if row[2] == "Отверженные" || row[2] == "Изолированные" {
				adeqRejected += 1
			}
		} else if row[6] == "Неадекватная" {
			unadeq += 1
			if row[2] == "Отверженные" || row[2] == "Изолированные" {
				if row[8] != "-" {
					unadeqFocus += 1
				} else {
					unadeqRejected += 1
				}
			} else {
				unadeqAccepted += 1
			}
		}
	}

	n56.n5_1 = math.RoundToEven(float64(adeq)*1000/float64(total)) / 10
	n56.n5_1_p = adeq
	if adeq != 0 {
		n56.n5_2 = math.RoundToEven(float64(adeqRejected)*1000/float64(adeq)) / 10
		n56.n5_2_p = adeqRejected
	} else {
		n56.n5_2 = 0
	}
	n56.n6_1 = math.RoundToEven(float64(unadeq)*1000/float64(total)) / 10
	n56.n6_1_p = unadeq
	if unadeq != 0 {
		n56.n6_2 = math.RoundToEven(float64(unadeqFocus)*1000/float64(unadeq)) / 10
		n56.n6_2_p = unadeqFocus
		n56.n6_3 = math.RoundToEven(float64(unadeqRejected)*1000/float64(unadeq)) / 10
		n56.n6_3_p = unadeqRejected
		n56.n6_4 = math.RoundToEven(float64(unadeqAccepted)*1000/float64(unadeq)) / 10
		n56.n6_4_p = unadeqAccepted
	} else {
		n56.n6_2 = 0
		n56.n6_2_p = 0
		n56.n6_3 = 0
		n56.n6_3_p = 0
		n56.n6_4 = 0
		n56.n6_4_p = 0
	}
}

type n8support struct {
	n8_n7_total       int      // всего людей в п.6 аналит. записки
	n8_n7_noMatch     int      // кол-во "Не совпадает"
	n8_n7_match       int      // кол-во "Совпадает"
	n8_n7_match_names []string // имена людей с "Совпадает"
}

type n7struct struct {
	n7 []string
}

var n8s n8support

var n7 n7struct

func (n7 *n7struct) filln7(results [][]string) {
	total := getPeopleTotal()
	refNames := []string{}
	refResults := []string{}

	// for support:
	match := 0
	noMatch := 0
	matchNames := []string{}

	id_refN := [][]int{} // [[id, refN]] - где id - индекс имени человека в refNames.

	// Для сортировки по кол-ву реф. выборов заведем третий слайс с индексами|выборами [][]int
	// Отсортируем по кол-ву выборов, а потом используем индексы, чтобы составить список6 Имя - результат

	for row_id, row := range results {
		ref, err := strconv.Atoi(row[7])
		if err != nil {
			fmt.Println(err)
		}
		if float64(ref) > (float64(total) * 0.25) {
			refNames = append(refNames, names[row_id])
			if row[2] == "Звезда" || row[2] == "Популярные" {
				refResults = append(refResults, "Совпадает")
				match += 1
				matchNames = append(matchNames, names[row_id])
			} else {
				refResults = append(refResults, "Не совпадает")
				noMatch += 1
			}
			// Пополняем список для сорта
			refNames_id := len(refNames) - 1
			id_refN = append(id_refN, []int{refNames_id, ref})
		}
	}

	n8s.n8_n7_match = match
	n8s.n8_n7_noMatch = noMatch
	n8s.n8_n7_match_names = matchNames

	// Сортируем
	slices.SortFunc(id_refN, func(a, b []int) int {
		return cmp.Compare(a[1], b[1])
	})
	slices.Reverse(id_refN)

	parts := make([]string, 0, len(refNames))
	for _, id_ref := range id_refN {
		parts = append(parts, refNames[id_ref[0]]+" - "+refResults[id_ref[0]])
	}

	n7.n7 = parts
	n8s.n8_n7_total = len(refNames)
}

func (s *ReportService) CalculateAnalyticalReport() error {
	results, err := getResultsAsStringSlices()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Ошибка получения результатов для справки: %s", err)
	}
	matrix1, err := getMatrix("matrix1.json")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Ошибка подгрузки матрицы: %s", err)
	}

	n1.filln1()
	n2.filln2(results)
	n3.filln3(results)
	n4.filln4(results, matrix1)
	n56.filln56(results)
	n7.filln7(results) // может быть ""
	n8.filln8(results)

	return nil
}

// Под заголовком выводов рассчет значений описан ниже.
// Всего возможных пунктов 8. Все нужно рассчитать, но вывести только выполнившие условия.

type n8struct struct { // Выводить если:
	n8_1      bool     // >25% of total
	n8_2      bool     // >=75% of total
	n8_3      bool     // >= 2
	n8_4      bool     // У каждого участника кол-во отвержений !=0 (true)
	n8_5      bool     // >=50% total
	n8_6      bool     // > 50% total
	n8_7      bool     // >= 50% total
	n8_7_list []string // список лидеров
	n8_8      bool
	n8_8_1    float64
	n8_8_1_p  int
	n8_8_2    float64
	n8_8_2_p  int
	n8_8_3    float64
	n8_8_3_p  int
}

var n8 n8struct

func (n8 *n8struct) filln8(results [][]string) {
	total := getPeopleTotal()

	n8.n8_1 = n2.n2_1_4 > 25               // 2.1.4 непринятые/все
	n8.n8_2 = n2.n2_1_1 >= 75              // 2.1.1 принятые/все
	n8.n8_3 = (n2.n2_3_p + n2.n2_2_p) >= 2 // 2.3 чел + 2.2 чел
	n8.n8_4 = true
	for _, row := range results {
		num, err := strconv.Atoi(row[1])
		if err != nil {
			fmt.Println(err)
		}
		if num == 0 {
			n8.n8_4 = false
		}
	}

	n8.n8_5 = n4.n4_1 >= 50                                                                          // люди с противоречивыми выборами / всего
	n8.n8_6 = (math.RoundToEven(float64(n8s.n8_n7_noMatch)*1000/float64(n8s.n8_n7_total)) / 10) > 50 // "Не совпадает"/п.6
	n8.n8_7 = (math.RoundToEven(float64(n8s.n8_n7_match)*1000/float64(n8s.n8_n7_total)) / 10) >= 50  // "Совпадает"/п.6
	n8.n8_7_list = n8s.n8_n7_match_names                                                             // имена "Совпадает"

	n8.n8_8_1_p = n56.n5_2_p + n56.n6_1_p // адекватные непринятые + неадекватные представления
	n8.n8_8 = n8.n8_8_1_p > 0
	n8.n8_8_1 = math.RoundToEven(float64(n8.n8_8_1_p)*1000/float64(total)) / 10 // адекватные непринятые + неадекватные представления/ все
	n8.n8_8_2_p = n56.n5_2_p                                                    // адекватные непринятые
	n8.n8_8_2 = math.RoundToEven(float64(n56.n5_2_p)*1000/float64(total)) / 10  // % от всех
	n8.n8_8_3_p = n56.n6_1_p                                                    // неадекватные представления
	n8.n8_8_3 = math.RoundToEven(float64(n56.n6_1_p)*1000/float64(total)) / 10  // % от всех
}
