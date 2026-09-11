// engine/occlude.go modificado para Soberanía Estricta (Lista Blanca)
package engine

import "strings"

// IsSovereignTraffic verifica si el destino pertenece legítimamente al ecosistema GeoChat o tráfico local.
func IsSovereignTraffic(dest string) bool {
    // 1. Tráfico local y de desarrollo
    if strings.HasPrefix(dest, "localhost") || strings.HasPrefix(dest, "127.0.0.1") || strings.HasPrefix(dest, "[::1]") {
        return true
    }

    // 2. Dominios o nodos autorizados de la red soberana GeoChat (ejemplo de sufijo o lista permitida)
    if strings.HasSuffix(dest, ".geochat.internal") || strings.HasPrefix(dest, "node-") {
        return true
    }

    // Todo lo demás ajeno al ecosistema se rechaza por defecto
    return false
}