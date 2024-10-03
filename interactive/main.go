package main

import "fmt"

func main() {
	app := InitApp()
	for _, c := range app.consumers {
		err := c.Start("")
		if err != nil {
			panic(err)
		}
	}
	go func() {
		fmt.Println("migrator start")
		app.webAdmin.Start()
	}()

	err := app.server.Serve()
	panic(err)
}
