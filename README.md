# Analizador de Seguridad TLS

Aplicación en Go que utiliza la API de SSL Labs para analizar la configuración de seguridad TLS/SSL de dominios.

## Características

- Analiza la seguridad TLS usando la API de SSL Labs
- Utiliza el endpoint `getEndpointData` para obtener información detallada de endpoints
- Procesa y analiza los datos de la API para proporcionar conclusiones y recomendaciones
- Interfaz HTML limpia (sin Fprint)
- Manejo adecuado de errores con tipos de error personalizados
- Estructura de código modular
- Tests automatizados

## Estructura del Proyecto

```
practicas_go/
├── cmd/
│   └── server/
│       └── main.go          # Punto de entrada de la aplicación
├── internal/
│   ├── api/
│   │   ├── ssllabs.go       # Cliente de la API de SSL Labs
│   │   └── ssllabs_test.go  # Tests del cliente API
│   ├── analyzer/
│   │   ├── analyzer.go      # Lógica de análisis
│   │   └── analyzer_test.go # Tests del analizador
│   ├── errors/
│   │   └── errors.go        # Tipos de error personalizados
│   ├── handlers/
│   │   ├── handlers.go      # Handlers HTTP
│   │   └── handlers_test.go # Tests de handlers
│   └── models/
│       └── models.go        # Modelos de datos
├── templates/
│   ├── index.html           # Template del formulario
│   ├── result.html          # Template de resultados
│   └── error.html           # Template de errores
├── go.mod
└── README.md
```

## Uso

### Compilar y Ejecutar

```bash
go build ./cmd/server
./server
```

O ejecutar directamente:

```bash
go run ./cmd/server
```

El servidor se iniciará en `http://localhost:8080`

### Testing

Ejecutar todos los tests:

```bash
go test ./...
```

Ejecutar tests de un paquete específico:

```bash
go test ./internal/api
go test ./internal/analyzer
go test ./internal/handlers
```

## Cómo Funciona

1. El usuario ingresa un dominio en el formulario web
2. La aplicación inicia el análisis usando el endpoint `/analyze` de SSL Labs
3. Espera a que el análisis se complete (sondea hasta que el estado sea READY)
4. Para cada endpoint, obtiene datos detallados usando el endpoint `/getEndpointData`
5. Analiza los datos y genera:
   - Calificación general de seguridad
   - Información específica de cada endpoint
   - Recomendaciones de seguridad
6. Muestra los resultados en una página HTML formateada

## Manejo de Errores

La aplicación utiliza tipos de error personalizados:
- `DomainError`: Entrada de dominio inválida
- `APIError`: Errores de conexión o respuesta de la API
- `AnalysisError`: Errores específicos del análisis

Todos los errores se manejan de forma adecuada y se muestran al usuario mediante templates HTML.
