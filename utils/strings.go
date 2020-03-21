package utils

import (
    "crypto/rand"
    "fmt"
)

// GenerateRandomString a random string of A-Z chars with len = l
func GenerateRandomString(l int) string {
    bytes := make([]byte, l)
    rand.Read(bytes)
    return fmt.Sprintf("%x", bytes)
}
