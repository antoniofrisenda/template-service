package config

import (
	"github.com/gofiber/fiber/v3/middleware/logger"
)

var Format = logger.Config{
	Format:     `[${time}] |${status}| ${latency} |${method}| ${path} | IP: ${ip} | req: ${locals:requestid}` + "\n",
	TimeFormat: "2006-Jan-02 15:04:05",
	TimeZone:   "UTC",
}
