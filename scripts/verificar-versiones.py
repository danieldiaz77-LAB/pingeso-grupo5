"""Consulta fuentes oficiales; no modifica las versiones del DAS."""
import concurrent.futures
import json
import re
import urllib.request
from pathlib import Path

DAS = next(Path("docs/referencias").glob("DAS*.md")).read_text(encoding="utf-8")
consultas = []
for repo, sha, version in re.findall(r"uses: ([\w/-]+)@([a-f0-9]{40}) # (v[\d.]+)", DAS):
    consultas.append((repo, version, "https://api.github.com/repos/" + repo + "/git/ref/tags/" + version, sha))
paquetes = {"react":"19.2.8", "react-dom":"19.2.8", "typescript":"6.0.3", "vite":"8.1.5", "@vitejs/plugin-react":"6.0.4", "tailwindcss":"4.3.3", "@tailwindcss/vite":"4.3.3", "@base-ui/react":"1.6.0", "react-router":"8.3.0", "@tanstack/react-query":"5.101.4", "react-hook-form":"7.83.0", "openapi-fetch":"0.15.2", "openapi-typescript":"7.13.0", "msw":"2.15.0", "vitest":"4.1.10", "@vitest/coverage-v8":"4.1.10", "@playwright/test":"1.62.0", "oxlint":"1.75.0", "prettier":"3.9.6", "prettier-plugin-tailwindcss":"0.8.1", "lefthook":"2.1.10", "@stoplight/prism-cli":"5.16.0", "pnpm":"10.13.1", "@types/react":"19.2.2", "@types/react-dom":"19.2.2"}
for nombre, version in paquetes.items():
    consultas.append((nombre, version, f"https://registry.npmjs.org/{nombre}/{version}", None))
for modulo, version in re.findall(r"`(github.com/[^` ]+|golang.org/[^` ]+)` \| (v[\d.]+)", DAS):
    escaped = ''.join('!' + c.lower() if c.isupper() else c for c in modulo)
    consultas.append((modulo, version, f"https://proxy.golang.org/{escaped}/@v/{version}.info", None))
for repo, version in {"aquasecurity/trivy":"v0.76.1", "gitleaks/gitleaks":"v8.30.1", "hadolint/hadolint":"v2.15.1", "goodwithtech/dockle":"v0.4.15", "oasdiff/oasdiff":"v1.35.0", "rhysd/actionlint":"v1.7.7"}.items():
    consultas.append((repo, version, f"https://api.github.com/repos/{repo}/releases/tags/{version}", None))
consultas.append(("govulncheck", "v1.9.0", "https://proxy.golang.org/golang.org/x/vuln/@v/v1.9.0.info", None))

def verificar(item):
    nombre, version, url, esperado = item
    resultado = dict(nombre=nombre, version=version, fuente=url)
    try:
        with urllib.request.urlopen(url, timeout=45) as response:
            data = json.load(response)
        resultado["estado"] = "verificado"
        if esperado:
            obj = data["object"]
            if obj["type"] == "tag":
                with urllib.request.urlopen(obj["url"], timeout=45) as response:
                    obj = json.load(response)["object"]
            resultado["sha"] = obj["sha"]
            if obj["sha"] != esperado:
                resultado["estado"] = "SHA distinto al DAS"
        elif "dist" in data:
            resultado["integrity"] = data["dist"]["integrity"]
    except Exception as error:
        resultado["estado"] = str(error)
    return resultado

if __name__ == "__main__":
    with concurrent.futures.ThreadPoolExecutor(max_workers=8) as pool:
        resultados = list(pool.map(verificar, sorted(set(consultas))))
    Path("docs/evidencias").mkdir(exist_ok=True)
    Path("docs/evidencias/versiones.json").write_text(json.dumps(resultados, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    for r in resultados:
        print(r["nombre"], r["version"], r["estado"])
    if any(r["estado"] != "verificado" for r in resultados):
        raise SystemExit(1)
