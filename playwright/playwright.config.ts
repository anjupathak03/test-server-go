import { defineConfig } from "@playwright/test";

const baseURL = process.env.BASE_URL || "http://localhost:8080";
const readinessURL = `${baseURL.replace(/\/+$/, "")}/api/todos`;
const baseURLPort = (() => {
  try {
    return new URL(baseURL).port || "8080";
  } catch {
    return "8080";
  }
})();

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
    baseURL,
    extraHTTPHeaders: {
      "Content-Type": "application/json",
    },
  },
  webServer: {
    command: "go run .",
    cwd: "..",
    url: readinessURL,
    reuseExistingServer: !process.env.CI,
    timeout: 120_000,
    env: {
      ...process.env,
      DB_HOST: process.env.DB_HOST || "localhost",
      DB_PORT: process.env.DB_PORT || "3306",
      DB_USER: process.env.DB_USER || "root",
      DB_PASSWORD: process.env.DB_PASSWORD || "password",
      DB_NAME: process.env.DB_NAME || "todo_db",
      PORT: process.env.PORT || baseURLPort,
    },
  },
});
