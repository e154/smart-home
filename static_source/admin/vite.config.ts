import {resolve} from 'path'
import type {ConfigEnv, UserConfig} from 'vite'
import {loadEnv} from 'vite'
import Vue from '@vitejs/plugin-vue'
import VueJsx from '@vitejs/plugin-vue-jsx'
import WindiCSS from 'vite-plugin-windicss'
import progress from 'vite-plugin-progress'
import EslintPlugin from 'vite-plugin-eslint'
import {ViteEjsPlugin} from "vite-plugin-ejs"
import PurgeIcons from 'vite-plugin-purge-icons'
import VueI18nPlugin from "@intlify/unplugin-vue-i18n/vite"
import {createSvgIconsPlugin} from 'vite-plugin-svg-icons'
import DefineOptions from "unplugin-vue-define-options/vite"
import {createStyleImportPlugin, ElementPlusResolve} from 'vite-plugin-style-import'
import Unfonts from 'unplugin-fonts/vite'
import analyze from "rollup-plugin-analyzer";
import {VitePWA, VitePWAOptions} from 'vite-plugin-pwa'
import process from "node:process";
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vitejs.dev/config/
const root = process.cwd()

function pathResolve(dir: string) {
  return resolve(root, '.', dir)
}

const pwaOptions: Partial<VitePWAOptions> = {
  mode: 'production',
  base: '/',
  includeAssets: ['*.svg', '*.png', '*.xml', '*.ico'],
  manifest: {
    id: '36b70975-9daf-4ea0-a451-340ab66fc175',
    orientation: 'any',
    name: "Smart Home",
    short_name: "Smart Home",
    description: "Software package for automation",
    start_url: "/",
    display: "standalone",
    background_color: "#333335",
    theme_color: "#333335",
    icons: [
      {
        "src": "/android-chrome-64x64.png",
        "type": "image/png",
        "sizes": "64x64"
      },
      {
        "src": "/android-chrome-192x192.png",
        "type": "image/png",
        "sizes": "192x192"
      },
      {
        src: '/android-chrome-512x512.png',
        sizes: '512x512',
        type: 'image/png',
      },
      {
        src: '/maskable-icon.png',
        sizes: '512x512',
        type: 'image/png',
        purpose: 'maskable'
      }
    ],
  },

  registerType: 'autoUpdate',

  strategies: 'injectManifest',
  injectManifest: {
    rollupFormat: 'iife',
    maximumFileSizeToCacheInBytes: 5 * 1024 * 1024,
  },

  // add this to cache all the imports
  workbox: {
    maximumFileSizeToCacheInBytes: 5 * 1024 * 1024,
  },
}

export default ({command, mode}: ConfigEnv): UserConfig => {
  let env = {} as any
  const isBuild = command === 'build'
  if (!isBuild) {
    env = loadEnv((process.argv[3] === '--mode' ? process.argv[4] : process.argv[3]), root)
  } else {
    env = loadEnv(mode, root)
  }
  return {
    base: env.VITE_BASE_PATH,
    plugins: [
      Vue(),
      vueDevTools(),
      VueJsx(),
      WindiCSS(),
      progress(),
      Unfonts({}),
      createStyleImportPlugin({
        resolves: [ElementPlusResolve()],
        libs: [{
          libraryName: 'element-plus',
          esModule: true,
          resolveStyle: (name) => {
            return `element-plus/es/components/${name.substring(3)}/style/css`
          }
        }]
      }),
      EslintPlugin({
        cache: false,
        include: ['src/**/*.vue', 'src/**/*.ts', 'src/**/*.tsx'] // 检查的文件
      }),
      VueI18nPlugin({
        runtimeOnly: true,
        compositionOnly: true,
        include: [resolve(__dirname, 'src/locales/**')]
      }),
      createSvgIconsPlugin({
        iconDirs: [pathResolve('src/assets/svgs')],
        symbolId: 'icon-[dir]-[name]',
        svgoOptions: true
      }),
      PurgeIcons(),
      DefineOptions(),
      ViteEjsPlugin({
        title: env.VITE_APP_TITLE
      }),
      // virtualMessagePlugin(),
      VitePWA(pwaOptions),
    ],

    css: {
      preprocessorOptions: {
        less: {
          additionalData: '@import "./src/styles/variables.module.less";',
          javascriptEnabled: true
        }
      }
    },
    resolve: {
      extensions: ['.mjs', '.js', '.ts', '.jsx', '.tsx', '.json', '.less', '.css'],
      alias: [
        {
          find: 'vue-i18n',
          replacement: 'vue-i18n/dist/vue-i18n.cjs.js'
        },
        {
          find: /\@\//,
          replacement: `${pathResolve('src')}/`
        }
      ]
    },
    build: {
      minify: 'terser',
      outDir: env.VITE_OUT_DIR || 'dist',
      sourcemap: env.VITE_SOURCEMAP === 'true' ? 'inline' : false,
      // brotliSize: false,
      terserOptions: {
        compress: {
          drop_debugger: env.VITE_DROP_DEBUGGER === 'true',
          drop_console: env.VITE_DROP_CONSOLE === 'true'
        }
      },
      rollupOptions: {
        plugins: [analyze()]
      },
    },
    server: {
      port: 9527,
      proxy: {
        // 选项写法
        '/api': {
          target: 'http://127.0.0.1:8000',
          changeOrigin: true,
          rewrite: path => path.replace(/^\/api/, '')
        }
      },
      hmr: {
        overlay: false
      },
      host: '0.0.0.0'
    },
    optimizeDeps: {
      include: [
        'vue',
        'vue-router',
        'vue-types',
        'element-plus/es/locale/lang/zh-cn',
        'element-plus/es/locale/lang/en',
        '@iconify/iconify',
        '@vueuse/core',
        'axios',
        'qs',
        'echarts',
        'echarts-wordcloud',
        'intro.js',
        'qrcode',
      ]
    },

  }
}
