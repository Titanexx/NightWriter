import { defineConfig, loadEnv } from 'vite'
import { resolve } from "path";
import vue from '@vitejs/plugin-vue'

// VITE_SERVER_DOMAIN = "localhost"
// VITE_SERVER_PORT = "8081"
// VITE_HTTPS = false

// https://vitejs.dev/config/
export default ({ mode }) => {
    process.env = {...process.env, ...loadEnv(mode, process.cwd())};
    return defineConfig({
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
                "/api": process.env.VITE_HTTPS == "true" ? `https://${process.env.VITE_SERVER_DOMAIN}:${process.env.VITE_SERVER_PORT}/` : `http://${process.env.VITE_SERVER_DOMAIN}:${process.env.VITE_SERVER_PORT}/`,
                "/ws":{
                    target: process.env.VITE_HTTPS == "true" ? `wss://${process.env.VITE_SERVER_DOMAIN}:${process.env.VITE_SERVER_PORT}/` : `ws://${process.env.VITE_SERVER_DOMAIN}:${process.env.VITE_SERVER_PORT}/`,
                    ws: true,
                }
            }
        },
        build:{
            outDir: 'dist',
            terserOptions:{
            }
        }
    });
}