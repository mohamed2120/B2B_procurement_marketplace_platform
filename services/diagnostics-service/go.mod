module github.com/b2b-platform/diagnostics-service

go 1.23.0

require (
	github.com/b2b-platform/shared v0.0.0
	github.com/gin-contrib/cors v1.7.6
	github.com/gin-gonic/gin v1.10.1
	github.com/lib/pq v1.10.9
	github.com/pressly/goose/v3 v3.15.1
	gorm.io/driver/postgres v1.5.4
	gorm.io/gorm v1.25.5
)

replace github.com/b2b-platform/shared => ../../shared
