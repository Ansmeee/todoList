package main

import "todoList/bootstrap"

func main()  {
	app := bootstrap.InitEngine()
	app.Run("127.0.0.1:8000")
}
