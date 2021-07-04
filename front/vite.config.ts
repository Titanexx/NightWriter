import { defineConfig } from 'vite'
import { resolve } from "path";
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [vue()],
    resolve: {
        alias: {
            '@': resolve(__dirname, './src'),
            "@routes": resolve(__dirname, './src/routes'),
            "@modules": resolve(__dirname, './src/modules'),
            "@plugins": resolve(__dirname, './src/plugins'),
            "@components": resolve(__dirname, './src/components'),
            "@layouts": resolve(__dirname, './src/components/layouts'),
            "@forms": resolve(__dirname, './src/components/forms'),
            "@views": resolve(__dirname, './src/views'),
        }
    },
    server: {
        proxy:{
            "/api":"http://localhost:8080/",
            "/ws":{
                target: "ws://localhost:8080/",
                ws: true,
            }
        }
    },
    build:{
        terserOptions:{
        }
    }
})
