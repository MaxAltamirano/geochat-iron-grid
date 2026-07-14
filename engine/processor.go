package engine

// Esta función es la que llama main.go
// IsAllowed define la lista blanca estricta
func IsAllowed(dest string) bool {
    whitelist := map[string]bool{
        "google.com":           true,
        "accounts.google.com":  true,
        "gstatic.com":          true,
    }
    return whitelist[dest]
}

// ProcessPacket ahora usa lógica inversa para asegurar la red
func ProcessPacket(dest string, data string) (string, string) {
    if !IsAllowed(dest) { 
        // Si no está autorizado, se activa la contrainteligencia
        noise := GenerateQuantumEntropy(32)
        BroadcastEvent("OCLUIDO", dest)
        return "OCLUIDO", noise
    }
    // Si está en la lista blanca, dejamos pasar el flujo legítimo
    return "ALLOWED", data
}