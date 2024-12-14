package config

import "time"

var (
	JwtDuration = 10 * 24 * time.Hour
	JwtSecret   = "fvrn3ui4uhucb32sq"
	MongoUri    = "mongodb://localhost:27017"
)
