/// <reference types="vitest" />
import { defineConfig } from "vitest/config";

export default defineConfig({
  test: {
    globals: true,
    environment: "jsdom",
    setupFiles: ["./setup.ts"],
    include: ["specs/**/*.spec.ts"],
    deps: {
      inline: [/@angular/, /@openfeature/],
    },
    coverage: {
      provider: "v8",
      reporter: ["text", "json", "html"],
    },
  },
});
