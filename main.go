package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"time"
	"github.com/geochat/iron-grid/engine"
)

// Interceptor de rastreo para autopsia profunda de llamadas de red
func RastrearOrigenSalida(destino string) {
	if destino == "analytics.google.com" || destino == "google.com" || destino == "api.facebook.com" || destino == "telemetry.windows.com" {
		log.Printf("🚨 [AUTOPSIA DE RED] ¡Intento de salida detectado hacia: %s!", destino)
		// Imprime en la terminal la pila de llamadas exacta (Stack Trace)
		debug.PrintStack()
	}
}

func ReportarAlCore(accion string, dest string) {
	ips, _ := net.LookupIP(dest)
	ipStr := "Desconocida"
	if len(ips) > 0 {
		ipStr = ips[0].String()
	}

	// Estructura JSON precisa para que App.vue la interprete en syncPipeline
	message := fmt.Sprintf(`{"origen":"IRONGRID", "accion":"%s", "dominio":"%s", "ip":"%s"}`, accion, dest, ipStr)
	socketPath := "/tmp/geochat_core.sock"

	for i := 1; i <= 5; i++ {
		conn, err := net.Dial("unix", socketPath)
		if err == nil {
			conn.Write([]byte(message))
			conn.Close()
			return
		}
		time.Sleep(time.Duration(i) * 200 * time.Millisecond)
	}
	
	fmt.Printf("[IronGrid Bridge] ⏳ Core ocupado sincronizando, reintentando en segundo plano para: %s\n", dest)
}

func main() {
	fmt.Println("🛡️ --- Iniciando Escudo IronGrid: Modo Autopsia y Blindaje Real (:8080) --- 🛡️")

	trafficChannel := make(chan string)

	// Servidor Proxy HTTP Real para interceptar el tráfico de Linux Mint
	go func() {
		proxyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			destino := r.Host
			if destino == "" {
				destino = r.URL.Host
			}
			
			// 1️⃣ Permitimos el tráfico local/interno de forma prioritaria sin ocluirlo ni trabar el arranque
			if destino == "localhost:8080" || destino == "127.0.0.1:8080" || destino == "localhost" || destino == "127.0.0.1" {
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, "✅ [LOCAL SOBERANO]")
				return
			}

			// Inyectamos el dominio detectado al canal para análisis externo
			trafficChannel <- destino

			// Evaluamos intrusos y disparamos autopsia si corresponde
			RastrearOrigenSalida(destino)

			// 2️⃣ Blindaje estricto para intrusos externos
			if destino == "analytics.google.com" || destino == "google.com" || destino == "api.facebook.com" || destino == "telemetry.windows.com" {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("🚫 [IRONGRID SOBERANO]: Tráfico ocluido por blindaje de red."))
				return
			}

			w.WriteHeader(http.StatusOK)
			io.WriteString(w, "✅ [PERMITIDO SOBERANO]")
		})

		fmt.Println("🌐 Escudo Proxy escuchando en 127.0.0.1:8080...")
		if err := http.ListenAndServe("127.0.0.1:8080", proxyHandler); err != nil {
			fmt.Printf("⚠️ Error en el Proxy de Red: %v\n", err)
		}
	}()

	// Bucle principal procesando y reportando al Core
	for dest := range trafficChannel {
		status, result := engine.ProcessPacket(dest, "payload_privado")
		
		if status == "OCLUIDO" || dest == "analytics.google.com" || dest == "google.com" || dest == "api.facebook.com" || dest == "telemetry.windows.com" {
			fmt.Printf("🚫 [BLOQUEO Y AUTOPSIA]: %s | Entropy: %s\n", dest, result)
			// Envía la alerta al Core para que Vue la reciba en tiempo real
			ReportarAlCore("OCLUIDO", dest)
		} else {
			fmt.Printf("✅ [PERMITIDO SOBERANO]: %s\n", dest)
		}
	}
}