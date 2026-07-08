package main

import (
	"fmt"
	"net"
	"time"
	"github.com/geochat/iron-grid/engine"
)

func ReportarAlCore(accion string, detalle string) {
	conn, err := net.Dial("unix", "/tmp/geochat_core.sock")
	if err != nil {
		return // Fallo silencioso, el Core no está listo
	}
	defer conn.Close()

	message := fmt.Sprintf(`{"origen":"IRONGRID", "accion":"%s", "detalle":"%s"}`, accion, detalle)
	conn.Write([]byte(message))
}

func main() {
	fmt.Println("🛡️ --- Iniciando Escudo IronGrid: Modo Watcher Continuo --- 🛡️")

	// Canal de simulación de tráfico (esto en el futuro será tu captura de red)
	trafficChannel := make(chan string)

	// Goroutine que simula la entrada de paquetes constantemente
	go func() {
		packets := []string{"analytics.google.com", "api.facebook.com", "chat.geochat.org", "telemetry.windows.com"}
		for {
			for _, p := range packets {
				trafficChannel <- p
				time.Sleep(3 * time.Second) // Simula ráfagas de tráfico
			}
		}
	}()

	// Bucle infinito: El corazón del watcher
	for dest := range trafficChannel {
		status, result := engine.ProcessPacket(dest, "payload_privado")
		
		if status == "OCLUIDO" {
			fmt.Printf("🚫 [BLOQUEO]: %s | Entropy: %s\n", dest, result)
			// Reportamos al Core
			ReportarAlCore("OCLUIDO", dest)
		} else {
			fmt.Printf("✅ [PERMITIDO]: %s\n", dest)
		}
	}
}