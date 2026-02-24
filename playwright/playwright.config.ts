import { defineConfig } from "@playwright/test";

export default defineConfig({
  testDir: "./tests",
  timeout: 30_000,
  retries: 0,
  reporter: [["list"], ["html", { open: "never" }]],
  // Pure API tests – no browser required
  projects: [
    {
      name: "api",
    },
  ],
  use: {
    baseURL: process.env.BASE_URL || "http://localhost:8080",
    extraHTTPHeaders: {
      "Content-Type": "application/json",
    },
  },
});
