package main

import (
	"LlBlog/07/Databases"
	"LlBlog/07/Router"
)
func main()  {
	defer Databases.DB.Close()
	Router.InitRouter()
}
