import { defineConfig } from 'vite';
import { resolve } from 'path';

export default defineConfig({
  root: 'src-web',
  base: '/web/',
  build: {
    outDir: '../web',
    emptyOutDir: true,
    rollupOptions: {
      input: {
        index: resolve(__dirname, 'src-web/index.html'),
        login: resolve(__dirname, 'src-web/login.html'),
        'pages/dashboard': resolve(__dirname, 'src-web/pages/dashboard.html'),
        'pages/configuracion': resolve(__dirname, 'src-web/pages/configuracion.html'),
        'pages/usuarios': resolve(__dirname, 'src-web/pages/usuarios.html'),
      },
      output: {
        entryFileNames: 'js/[name]-[hash].js',
        chunkFileNames: 'js/[name]-[hash].js',
        assetFileNames: ({ name }) => {
          if (/\.(css)$/.test(name ?? '')) {
            return 'css/[name]-[hash][extname]';
          }
          return 'assets/[name]-[hash][extname]';
        },
      },
    },
  },
});
