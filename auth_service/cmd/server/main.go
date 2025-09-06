package main

import (
	"github.com/adityawaradkar/gratia/auth_service/internal/db"
)

func main() {
	db.ConnectDB()
	// Keep service alive (for now, no server code)
	select {}
}
