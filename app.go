package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	service "soc_test/service"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
	list *service.ListService
	matrix *service.MatrixService
	result *service.ResultService
	excel *service.ExcelService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		list: service.NewListService(),
		matrix: service.NewMatrixService(),
		result: service.NewResultService(),
		excel: service.NewExcelService(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("Невозможно найти конфиг пользователя: ", err)
	}
	service.SaveDir = filepath.Join(configDir, "BitianovaSociometry", "SaveData")
	if err := os.MkdirAll(service.SaveDir, 0755); err != nil {
		log.Fatal("Не удалось создать директорию сохранений:", err)
	}

}


//List service
func (a *App) ListAdd(name string) error {
	return a.list.ListAdd(name)
}

func (a *App) ListGet() []string {
	return a.list.ListGet()
}

func (a *App) ListRemove(id int) {
	a.list.ListRemove(id)
}

func (a *App) ListExport() service.ListFile {
	return a.list.ListExport()
}

func (a *App) SaveTestDate(date string) {
	a.list.SaveTestDate(date)
}

//Matrix service

func (a *App) SaveMatrix(filename string, data [][]int, uuid string) error {
	return a.matrix.SaveMatrix(filename, data, uuid)
}

func (a *App) SaveOriginalMatrix(filename string, data [][]string) error {
	return a.matrix.SaveOriginalMatrix(filename, data)
}


func (a *App) LoadMatrix(filename string) (string, error) {
	return a.matrix.LoadMatrix(filename)
}

// Result service

func (a *App) GetResults() service.Results {
	return a.result.GetResults()
}

func (a *App) SaveResult(filename string, data [][]string, uuid string) error {
	return a.result.SaveResult(filename, data, uuid)
}

func (a *App) LoadResult(filename string) (string, error) {
	return a.result.LoadResult(filename)
}

func (a *App) CheckMatrices() bool {
	return a.result.CheckMatrices()
}

func (a *App) CheckResults(currentUUID string) bool {
	return a.result.CheckResults(currentUUID)
}

//Excel service

func (a *App) CreateExcelFile(fullPath string) error {
	return a.excel.CreateExcelFile(fullPath)
}

func (a *App) ShowSaveExcelDialog() (string, error) {
    filePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
        Title:                "Сохранить файл Excel",
        DefaultFilename:      "Социометрия Битяновой.xlsx",
        CanCreateDirectories: true,
        Filters: []runtime.FileFilter{
            {
                DisplayName: "Файлы Excel (*.xlsx)",
                Pattern:     "*.xlsx",
            },
        },
    })

    if err != nil {
        return "", err
    }
    if filePath == "" {
        return "", fmt.Errorf("user cancelled")
    }

    // Force .xlsx extension if missing
    if !strings.HasSuffix(strings.ToLower(filePath), ".xlsx") {
        filePath += ".xlsx"
    }

    return filePath, nil
}