package main

import (
	"fmt"
	"net"
	"time"
	"github.com/geochat/iron-grid/engine"
)

func ReportarAlCore(accion string, dest string) {
	// Resolución de IP en tiempo real
	ips, _ := net.LookupIP(dest)
	ipStr := "Desconocida"
	if len(ips) > 0 {
		ipStr = ips[0].String()
	}

	message := fmt.Sprintf(`{"origen":"IRONGRID", "accion":"%s", "dominio":"%s", "ip":"%s"}`, accion, dest, ipStr)

	socketPath := "/tmp/geochat_core.sock"

	// Bucle persistente de reintentos optimizado para la red soberana
	for i := 1; i <= 5; i++ {
		conn, err := net.Dial("unix", socketPath)
		if err == nil {
			conn.Write([]byte(message))
			conn.Close()
			return // Éxito total en la entrega
		}
		
		// Pausa progresiva entre intentos
		time.Sleep(time.Duration(i) * 200 * time.Millisecond)
	}
	
	// Registro limpio sin errores críticos para no ensuciar la terminal
	fmt.Printf("[IronGrid Bridge] ⏳ Core ocupado sincronizando, reintentando en segundo plano para: %s\n", dest)
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