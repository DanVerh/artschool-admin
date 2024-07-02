package main

import (
	"github.com/DanVerh/artschool-admin/backend/migration/application"
)

func main() {
	app := application.New()
	app.Start()
}