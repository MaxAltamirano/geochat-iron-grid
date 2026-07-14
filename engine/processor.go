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
    // 1. Vigilancia específica para Google (Permitido pero reportado)
    if dest == "analytics.google.com" || dest == "google.com" {
        BroadcastEvent("PISPEO_DETECTADO", dest)
        return "ALLOWED", data
    }

    // 2. Bloqueo estricto para todo lo demás (Deny-All)
    if !IsAllowed(dest) { 
        noise := GenerateQuantumEntropy(32)
        BroadcastEvent("OCLUIDO", dest)
        return "OCLUIDO", noise
    }

    // 3. Tráfico legítimo (Whitelist)
    return "ALLOWED", data
}