package main

import (
	"context"
	"fmt"

	"github.com/DanVerh/artschool-admin/backend/api/application"
)

func main() {
	app := application.New()

	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println("failed to start app:", err)
	}
}
