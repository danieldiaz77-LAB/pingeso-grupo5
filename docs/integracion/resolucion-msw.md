# Resolución de MSW en la integración local

La ejecución manual de `pnpm install --lockfile-only --ignore-scripts` con Node 22.22.0 y pnpm 10.13.1 falló: MSW 2.15.0 resolvió @inquirer/confirm 6.4.0, que exige ^22.22.2 || ^24.15.0 || >=26.0.0. No se desactivó engine-strict ni se cambió Node.

Se retienen mediante overrides las versiones existentes en el lockfile de feat/plantilla-hito0, commit 7d2b7dcc216eec847b78b7d7aeb4731aec1cd94d:

- `@inquirer/ansi@2.0.10`
- `@inquirer/confirm@6.3.5`
- `@inquirer/core@12.0.6`
- `@inquirer/figures@2.0.10`
- `@inquirer/type@4.2.1`

Esto no acredita instalación ni seguridad. Pendiente: regenerar el lockfile con pnpm, revisar todas las diferencias de versiones, instalación frozen-lockfile, audit y pruebas. No se ha editado manualmente el lockfile ni adoptado las actualizaciones pendientes del DAS.

La primera regeneración con overrides terminó correctamente, con un aviso peer de openapi-typescript 7.13.0 (^5.x) frente a TypeScript 6.0.3. Conservó las 171 versiones del pipeline y añadió 98 entradas. Nueve resoluciones nuevas diferían de la plantilla (ver revision-lock.json); se fijaron a las versiones de Benja y queda pendiente regenerar de nuevo. No hay instalación frozen ni pruebas acreditadas todavía.
