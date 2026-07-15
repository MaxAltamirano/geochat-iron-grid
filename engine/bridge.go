package engine

import (
	"fmt"
	"net"
)

// SOCKET_PATH define la ruta de comunicación soberana local
const SOCKET_PATH = "/tmp/geochat_core.sock"

// BroadcastEvent comunica eventos de seguridad al GeoChat-Core a través de un socket Unix.
// Si el Core no está activo, el sistema continúa operando sin bloquearse.
// Ahora apuntamos al socket del Buzón (SNC) en lugar de un Core independiente
const BUZON_SOCKET_PATH = "/tmp/geochat_buzon.sock"

// BroadcastEvent ahora inyecta los eventos de seguridad directamente en el Sistema Nervioso Central (Buzón)
func BroadcastEvent(status string, target string) {
	// Conectamos al socket del Buzón que gestiona el Radar y la Telemetría
	conn, err := net.Dial("unix", BUZON_SOCKET_PATH)
	if err != nil {
		// Fallback: Si el Buzón no responde, mantenemos el log local pero sin error crítico
		fmt.Printf("[IronGrid SNC-Bridge] Buzón en silencio. Evento: %s | Target: %s\n", status, target)
		return
	}
	defer conn.Close()

	// Transformamos el evento de seguridad en un objeto de "Interferencia" para el Radar
	// Esto hará que el Buzón lo procese como un objeto detectado en el espacio aéreo
	message := fmt.Sprintf(`{"type":"INTERFERENCIA", "status":"%s", "target":"%s", "origen":"IRON_GRID"}`, status, target)
	
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Printf("[IronGrid SNC-Bridge] Error de sincronización con el Buzón: %v\n", err)
	} else {
		fmt.Printf("[IronGrid SNC-Bridge] Evento inyectado en el Radar: %s\n", target)
	}
}