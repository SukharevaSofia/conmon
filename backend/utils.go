package main

import (
	"time"
)


type tableRow struct{
  IP string `json:"ip"`
  Ping_time time.Time `json:"ping_time"`
  Last_success time.Time `json:"last_success"`
}

