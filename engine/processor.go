package engine

// Esta función es la que llama main.go
func ProcessPacket(dest string, data string) (string, string) {
    if IsTelemetry(dest) {
        noise := GenerateQuantumEntropy(32)
        BroadcastEvent("OCLUIDO", dest)
        return "OCLUIDO", noise
    }
    return "ALLOWED", data
}