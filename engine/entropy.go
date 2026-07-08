package engine

import (
    "crypto/rand"
    "encoding/hex"
)

// GenerateQuantumEntropy genera ruido aleatorio seguro para oclusión.
func GenerateQuantumEntropy(length int) string {
    bytes := make([]byte, length/2)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}