package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"

	app "github.com/nevinmanoj/hostmate/internal/app"
)

func main() {

	fmt.Println("Starting HOSTMATE API service...")

	fmt.Println(`
 __  __     ______     ______     ______       __    __     ______     ______     ______    
/\ \_\ \   /\  __ \   /\  ___\   /\__  _\     /\ "-./  \   /\  __ \   /\__  _\   /\  ___\   
\ \  __ \  \ \ \/\ \  \ \___  \  \/_/\ \/     \ \ \-./\ \  \ \  __ \  \/_/\ \/   \ \  __\   
 \ \_\ \_\  \ \_____\  \/\_____\    \ \_\      \ \_\ \ \_\  \ \_\ \_\    \ \_\    \ \_____\ 
  \/_/\/_/   \/_____/   \/_____/     \/_/       \/_/  \/_/   \/_/\/_/     \/_/     \/_____/ 
 `)

	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}

	if err := app.Start(); err != nil {
		fmt.Println("Error starting server:", err)
	}

}
