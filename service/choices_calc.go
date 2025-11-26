package service

import (
	"errors"
	"log"
)

//Работаем с социометрией
// Вопр. 1-2 = matrix1; Вопр. 3-4 = matrix2; Вопр.5 = matrix3
//1 = позитив (1 вопр), 2 = негатив (2 вопр), 0 = нет выбора
// Пример матрицы:
// [[0, 1],
//  [1, 1]]
// Так Matrix[0] = первая строка (не колонка)

var positiveChoices map[string]int
var negativeChoices map[string]int
var status map[string]string
var mutualNegativeChoices map[string]int
var mutualPositiveChoices map[string]int
var contradictoryChoices map[string]int
var autosociometryInts map[string]int
var autosociometryChoices map[string]string
var referentChoices map[string]int

// Названия файлов сохранения
var M1 string = "matrix1.json"
var M2 string = "matrix2.json"
var M3 string = "matrix3.json"

// Считаем сколько раз человека выбрали другие
func choiceCalc() {
	matrix, err := getMatrix(M1)
	if err != nil {
		log.Println("Error getting matrix for calcs")
	}

	positiveChoices = make(map[string]int)
	negativeChoices = make(map[string]int)
	for _, name := range SavedList {
		negativeChoices[name] = 0
		positiveChoices[name] = 0
	}
	for row_id := range matrix {
		for col_id, choice := range matrix[row_id] {
			if choice == 1 {
				positiveChoices[SavedList[col_id]] += 1
			} else if choice == 2 {
				negativeChoices[SavedList[col_id]] += 1
			}
		}
	}
}

// Считаем статус на основе выборов 1-2 вопросов
func statusCalc() {
	status = make(map[string]string)
	totalPpl := len(SavedList)
	for _, name := range SavedList {
		if float64(positiveChoices[name])/float64(totalPpl) > 0.5 {
			status[name] = "Звезда"
		} else if (float64(positiveChoices[name])/float64(totalPpl) > 0.25) && (float64(negativeChoices[name])/float64(totalPpl) < 0.1) {
			status[name] = "Популярные"
		} else if positiveChoices[name] == 0 && negativeChoices[name] > 0 {
			status[name] = "Отверженные"
		} else if positiveChoices[name] == 0 && negativeChoices[name] == 0 {
			status[name] = "Изолированные"
		} else {
			status[name] = "Не определено"
		}
	}
}

//Считаем взаимные (и противоречивые) выборы на основе вопр. 1-2
func mutualChoicesCalc() {
	matrix, err := getMatrix(M1)
	if err != nil {
		log.Println("Error getting matrix for calcs")
	}
	mutualNegativeChoices = make(map[string]int)
	mutualPositiveChoices = make(map[string]int)
	contradictoryChoices = make(map[string]int)
	for _, name := range SavedList {
		mutualNegativeChoices[name] = 0
		mutualPositiveChoices[name] = 0
		contradictoryChoices[name] = 0
	}
	for row_id := range matrix {
		for col_id, choice := range matrix[row_id] {
			if matrix[row_id][col_id] == matrix[col_id][row_id] && choice == 1 {
				mutualPositiveChoices[SavedList[col_id]] += 1
			} else if matrix[row_id][col_id] == matrix[col_id][row_id] && choice == 2 {
				mutualNegativeChoices[SavedList[col_id]] += 1
			} else if (matrix[row_id][col_id] != matrix[col_id][row_id]) && matrix[row_id][col_id] != 0 && matrix[col_id][row_id] != 0 {
				contradictoryChoices[SavedList[col_id]] += 1
			}
		}
	}
}

//Считаем аутосоциометрию количественно для каждого тестируемого
func autosociometry() error {
	matrix1, err := getMatrix(M1)
	if err != nil {
		log.Println("Error getting matrix2 for calcs")
	}
	matrix2, err := getMatrix(M2)
	if err != nil {
		log.Println("Error getting matrix2 for calcs")
	}
	autosociometryInts = make(map[string]int)
	for _, name := range SavedList {
		autosociometryInts[name] = 0
	}

	//Сначала проверим совпадение размеров матриц, для предотвращения ошибок
	if len(matrix1) != len(matrix2) {
		return errors.New("cannot calculate autosociometry: different size matrices")
	} else if len(matrix1) <= 1 {
		return errors.New("incorrect matrix size")
	} else if len(matrix1[0]) != len(matrix2[0]) {
		return errors.New("cannot calculate autosociometry: different size matrices")
	}
	// Поскольку матрица одного размера и с одиноковыми названиями столбцов\колонок
	// мы можем проходить по одной и сверять лишь значения в клетках
	for row_id := range matrix2 {
		for col_id, choice := range matrix2[row_id] {
			if (matrix2[row_id][col_id] == matrix1[col_id][row_id]) && choice != 0 {
				autosociometryInts[SavedList[row_id]] += 1
			}
		}
	}

	return nil
}

//переводим численные значения аутосоциометрии в статус
func autosociometryCalc() {
	autosociometryChoices = make(map[string]string)
	err := autosociometry()
	if err != nil {
		log.Println(err)
	}

	for key, value := range autosociometryInts {
		if value <= 2 {
			autosociometryChoices[key] = "Неадекватная"
		} else if value > 2 {
			autosociometryChoices[key] = "Адекватная"
		} else {
			autosociometryChoices[key] = "N/D"
		}
	}
}

//Считаем кол-во референтных выборов (сколько раз человека выбрали в вопр.5)
func referentChoiceCalc() {
	matrix, err := getMatrix(M3)
	if err != nil {
		log.Println("Error getting matrix3 for calcs")
	}
	referentChoices = make(map[string]int)
	for _, name := range SavedList {
		referentChoices[name] = 0
	}

	for row_id := range matrix {
		for col_id, choice := range matrix[row_id] {
			if choice == 1 {
				referentChoices[SavedList[col_id]] += 1
			}
		}
	}
}
