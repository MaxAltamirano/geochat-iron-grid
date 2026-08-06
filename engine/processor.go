package engine

import (
    "net"
    //"os"
    "runtime/debug"
    "log"
)

// IsAllowed define la lista blanca estricta
func IsAllowed(dest string) bool {
    whitelist := map[string]bool{
        "google.com":           true,
        "accounts.google.com":  true,
        "gstatic.com":          true,
    }
    return whitelist[dest]
}

// NotificarSocketCore envía la alerta de intrusión o bloqueo al socket Unix del Core en tiempo real
func NotificarSocketCore(dominio string, estado string) {
    socketPath := "/tmp/geochat_core.sock"
    conn, err := net.Dial("unix", socketPath)
    if err != nil {
        return // Si el core no está escuchando de forma temporal, evitamos interrupciones
    }
    defer conn.Close()

    // Estructura limpia en formato JSON para que el Core emita el evento hacia Vue y suene la alarma de 432Hz
    payload := `{"dominio":"` + dominio + `", "estado":"` + estado + `"}`
    _, err = conn.Write([]byte(payload))
    if err != nil {
        log.Printf("⚠️ [IRON-GRID] Error al notificar al socket del Core: %v", err)
    }
}

// ProcessPacket intercepta el tráfico, vigila el pispeo y responde con eficiencia soberana de última milla
func ProcessPacket(dest string, data string) (string, string) {
    // 1. Vigilancia específica para Google: Demostración de eficiencia y arquitectura Edge AI
    if dest == "analytics.google.com" || dest == "google.com" {
        BroadcastEvent("PISPEO_DETECTADO", dest)
        
        // Notificamos también al socket para que quede registro visual en el nodo
        NotificarSocketCore(dest, "PISPEO_DETECTADO")
        
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
        
        // 👉 Notificación instantánea al socket del Core para activar la alerta y el sonido de 432Hz
        NotificarSocketCore(dest, "OCLUIDO")
        
        return "OCLUIDO", noise
    }

    // 3. Tráfico legítimo (Whitelist)
    return "ALLOWED", data
}

// Función para rastrear el origen exacto de la llamada fantasma
func AuditarOrigenLlamada(dest string) {
    if dest == "analytics.google.com" || dest == "api.facebook.com" {
        log.Printf("🔍 [AUTOPSIA DE RED] Origen detectado para: %s", dest)
        // Esto imprime en la terminal la ruta exacta del archivo y la línea de código que invocó la conexión
        debug.PrintStack()
    }
}