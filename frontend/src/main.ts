import {createApp} from 'vue'
import App from './App.vue'
import './style.css';
import router from './router'
import * as monaco from 'monaco-editor'
import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'
import jsonWorker from 'monaco-editor/esm/vs/language/json/json.worker?worker'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

// 配置 Monaco Editor 的 workers
window.MonacoEnvironment = {
  getWorker(_, label) {
    if (label === 'json') {
      return new jsonWorker()
    }
    return new editorWorker()
  }
}

// 配置 JSON 格式化选项
monaco.languages.json.jsonDefaults.setDiagnosticsOptions({
  validate: true,
  allowComments: false,
  schemas: [],
})

const app = createApp(App)
app.use(ElementPlus)
app.use(router)
app.mount('#app')
