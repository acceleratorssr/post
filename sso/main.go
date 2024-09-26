package main

func main() {
	app := InitApp()
	err := app.server.Serve()
	panic(err)
}
