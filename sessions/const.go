package main

import "time"

var redisAddr string = "redis://user:@localhost:6379/0"
var tcpPort string = ":8083"
var expireTime time.Duration = 24 * 31 * time.Hour
