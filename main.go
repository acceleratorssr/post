package main

func main() {
	app := InitApp()
	server := app.server
	topic := "article"
	for _, c := range app.consumers {
		err := c.Start(topic)
		// TODO 错误处理
		if err != nil {
			panic(err)
		}
	}
	err := server.Run(":9190")
	if err != nil {
		panic(err)
	}
}
