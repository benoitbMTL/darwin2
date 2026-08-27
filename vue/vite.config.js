import { defineConfig, loadEnv } from "vite";
import vue from "@vitejs/plugin-vue";

const backendRoutes = [
  "/api",
  "/apply-config",
  "/assign-vip-to-virtual-server",
  "/bot-deception",
  "/bot-page-source",
  "/bot-scraper-api",
  "/clone-config",
  "/clone-inline-protection",
  "/clone-signature-protection",
  "/config",
  "/configure-protection-profile",
  "/cookie-security",
  "/create-member-pool",
  "/create-policy",
  "/create-server-pool",
  "/create-virtual-ip",
  "/create-virtual-server",
  "/create-x-forwarded-for-rule",
  "/delete-inline-protection",
  "/delete-local",
  "/delete-policy",
  "/delete-server-pool",
  "/delete-signature-protection",
  "/delete-virtual-ip",
  "/delete-virtual-server",
  "/delete-x-forwarded-for-rule",
  "/import",
  "/known-bots",
  "/list-configs",
  "/machine-learning",
  "/rename-config",
  "/reset-api-machine-learning",
  "/reset-config",
  "/reset-machine-learning",
  "/run-health-check",
  "/save-config",
  "/selenium",
  "/traffic-generation",
  "/user-auth",
  "/web-attacks",
  "/web-scan",
];

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, ".", "");
  const proxy = Object.fromEntries(
    backendRoutes.map((route) => [
      route,
      {
        target: env.VITE_BACKEND_URL,
        changeOrigin: true,
      },
    ]),
  );

  return {
    plugins: [vue()],
    server: {
      host: "0.0.0.0",
      proxy,
    },
  };
});
