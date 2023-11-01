import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react-swc";
import { fileURLToPath } from "url";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  test: {
    globals: true,
    environment: "jsdom",
  },
  server: {
    host: "0.0.0.0",
    port: 3000,
  },
  css: {
    preprocessorOptions: {
      scss: {},
    },
  },
  resolve: {
    alias: [
      {
        find: "@",
        replacement: fileURLToPath(new URL("./src", import.meta.url)),
      },

      {
        find: "@assets",
        replacement: fileURLToPath(new URL("./src/assets", import.meta.url)),
      },
      {
        find: "@components",
        replacement: fileURLToPath(
          new URL("./src/components", import.meta.url)
        ),
      },
      {
        find: "@config",
        replacement: fileURLToPath(new URL("./src/config", import.meta.url)),
      },
      {
        find: "@lib",
        replacement: fileURLToPath(new URL("./src/lib", import.meta.url)),
      },
      {
        find: "@features",
        replacement: fileURLToPath(new URL("./src/features", import.meta.url)),
      },
      {
        find: "@pages",
        replacement: fileURLToPath(new URL("./src/pages", import.meta.url)),
      },
      {
        find: "@types",
        replacement: fileURLToPath(new URL("./src/types", import.meta.url)),
      },
      {
        find: "@test",
        replacement: fileURLToPath(new URL("./src/test", import.meta.url)),
      },
      {
        find: "@hooks",
        replacement: fileURLToPath(new URL("./src/hooks", import.meta.url)),
      },
      {
        find: "@providers",
        replacement: fileURLToPath(new URL("./src/providers", import.meta.url)),
      },
      {
        find: "@utils",
        replacement: fileURLToPath(new URL("./src/utils", import.meta.url)),
      },
      {
        find: "@stores",
        replacement: fileURLToPath(new URL("./src/stores", import.meta.url)),
      },
      {
        find: "@types",
        replacement: fileURLToPath(new URL("./src/types", import.meta.url)),
      },
    ],
  },
});
