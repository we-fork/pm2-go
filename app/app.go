package app

import (
	"net/rpc"
	"os"

	"github.com/dunstorm/pm2-go/shared"
	"github.com/dunstorm/pm2-go/utils"
	"github.com/rs/zerolog"
)

type App struct {
	client *rpc.Client
	logger *zerolog.Logger
}

func New() *App {
	return &App{
		logger: utils.NewLogger(),
	}
}

func (app *App) GetLogger() *zerolog.Logger {
	return app.logger
}

func (app *App) createClient() {
	var err error
	app.client, err = rpc.DialHTTP("tcp", "localhost:9001")
	if err != nil {
		app.logger.Fatal().Msgf("Connection error: %s", err.Error())
		os.Exit(1)
	}
}

func (app *App) AddProcess(process *shared.Process) shared.Process {
	var reply shared.Process
	app.createClient()
	defer app.client.Close()
	app.client.Call("API.AddProcess", process, &reply)
	return reply
}

func (app *App) ListProcs() []shared.Process {
	var db []shared.Process
	app.createClient()
	defer app.client.Close()
	app.client.Call("API.GetDB", "", &db)
	return db
}

func (app *App) FindProcess(name string) *shared.Process {
	var reply shared.Process
	app.createClient()
	defer app.client.Close()
	app.client.Call("API.FindProcess", name, &reply)
	return &reply
}

func (app *App) StopProcessByIndex(index int) bool {
	var reply bool
	app.createClient()
	defer app.client.Close()
	app.client.Call("API.StopProcessByIndex", index, &reply)
	return reply
}

func (app *App) StopProcessByName(name string) bool {
	var reply bool
	app.createClient()
	defer app.client.Close()
	app.client.Call("API.StopProcessByName", name, &reply)
	return reply
}

func (app *App) StartProcess(newProcess *shared.Process) *shared.Process {
	var reply *shared.Process
	app.createClient()
	defer app.client.Close()
	app.client.Call("API.UpdateProcess", newProcess, &reply)
	return reply
}

func (app *App) RestartProcess(process *shared.Process) *shared.Process {
	app.StopProcessByName(process.Name)
	newProcess := shared.SpawnNewProcess(shared.SpawnParams{
		Name:           process.Name,
		Args:           process.Args,
		ExecutablePath: process.ExecutablePath,
		AutoRestart:    process.AutoRestart,
		Logger:         app.logger,
	})
	process = app.StartProcess(newProcess)
	return process
}
