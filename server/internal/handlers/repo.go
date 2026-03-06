package handlers

import "server/internal/config"

var app *config.AppConfig

func SetRepo(a *config.AppConfig) {
	app = a
}
