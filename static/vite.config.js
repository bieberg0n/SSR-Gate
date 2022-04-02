import { createVuePlugin } from "vite-plugin-vue2";
import { defineConfig } from "vite";
import path from "path";

export default defineConfig({
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "src")
    }
  },
  base: "/",
  plugins: [
    // vue()
    createVuePlugin()
  ],
  server: {
    proxy: {
      '/api': 'http://192.168.1.10:5000'
    }
  }
});
