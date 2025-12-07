package service

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/xuri/excelize/v2"
)

// var excelFilename string = "Социометрия Битяновой.xlsx"

type ExcelService struct{}

func NewExcelService() *ExcelService {
	return &ExcelService{}
}

func (s *ExcelService) CreateExcelFile(fullPath string) error {
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("ошибка создания excel файла %s\n", err)
			return
		}
	}()

	om1, err := getOriginalMatrix("original_matrix1.json")
	if err != nil {
		fmt.Println("ошибка открытия оригинальной матрицы")
		return fmt.Errorf("ошибка открытия оригинальной матрицы: %s", err)
	}
	om2, err := getOriginalMatrix("original_matrix2.json")
	if err != nil {
		fmt.Println("ошибка открытия оригинальной матрицы")
		return fmt.Errorf("ошибка открытия оригинальной матрицы: %s", err)
	}
	om3, err := getOriginalMatrix("original_matrix3.json")
	if err != nil {
		fmt.Println("ошибка открытия оригинальной матрицы")
		return fmt.Errorf("ошибка открытия оригинальной матрицы: %s", err)
	}
	res, err := getOriginalMatrix("original_results.json")
	if err != nil {
		fmt.Println("ошибка открытия оригинальной матрицы")
		return fmt.Errorf("ошибка открытия оригинальной матрицы: %s", err)
	}
	createExcelSheet(file, "Вопросы 1-2", om1)
	file.DeleteSheet("Sheet1")
	createExcelSheet(file, "Вопросы 3-4", om2)
	createExcelSheet(file, "Вопрос 5", om3)
	createResultsExcelSheet(file, "Результаты", res)

	// fmt.Println("Полный путь до файла: ", fullPath)
	if err := file.SaveAs(fullPath); err != nil {
		fmt.Println("ошибка сохранения таблицы Excel", err)
		return fmt.Errorf("ошибка сохранения таблицы Excel: %s", err)
	}

	return nil

}

func createExcelSheet(f *excelize.File, sheetName string, sliceInterface [][]interface{}) {

	// Создать новый лист
	_, err := f.NewSheet(sheetName)
	if err != nil {
		fmt.Println("ошибка создания нового рабочего листа excel", err)
		return
	}
	//Ставим дату проведения теста
	err = f.SetCellValue(sheetName, "A1", "Дата проведения теста")
	if err != nil {
		fmt.Println("ошибка заполнения даты", err)
		return
	}
	err = f.SetCellValue(sheetName, "A2", testDate)
	if err != nil {
		fmt.Println("ошибка заполнения даты", err)
		return
	}
	//Ставим кол-во участников
	err = f.SetCellValue(sheetName, "B1", "Кол-во тестируемых")
	if err != nil {
		fmt.Println("ошибка заполнения кол-ва", err)
		return
	}
	err = f.SetCellValue(sheetName, "B2", len(names))
	if err != nil {
		fmt.Println("ошибка заполнения кол-ва", err)
		return
	}

	//Заполняем заголовки и первую колонку именами
	err = f.SetSheetRow(sheetName, "D3", getListInterface())
	if err != nil {
		fmt.Println("ошибка заполнения строки эксель", err)
		return
	}
	err = f.SetSheetCol(sheetName, "C4", getListInterface())
	if err != nil {
		fmt.Println("ошибка заполнения колонки эксель", err)
		return
	}

	//Заполняем содержимое из [][]interface{}

	var (
		// значения ячеек
		data = sliceInterface
		addr string
	)
	// установить значение каждой ячейки
	for r, row := range data {
		if addr, err = excelize.JoinCellName("D", r+4); err != nil {
			fmt.Println("ошибка заполнения данных эксель", err)
			return
		}
		if err = f.SetSheetRow(sheetName, addr, &row); err != nil {
			fmt.Println("ошибка заполнения данных эксель 2", err)
			return
		}
	}

	rows_slice, err := f.GetRows(sheetName)
	lastRow := len(rows_slice)
	cols_slice := rows_slice[lastRow-1]
	lastColIndex := len(cols_slice)

	lastCol, _ := excelize.ColumnNumberToName(lastColIndex)
	lastCell := fmt.Sprintf("%s%d", lastCol, lastRow)

	wrap_style, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText: true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	f.SetCellStyle(sheetName, "C3", lastCell, wrap_style)

	//настраиваем ширину колонок
	cols, err := f.GetCols(sheetName)
	if err != nil {
		fmt.Println("ошибка корректировки ширины колонок", err)
	}
	for idx, col := range cols {
		largestWidth := 10
		for _, rowCell := range col {
			cellWidth := utf8.RuneCountInString(rowCell) + 2 // + 2 for margin
			if cellWidth > largestWidth {
				largestWidth = cellWidth
			}
		}
		name, err := excelize.ColumnNumberToName(idx + 1)
		if err != nil {
			fmt.Println("ошибка корректировки ширины колонок", err)
		}
		f.SetColWidth(sheetName, name, name, float64(largestWidth))
	}

}

var resultsLables = []interface{}{"ФИО",
	"Положительные выборы",
	"Отрицательные выборы",
	"Статус",
	"Кол-во взаимных социометрических выборов (положительных)",
	"Кол-во взаимных социометрических выборов (отрицательных)",
	"Кол-во противоречивых выборов",
	"Аутосоциометрия",
	"Кол-во референтных выборов",
	"Целевая группа"}

func createResultsExcelSheet(f *excelize.File, sheetName string, sliceInterface [][]interface{}) {
	// Создать новый лист
	_, err := f.NewSheet(sheetName)
	if err != nil {
		fmt.Println("ошибка создания нового рабочего листа excel", err)
		return
	}
	//Ставим дату проведения теста
	err = f.SetCellValue(sheetName, "A1", "Дата проведения теста")
	if err != nil {
		fmt.Println("ошибка заполнения даты", err)
		return
	}
	err = f.SetCellValue(sheetName, "A2", testDate)
	if err != nil {
		fmt.Println("ошибка заполнения даты", err)
		return
	}

	//Ставим кол-во участников
	err = f.SetCellValue(sheetName, "B1", "Кол-во тестируемых")
	if err != nil {
		fmt.Println("ошибка заполнения кол-ва", err)
		return
	}
	err = f.SetCellValue(sheetName, "B2", len(names))
	if err != nil {
		fmt.Println("ошибка заполнения кол-ва", err)
		return
	}

	//Заполняем заголовки и первую колонку именами
	err = f.SetSheetRow(sheetName, "C3", &resultsLables)
	if err != nil {
		fmt.Println("ошибка заполнения строки эксель", err)
		return
	}
	err = f.SetSheetCol(sheetName, "C4", getListInterface())
	if err != nil {
		fmt.Println("ошибка заполнения колонки эксель", err)
		return
	}

	var (
		// значения ячеек
		data = sliceInterface
		addr string
	)
	// установить значение каждой ячейки
	for r, row := range data {
		if addr, err = excelize.JoinCellName("D", r+4); err != nil {
			fmt.Println("ошибка заполнения данных эксель", err)
			return
		}
		if err = f.SetSheetRow(sheetName, addr, &row); err != nil {
			fmt.Println("ошибка заполнения данных эксель 2", err)
			return
		}
	}

	leng := len(names) + 3
	tableSize := "C3:L" + strconv.Itoa(leng)
	fmt.Println(tableSize)

	//Создаём таблицу
	err = f.AddTable(sheetName, &excelize.Table{Range: tableSize,
		Name:      "ResultTable",
		StyleName: "TableStyleMedium11",
	})
	if err != nil {
		fmt.Println("ошибка создания таблицы", err)
		return
	}

	rows_slice, err := f.GetRows(sheetName)
	lastRow := len(rows_slice)
	cols_slice := rows_slice[lastRow-1]
	lastColIndex := len(cols_slice)

	lastCol, _ := excelize.ColumnNumberToName(lastColIndex)
	lastCell := fmt.Sprintf("%s%d", lastCol, lastRow)

	wrap_style, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText: true,
		},
	})
	f.SetCellStyle(sheetName, "C3", lastCell, wrap_style)

	//настраиваем ширину колонок
	cols, err := f.GetCols(sheetName)
	if err != nil {
		fmt.Println("ошибка корректировки ширины колонок", err)
	}
	for idx, _ := range cols {
		largestWidth := 20
		// for _, rowCell := range col {
		// 	cellWidth := utf8.RuneCountInString(rowCell) + 2 // + 2 for margin
		// 	if cellWidth > largestWidth && cellWidth <= 30 {
		// 		largestWidth = cellWidth
		// 	}
		// }
		name, err := excelize.ColumnNumberToName(idx + 1)
		if err != nil {
			fmt.Println("ошибка корректировки ширины колонок", err)
		}
		f.SetColWidth(sheetName, name, name, float64(largestWidth))
	}

}

//Пример заполнения:
// var (
//     // значения ячеек
//     data = [][]interface{}{
//         {"Fruits", "Vegetables"},
//         {"Mango", "Potato", nil, "Drop Down 1", "Drop Down 2"},
//         {"Apple", "Tomato"},
//         {"Grapes", "Spinach"},
//         {"Strawberry", "Onion"},
//         {"Kiwi", "Cucumber"},
//     }
//     addr                    string
//     err                     error
//     cellsStyle, headerStyle int
// )
// // установить значение каждой ячейки
// for r, row := range data {
//     if addr, err = excelize.JoinCellName("A", r+1); err != nil {
//         fmt.Println(err)
//         return
//     }
//     if err = f.SetSheetRow("Sheet1", addr, &row); err != nil {
//         fmt.Println(err)
//         return
//     }
// }

//OR

// // Установленное значение ячейки
// f.SetCellValue("Sheet2", "A2", "Hello world.")
// f.SetCellValue("Sheet1", "B2", 100)
// // Установить активный лист рабочей книги
// f.SetActiveSheet(index)
// // Сохранить файл xlsx по данному пути
// if err := f.SaveAs("Book1.xlsx"); err != nil {
//     fmt.Println(err)
// }
