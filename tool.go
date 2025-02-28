package main

import (
	"fmt"
	"strings"
	"time"
)

func GetTimestampFilename(file_name string, file_name_index int) string {
	timestamp := time.Now().Unix()
	parts := strings.Split(file_name, ".")
	parts[0] = fmt.Sprintf("%d_%d", timestamp, file_name_index)

	return strings.Join(parts, ".")
}
