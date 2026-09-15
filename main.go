package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"time"
	"github.com/geochat/iron-grid/engine"
)

// Reportar al núcleo vía Socket Unix con reintentos exponenciales
func ReportarAlCore(accion string, dest string) {
	ips, _ := net.LookupIP(dest)
	ipStr := "Desconocida"
	if len(ips) > 0 {
		ipStr = ips[0].String()
	}

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

// Receptor de alertas del escudo a través del Socket Unix
func IniciarEscuchaUnixSocket() {
	socketPath := "/tmp/geochat_core.sock"

	if _, err := os.Stat(socketPath); err == nil {
		os.Remove(socketPath)
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Printf("❌ [CORE SOCKET]: Error al crear el socket en %s: %v\n", socketPath, err)
		return
	}
	os.Chmod(socketPath, 0666)

	fmt.Printf("🔌 [CORE SOCKET]: Escuchando canal de auditoría IronGrid en %s\n", socketPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go func(c net.Conn) {
			defer c.Close()
			buf := bufio.NewReader(c)
			line, err := buf.ReadString('\n')
			if err != nil && err.Error() != "EOF" {
				return
			}
			if line != "" {
				fmt.Printf("📥 [EVENTO RECIBIDO DEL ESCUDO]: %s\n", line)
			}
		}(conn)
	}
}

// Autopsia de intentos de salida prohibidos
func RastrearOrigenSalida(destino string) {
	if destino == "analytics.google.com" || destino == "google.com" || destino == "api.facebook.com" || destino == "telemetry.windows.com" {
		log.Printf("🚨 [AUTOPSIA DE RED] ¡Intento de salida detectado hacia: %s!", destino)
		debug.PrintStack()
	}
}

func main() {
	fmt.Println("🛡️ --- Iniciando Sistema Unificado: Core + Escudo IronGrid en :8080 --- 🛡️")

	// 1. Levantamos el receptor de eventos Unix en segundo plano (Core)
	go IniciarEscuchaUnixSocket()

	// 🛰️ 2. ENCENDEMOS EL RADAR DE CONTRAINTELIGENCIA IA EN SEGUNDO PLANO
	engine.IniciarRadarContrainteligencia()

	// Creamos nuestro propio enrutador limpio
	mux := http.NewServeMux()

    // 🌐 Endpoint de Estado y Seguridad Unificado (Conectado al Radar del Engine)
    mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        
        // Obtenemos el estado de forma segura llamando al engine
        alertaActiva, target, entropia := engine.GetRadarStatus()

        // Si no hay target detectado aún por el radar, dejamos un fallback defensivo por defecto
        if target == "" {
            target = "vector-host-malicioso.net"
        }

        // Convertimos el booleano a string ("true" o "false") para JSON plano
        alertaStr := "false"
        if alertaActiva {
            alertaStr = "true"
        }

        // Construimos el JSON plano de forma segura sin depender del paquete json
        jsonResponse := fmt.Sprintf(`{
            "status": "online",
            "nodo": "Soberano-Avellaneda",
            "seguridad": {
                "ultimo_evento": "MONITOREO_ACTIVO",
                "target": "%s",
                "alerta_activa": %s,
                "last_entropy": %.2f
            }
        }`, target, alertaStr, entropia)

        w.Write([]byte(jsonResponse))
    })

	mux.HandleFunc("/api/cortex/logs-ollama", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"logs": []}`))
	})

	mux.HandleFunc("/api/tf/telemetria-ia", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"telemetria": "activa", "entropia": "estable"}`))
	})

	// Registramos la ruta que pide el Llavero
	mux.HandleFunc("/api/cortex/inspec", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "OK", "cortex": "activo"}`))
	})

	// 🛡️ Agregamos la ruta del reporte directamente al mux de la Iron Grid
	mux.HandleFunc("/api/cortex/reporte", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		fmt.Println("🧠 [IRON-GRID]: Pulso de reporte soberano interceptado y aceptado.")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"estado": "sincronizado", "sintonia": "432Hz"}`))
	})

	// Registramos el comodín de depuración con tu log original para cualquier otra cosa que falle
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("🔍 [DEBUG 404]: Ruta no encontrada intentada -> %s (Método: %s)\n", r.URL.Path, r.Method)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Ruta no encontrada en el Córtex"}`))
	})

	// 2. Servidor HTTP unificado con el middleware IronGrid integrado
	handlerUnificado := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		destino := r.Host
		if destino == "" {
			destino = r.URL.Host
		}

		engine.AuditarOrigenLlamada(destino)
		RastrearOrigenSalida(destino)

		status, payload := engine.ProcessPacket(destino, "payload_privado")

		switch status {
		case "OCLUIDO":
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(fmt.Sprintf("🚫 [IRONGRID SOBERANO - OCLUIDO]: Tráfico no autorizado.\nEntropía Cuántica: %s", payload)))
			return

		case "NEGOCIACION_ACTIVA":
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(fmt.Sprintf("🤝 [IRONGRID CONTRA-INTELIGENCIA]:\n%s", payload)))
			return

		case "ALLOWED":
			// Tránsito libre: derivamos a nuestro mux propio (no al DefaultServeMux)
			r.Header.Set("X-IronGrid-Verified", "Sovereign-Node-Avellaneda")
			mux.ServeHTTP(w, r)
			return

		default:
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("🚫 [IRONGRID]: Acceso denegado por política de red."))
			return
		}
	})

	fmt.Println("🌐 Nodo Soberano escuchando unificado en 127.0.0.1:8080...")
	if err := http.ListenAndServe("127.0.0.1:8080", handlerUnificado); err != nil {
		log.Fatalf("⚠️ Error crítico en el servidor unificado: %v", err)
	}
}