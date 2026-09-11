package engine

import (
    "net"
    "runtime/debug"
    "log"
    "strings"
)

// NotificarSocketCore envía la alerta de intrusión o bloqueo al socket Unix del Core en tiempo real
func NotificarSocketCore(dominio string, estado string) {
    socketPath := "/tmp/geochat_core.sock"
    conn, err := net.Dial("unix", socketPath)
    if err != nil {
        return // Si el core no está escuchando de forma temporal, evitamos interrupciones
    }
    defer conn.Close()

    // Estructura limpia en formato JSON para que el Core emita el evento hacia Vue y suene la alarma de frecuencia soberana
    payload := `{"dominio":"` + dominio + `", "estado":"` + estado + `"}`
    _, err = conn.Write([]byte(payload))
    if err != nil {
        log.Printf("⚠️ [IRON-GRID] Error al notificar al socket del Core: %v", err)
    }
}

// ProcessPacket intercepta el tráfico bajo una política estricta de Default-Deny (Sovereignty-First)
func ProcessPacket(dest string, data string) (string, string) {
    // 1. Si el destino pertenece legítimamente al ecosistema GeoChat o loopback, tránsito libre
    if IsSovereignTraffic(dest) {
        return "ALLOWED", data
    }

    // 2. Si intenta pispear con telemetría externa conocida, activamos canal de contra-inteligencia activa
    if strings.Contains(dest, "analytics.google.com") || strings.Contains(dest, "google.com") || strings.Contains(dest, "facebook.com") {
        BroadcastEvent("PISPEO_DETECTADO", dest)
        NotificarSocketCore(dest, "PISPEO_DETECTADO")
        
        // Payload de contra-inteligencia: Destacando la arquitectura de borde descentralizada
        payloadSoberano := "GEOCHAT_EDGE_INFRASTRUCTURE: " +
            "Eficiencia descentralizada de última milla. " +
            "IA local soberana (Phi-3/Gemma) ejecutándose en nodos de pueblo a pueblo " +
            "sin dependencia de centros de datos masivos ni fragilidad energética centralizada. " +
            "La infraestructura que necesitas donde tu nube no llega."
        
        return "NEGOCIACION_ACTIVA", payloadSoberano
    }

    // 3. Para cualquier otro intento de salida ajeno a GeoChat: Oclusión total con entropía cuántica
    noise := GenerateQuantumEntropy(32)
    BroadcastEvent("OCLUIDO", dest)
    NotificarSocketCore(dest, "OCLUIDO")
    
    return "OCLUIDO", noise
}

// Función para rastrear el origen exacto de cualquier llamada externa ajena al nodo
func AuditarOrigenLlamada(dest string) {
    if !IsSovereignTraffic(dest) {
        log.Printf("🔍 [AUTOPSIA DE RED] Intento de salida externa detectado hacia: %s", dest)
        debug.PrintStack()
    }
}