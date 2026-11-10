# Actualizaciones pendientes de aprobación

No se cambió ninguna versión del DAS ni ningún lockfile. Las versiones corregidas siguientes
se verificaron en avisos oficiales y registros de publicación; no se instalaron ni se afirma
que el grafo propuesto haya pasado pruebas. Evidencias: `evidencias/avisos-go-verificados.json`,
`avisos-imagenes-verificados.json`, `avisos-go-imagenes.json`, `avisos-xnet-complemento.json`, `versiones-propuestas-modulos.json`,
`vitest-propuesta.json`, `pnpm-audit.json` y `imagenes-validacion.json` dentro de evidencias.

## Dependencias y alcance

| Dependencia                           | Actual                   | Avisos                                                              | Impacto / alcanzabilidad observada                                                                                                                                                                      | Corrección propuesta verificada                                              | Archivos por modificar tras aprobación                                                                                  |
| ------------------------------------- | ------------------------ | ------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------- |
| Go / biblioteca estándar              | 1.26.5                   | 16 GO alcanzables, lista abajo; incluye GO-2026-6809                | HTTP, TLS/X.509 y parsers: DoS, validación de certificados y otros impactos según aviso. govulncheck detecta rutas llamadas; no prueba explotación. Trivy identifica 12 HIGH repetidos en cada binario. | 1.26.10, máximo de fixes de la rama 1.26 entre estos avisos                  | apps/api/go.mod, .github/workflows/ci.yml, los tres Dockerfiles, scripts/entorno-local.ps1 y documentación de versiones |
| golang.org/x/text                     | API 0.37.0; Goose 0.21.0 | GO-2026-6629 / CVE-2026-56851; GO-2026-5970 / CVE-2026-56852        | DoS al procesar texto; API: rutas llamadas desde inicialización de pgx. Goose: presencia en binario, sin prueba de alcanzabilidad                                                                       | 0.41.0 (el segundo aviso se corrige desde 0.39.0)                            | apps/api/go.mod y go.sum; módulo de herramientas Goose propuesto                                                        |
| Vite                                  | 8.1.5                    | GHSA-587h-7wmj-h5pm (HIGH); otros 4 avisos en audit                 | launch-editor en Windows y servidor dev; no se incluye como servidor en runtime SPA Go. Linux CI no justifica ignorarlo                                                                                 | 8.1.6; rango vulnerable de esta serie >=8.0.0 <=8.1.5                        | pnpm-workspace.yaml, pnpm-lock.yaml                                                                                     |
| Vitest                                | 4.1.10                   | GHSA-82fw-gwwq-j7x9                                                 | Lectura arbitraria/traversal en mocker durante desarrollo; no runtime productivo                                                                                                                        | 4.1.11 corrige este aviso; proponer 4.1.12 para alinear mocker y coverage-v8 | pnpm-workspace.yaml, pnpm-lock.yaml                                                                                     |
| @vitest/mocker                        | 4.1.10                   | GHSA-82fw-gwwq-j7x9; GHSA-qxg6-89mg-rqv8                            | Lectura e inyección mediante registro de mocks; entorno de pruebas, explotación local no demostrada                                                                                                     | 4.1.12; rangos >=2.1.0 <4.1.11 y >=2.1.0-beta.7 <4.1.12                      | pnpm-lock.yaml mediante actualización alineada Vitest/coverage-v8 en catalog                                            |
| github.com/golang-jwt/jwt/v4 (Goose)  | 4.5.1                    | CVE-2025-30204                                                      | DoS en parseo JWT; incluido con otros drivers, uso por migración PostgreSQL no demostrado                                                                                                               | 4.5.2                                                                        | Módulo Goose propuesto, Dockerfile.migrate y scripts/tareas.mjs                                                         |
| github.com/jackc/pgx/v5 (Goose)       | 5.7.1                    | CVE-2026-33815, CVE-2026-33816                                      | Seguridad de memoria en driver PostgreSQL utilizado; no se ha reproducido explotación de los caminos afectados                                                                                          | 5.9.0                                                                        | Módulo Goose propuesto; API ya usa 5.10.0 y no requiere bajar versión                                                   |
| github.com/klauspost/compress (Goose) | 1.17.7                   | CVE-2026-63209                                                      | Overflow y escritura fuera de límites en compresión S2; presencia no equivale a uso en esta migración                                                                                                   | 1.18.7                                                                       | Módulo Goose propuesto                                                                                                  |
| golang.org/x/crypto (Goose)           | 0.31.0                   | 12 CVE, lista abajo                                                 | Principalmente rutas SSH/cripto de drivers incluidos; alcanzabilidad no demostrada en PostgreSQL                                                                                                        | 0.55.0 cubre los fixes listados                                              | Módulo Goose propuesto                                                                                                  |
| golang.org/x/net (Goose)              | 0.33.0                   | 7 CVE, lista abajo                                                  | HTML (XSS/DoS), HTTP/2, DNS, IDNA y proxies; presencia en binario, sin prueba por ruta de ejecución                                                                                                     | 0.62.0 cubre los fixes listados                                              | Módulo Goose propuesto                                                                                                  |
| google.golang.org/grpc (Goose)        | 1.62.1                   | CVE-2026-33186, CVE-2026-84304, CVE-2026-84445, GHSA-hrxh-6v49-42gf | DoS y autorización RBAC en drivers extra; no se levanta servidor gRPC en esta migración                                                                                                                 | 1.83.2 cubre los cuatro; 1.82.2 no cubre el fix 1.83.1 de CVE-2026-84304     | Módulo Goose propuesto                                                                                                  |

Los **18 de govulncheck** son **16 avisos de stdlib + 2 de un único módulo x/text**,
no 18 dependencias. La biblioteca estándar queda incorporada en los binarios: actualizar
solo la instalación local de Go sin reconstruir imágenes no remedia los binarios existentes.
Los avisos adicionales de Goose los encontró Trivy en la imagen; no se suman a los 18 de
la API como si fueran un solo análisis. Trivy detecta presencia/versiones, no alcanzabilidad.

`pnpm audit` terminó con **código 1: 1 HIGH, 6 moderate, 1 low**. Se conserva el bloqueo HIGH
y también se documentan los inferiores; no hay ignore ni desactivación. El control Trivy
usa HIGH/CRITICAL con corrección disponible, tal como exige el DAS; esta tabla no certifica
la ausencia de avisos de otras severidades en imágenes.

## Propuesta mínima para el jefe

1. Autorizar Go 1.26.5 → 1.26.10 en toolchain, CI y builds; imagen candidata verificada:
   `golang:1.26.10-alpine3.23@sha256:78e7f1e419545bc4b8b5ee8a40a99ddb5df1216d38681dbab2eeb1fdb2427f43`.
2. Autorizar x/text de la API a 0.41.0 y Vite a 8.1.6. Regenerar únicamente los locks
   afectados, revisar cambios transitivos y repetir verify, e2e, govulncheck, audit y Trivy.
3. Para los avisos moderados de tests, autorizar Vitest y @vitest/coverage-v8 alineados en
   4.1.12, arrastrando @vitest/mocker 4.1.12. Son herramientas de pruebas, no runtime.
4. Mantener Goose 3.24.1 pero proponer un módulo reproducible `tools/goose/go.mod` y
   `tools/goose/go.sum` que fije los mínimos corregidos de su tabla. Modificar build y e2e
   para usar ese mismo grafo. El actual `go install/run ...@v3.24.1` usa el grafo publicado
   de Goose y **no** hereda las correcciones de apps/api/go.mod. No basta editar el módulo
   de la API. Estos son mínimos por aviso: Go MVS puede elevar versiones transitivas;
   ese grafo debe revisarse y probarse después de la autorización. No se propone una
   versión de Goose supuestamente segura sin haberla verificado.
5. Jefe: registrar la excepción/actualización del DAS; no editar su referencia unilateralmente.
   Reconstruir las tres imágenes y repetir Trivy/Dockle. No prometer verde antes de repetir.

Los módulos Go propuestos publican requisitos Go <=1.26.0, compatibles a nivel de versión
con 1.26.10; eso no sustituye compilación y pruebas. El mínimo Vite publicado acepta Node
22.22.0. No se aplicó ninguna de estas propuestas.

## Identificadores completos por dependencia en imágenes

- **golang.org/x/text**: CVE-2026-56851, CVE-2026-56852.
- **stdlib**: CVE-2026-111943, CVE-2026-33818, CVE-2026-39821, CVE-2026-46600, CVE-2026-56853, CVE-2026-56858, CVE-2026-56859, CVE-2026-56860, CVE-2026-56862, CVE-2026-78667, CVE-2026-78669, CVE-2026-97031.

## Fuentes oficiales

- [Aviso HIGH de Vite](https://github.com/advisories/GHSA-587h-7wmj-h5pm).
- [Go 1.26.10](https://go.dev/dl/#go1.26.10) y [base oficial de vulnerabilidades Go](https://vuln.go.dev/).
- Cada aviso de imagen y su rango/fix se conserva con URL oficial en los JSON citados. Los dos avisos sin ficha GitHub se contrastaron con GO-2026-4918 y GO-2026-6611.

## Los 18 avisos alcanzables de govulncheck

- [GO-2026-6809](https://vuln.go.dev/ID/GO-2026-6809.json): Distrust compromised intermediate certificates in crypto/x509.
- [GO-2026-6629](https://vuln.go.dev/ID/GO-2026-6629.json): Panic parsing crafted input in x/text/secure/precis in golang.org/x/text.
- [GO-2026-6617](https://vuln.go.dev/ID/GO-2026-6617.json): HTTP/2 server crash due to HPACK encoder race in net/http.
- [GO-2026-6613](https://vuln.go.dev/ID/GO-2026-6613.json): HTTP/1 server connection desynchronization after 2xx CONNECT response in net/http.
- [GO-2026-6612](https://vuln.go.dev/ID/GO-2026-6612.json): Double flow control refund on HTTP/2 server streams in net/http.
- [GO-2026-6611](https://vuln.go.dev/ID/GO-2026-6611.json): Excessive CPU consumption from repeated initial window changes in net/http.
- [GO-2026-6610](https://vuln.go.dev/ID/GO-2026-6610.json): HTTP/2 transport accepts malformed framing-related headers in net/http.
- [GO-2026-6609](https://vuln.go.dev/ID/GO-2026-6609.json): Lack of limit on size of parsed Range headers in net/http.
- [GO-2026-6608](https://vuln.go.dev/ID/GO-2026-6608.json): Memory limit bypass when parsing MIME headers in net/textproto, mime/multipart.
- [GO-2026-6607](https://vuln.go.dev/ID/GO-2026-6607.json): Reject malformed ECH outer extension references in crypto/tls.
- [GO-2026-6605](https://vuln.go.dev/ID/GO-2026-6605.json): HTTP/1 client connection desynchronization after CONNECT rejection in net/http.
- [GO-2026-6603](https://vuln.go.dev/ID/GO-2026-6603.json): HTTP/2 server memory exhaustion due to Trailer headers in net/http.
- [GO-2026-6218](https://vuln.go.dev/ID/GO-2026-6218.json): Avoid quadratic complexity in resolvePath in net/url.
- [GO-2026-6090](https://vuln.go.dev/ID/GO-2026-6090.json): Limit handshake messages we are willing to accept post-handshake in crypto/tls.
- [GO-2026-6089](https://vuln.go.dev/ID/GO-2026-6089.json): Apply ReadHeaderTimeout when doing unencrypted HTTP/2 check in net/http.
- [GO-2026-5972](https://vuln.go.dev/ID/GO-2026-5972.json): Enforce maximum recursion depth in encoding/asn1.
- [GO-2026-5970](https://vuln.go.dev/ID/GO-2026-5970.json): Infinite loop on invalid input in golang.org/x/text.
- [GO-2026-5026](https://vuln.go.dev/ID/GO-2026-5026.json): Invoking failure to reject ASCII-only Punycode-encoded labels in golang.org/x/net/idna.

## Decisión pendiente tras las correcciones sin cambio de versiones

Las filas de la tabla anterior siguen pendientes de aprobación del jefe. Las correcciones
del servidor SPA y sus pruebas no modifican Go, módulos, catalog pnpm ni lockfiles. La
aprobación debe cubrir por separado el toolchain/stdlib, x/text de la API, Vite, la familia
Vitest y el grafo propio de Goose. Las versiones corregidas y fuentes verificadas se mantienen
en la tabla y en sus evidencias; no se afirma que el grafo propuesto haya sido instalado o probado.
Tras la aprobación corresponde regenerar los archivos indicados, reconstruir imágenes y
repetir controles. Mientras tanto `seguridad` debe continuar bloqueando los hallazgos.
