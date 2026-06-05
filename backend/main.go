package main

import db "backend/dao"

func main() {
	db.Connect()
	db.Migrate()
}
