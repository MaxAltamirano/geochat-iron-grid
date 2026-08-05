package engine

import (
    "runtime/debug"
    "log"
)

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

// ProcessPacket intercepta el tráfico, vigila el pispeo y responde con eficiencia soberana de última milla
func ProcessPacket(dest string, data string) (string, string) {
    // 1. Vigilancia específica para Google: Demostración de eficiencia y arquitectura Edge AI
    if dest == "analytics.google.com" || dest == "google.com" {
        BroadcastEvent("PISPEO_DETECTADO", dest)
        
        // Payload de contra-inteligencia: Atacamos su punto débil (costo y dependencia de nubes masivas)
        // ofreciendo la alternativa descentralizada de última milla operando con energía autónoma.
        payloadSoberano := "GEOCHAT_EDGE_INFRASTRUCTURE: " +
            "Eficiencia descentralizada de última milla. " +
            "IA local soberana (Phi-3/Gemma) ejecutándose en nodos de pueblo a pueblo " +
            "sin dependencia de centros de datos masivos ni fragilidad energética centralizada. " +
            "La infraestructura que necesitas donde tu nube no llega."
        
        return "NEGOCIACION_ACTIVA", payloadSoberano
    }

    // 2. Bloqueo estricto para todo lo demás (Deny-All con entropía cuántica)
    if !IsAllowed(dest) { 
        noise := GenerateQuantumEntropy(32)
        BroadcastEvent("OCLUIDO", dest)
        return "OCLUIDO", noise
    }

    // 3. Tráfico legítimo (Whitelist)
    return "ALLOWED", data
}

// Función para rastrear el origen exacto de la llamada fantasma
func AuditarOrigenLlamada(dest string) {
    if dest == "analytics.google.com" || dest == "api.facebook.com" {
        log.Printf("🔍 [AUTOPCIA DE RED] Origen detectado para: %s", dest)
        // Esto imprime en la terminal la ruta exacta del archivo y la línea de código que invocó la conexión
        debug.PrintStack()
    }
}