package engine

// IsTelemetry verifica si el destino es un punto de telemetría no deseado.
func IsTelemetry(dest string) bool {
    telemetryEndpoints := map[string]bool{
        "analytics.google.com":  true,
        "telemetry.windows.com": true,
        "graph.facebook.com":    true,
    }
    return telemetryEndpoints[dest]
}