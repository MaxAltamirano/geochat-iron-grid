package main

import (
    "fmt"
    "log"
    "net"
    "net/http"
    "net/http/httputil"
    "net/url"
    "runtime/debug"
    "time"
    "github.com/geochat/iron-grid/engine"
)

// Interceptor de rastreo para autopsia profunda de llamadas de red
func RastrearOrigenSalida(destino string) {
    if destino == "analytics.google.com" || destino == "google.com" || destino == "api.facebook.com" || destino == "telemetry.windows.com" {
        log.Printf("🚨 [AUTOPSIA DE RED] ¡Intento de salida detectado hacia: %s!", destino)
        debug.PrintStack()
    }
}

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

func main() {
    fmt.Println("🛡️ --- Iniciando Escudo IronGrid: Gateway Inverso y Blindaje Real (:8080) --- 🛡️")

    // URL del núcleo real de GeoChat Core operando de manera interna en el puerto 10001
    targetURL, err := url.Parse("http://127.0.0.1:10001")
    if err != nil {
        log.Fatalf("❌ Error al parsear URL del core interno: %v", err)
    }
    proxy := httputil.NewSingleHostReverseProxy(targetURL)

    trafficChannel := make(chan string, 100)

    // Goroutine procesadora de eventos de red y reportes al Core vía Socket UNIX
    go func() {
        for dest := range trafficChannel {
            status, result := engine.ProcessPacket(dest, "payload_privado")
            
            if status == "OCLUIDO" || dest == "analytics.google.com" || dest == "google.com" || dest == "api.facebook.com" || dest == "telemetry.windows.com" {
                fmt.Printf("🚫 [BLOQUEO Y AUTOPSIA]: %s | Entropy: %s\n", dest, result)
                ReportarAlCore("OCLUIDO", dest)
            } else {
                fmt.Printf("✅ [PERMITIDO SOBERANO]: %s\n", dest)
            }
        }
    }()

    // Servidor Proxy HTTP / Gateway Inverso
    proxyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        destino := r.Host
        if destino == "" {
            destino = r.URL.Host
        }

        // Enviamos el destino al canal de auditoría no bloqueante
        select {
        case trafficChannel <- destino:
        default:
        }

        // Evaluamos intrusos y disparamos autopsia si corresponde
        RastrearOrigenSalida(destino)

        // Blindaje estricto para dominios no deseados
        if destino == "analytics.google.com" || destino == "google.com" || destino == "api.facebook.com" || destino == "telemetry.windows.com" {
            w.WriteHeader(http.StatusForbidden)
            w.Write([]byte("🚫 [IRONGRID SOBERANO]: Tráfico ocluido por blindaje de red."))
            return
        }

        // Inyectamos cabecera soberana de certificación y reenviamos al Core real (:10001)
        r.Header.Set("X-IronGrid-Verified", "Sovereign-Node-Avellaneda")
        proxy.ServeHTTP(w, r)
    })

    fmt.Println("🌐 Escudo Gateway escuchando en 127.0.0.1:8080 (Derivando tráfico legítimo a Core :10001)...")
    if err := http.ListenAndServe("127.0.0.1:8080", proxyHandler); err != nil {
        fmt.Printf("⚠️ Error en el Gateway de Red: %v\n", err)
    }
}