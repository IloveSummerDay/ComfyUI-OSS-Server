package main

import (
	"fmt"
	"strings"
	"time"
)

func GetTimestampFilename(file_name string) string {
	timestamp := time.Now().Unix()
	parts := strings.Split(file_name, ".")
	parts[0] = fmt.Sprintf("%d", timestamp)

	return strings.Join(parts, ".")
}
