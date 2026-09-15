package engine

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"
)

// --- ESTADO GLOBAL DE CONTRAINTELIGENCIA DEL ENGINE ---
type RadarState struct {
	mu           sync.Mutex
	AlertaActiva bool    `json:"alerta_activa"`
	Target       string  `json:"target"`
	ThreatType   string  `json:"threat_type"`
	LastEntropy  float64 `json:"last_entropy"`
}

var GlobalRadar = RadarState{}

// Estructuras de dialectos M2M / IA
type DialectV1 struct {
	Protocol string `json:"protocol"`
	Intent   string `json:"intent"`
}

type DialectV2 struct {
	MatrixSig string `json:"matrix_sig"`
	SyncMode  string `json:"sync_mode"`
}

// --- GETTER SEGURO PARA EL MAIN ---
func GetRadarStatus() (bool, string, float64) {
	GlobalRadar.mu.Lock()
	defer GlobalRadar.mu.Unlock()
	return GlobalRadar.AlertaActiva, GlobalRadar.Target, GlobalRadar.LastEntropy
}

// --- 1. GENERADOR DE RUIDO CUÁNTICO (ORIGINAL) ---
// GenerateQuantumEntropy genera ruido aleatorio seguro para oclusión.
func GenerateQuantumEntropy(length int) string {
	bytesData := make([]byte, length/2)
	rand.Read(bytesData)
	return hex.EncodeToString(bytesData)
}

// --- 2. MOTOR DE ENTROPÍA DE SHANNON ---
func CalculateEntropy(data string) float64 {
	if len(data) == 0 {
		return 0
	}
	freq := make(map[rune]float64)
	for _, char := range data {
		freq[char]++
	}
	var entropy float64 = 0
	length := float64(len(data))
	for _, count := range freq {
		p := count / length
		entropy -= p * (math.Log2(p))
	}
	return entropy
}

// --- 3. AUDITORÍA MULTIDIALECTO DE NODOS ---
func auditTargetNode(address string) {
	dialects := []map[string]any{
		{
			"dialect_id": "geochat-sync-v1",
			"payload": DialectV1{
				Protocol: "geochat-sync-v1",
				Intent:   "ai-contraintel-sweep",
			},
		},
		{
			"dialect_id": "autonomous-vector-v2",
			"payload": DialectV2{
				MatrixSig: "0xSovereignCore",
				SyncMode:  "m2m-mesh-probe",
			},
		},
	}

	client := http.Client{Timeout: 1 * time.Second}

	for _, d := range dialects {
		dialectID := d["dialect_id"].(string)
		url := fmt.Sprintf("http://%s/beacon", address)

		jsonData, err := json.Marshal(d["payload"])
		if err != nil {
			continue
		}

		resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		var buf bytes.Buffer
		buf.ReadFrom(resp.Body)
		responseString := buf.String()

		entropy := CalculateEntropy(responseString)

		GlobalRadar.mu.Lock()
		GlobalRadar.LastEntropy = entropy
		GlobalRadar.mu.Unlock()

		if entropy > 5.0 {
			fmt.Printf("[⚠️ RADAR IA] ¡Alta entropía (%.2f) detectada en %s! Dialecto [%s]. Posible canal encubierto.\n", entropy, address, dialectID)
			
			GlobalRadar.mu.Lock()
			GlobalRadar.AlertaActiva = true
			GlobalRadar.Target = address
			GlobalRadar.ThreatType = "Canal Encubierto M2M / IA Externa"
			GlobalRadar.mu.Unlock()
		} else {
			fmt.Printf("[Radar IA] Nodo %s seguro (Entropía: %.2f).\n", address, entropy)
		}
		return
	}
}

// --- 4. DEMONIO DE BARRIDO EN SEGUNDO PLANO ---
func IniciarRadarContrainteligencia() {
	go func() {
		// Objetivos perimetrales a escanear en segundo plano
		targetsToScan := []string{
			"127.0.0.1:8085",
		}

		ticker := time.NewTicker(10 * time.Second) // Barre cada 10 segundos
		for range ticker.C {
			for _, target := range targetsToScan {
				auditTargetNode(target)
			}
		}
	}()
	fmt.Println("🛰️ [ENGINE]: Radar de Contrainteligencia IA desplegado en segundo plano.")
}