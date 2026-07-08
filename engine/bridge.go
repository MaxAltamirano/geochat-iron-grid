package engine

import (
	"fmt"
	"net"
)

// SOCKET_PATH define la ruta de comunicación soberana local
const SOCKET_PATH = "/tmp/geochat_core.sock"

// BroadcastEvent comunica eventos de seguridad al GeoChat-Core a través de un socket Unix.
// Si el Core no está activo, el sistema continúa operando sin bloquearse.
func BroadcastEvent(status string, target string) {
	// Intentamos conectar con el socket del Core
	conn, err := net.Dial("unix", SOCKET_PATH)
	if err != nil {
		// Si el Core no está escuchando, registramos el evento en consola (modo fallback)
		fmt.Printf("[IronGrid Bridge] Core no disponible, evento: %s -> %s\n", status, target)
		return
	}
	defer conn.Close()

	// Preparamos el mensaje en formato JSON para el Core
	message := fmt.Sprintf(`{"status":"%s", "target":"%s"}`, status, target)
	
	// Enviamos el mensaje al socket
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Printf("[IronGrid Bridge] Error al enviar mensaje: %v\n", err)
	}
}