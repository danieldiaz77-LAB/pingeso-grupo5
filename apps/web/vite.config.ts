import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vitest/config";

const apiLocal = "http://localhost:5171";

export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    port: 5170,
    strictPort: true,
    proxy: {
      "/api": { target: apiLocal, changeOrigin: true },
      "/healthz": { target: apiLocal, changeOrigin: true },
      "/readyz": { target: apiLocal, changeOrigin: true },
    },
  },
  preview: {
    port: 5170,
    strictPort: true,
  },
  test: {
    environment: "node",
    include: ["src/**/*.test.ts"],
  },
});
