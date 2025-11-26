package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	service "soc_test/service"
)

// App struct
type App struct {
	ctx context.Context
	list *service.ListService
	matrix *service.MatrixService
	result *service.ResultService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		list: service.NewListService(),
		matrix: service.NewMatrixService(),
		result: service.NewResultService(),
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

//Matrix service

func (a *App) SaveMatrix(filename string, data [][]int, uuid string) error {
	return a.matrix.SaveMatrix(filename, data, uuid)
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