import { defineStore } from 'pinia'

interface JsonEditorTab {
  code: string
  settings: {
    autoDecodeUnicode: boolean
    removeEscapes: boolean
  }
}

interface JsonEditorTabs {
  [key: string]: JsonEditorTab
}

interface XmlEditorTab {
  code: string
}

interface XmlEditorTabs {
  [key: string]: XmlEditorTab
}

export const useToolsStore = defineStore('tools', {
  state: () => ({
    // 保持每个组件原有的 ref 变量名
    qrCode: {
      text: 'I am the qrcode generator of devTools😁'
    },
    textDiff: {
      oldText: '',
      newText: '',
      ignoreWhitespace: false
    },
    xmlEditorTabs: <XmlEditorTabs>{
      default: {
        code: ''
      }
    },
    jsonEditorTabs: <JsonEditorTabs>{
      default: {
        code: '',
        settings: {
          autoDecodeUnicode: false,
          removeEscapes: false
        }
      }
    },
    urlConverter: {
      rawText: '',
      result: ''
    },
    urlParser: {
      urlText: '',
      viewMode: 'json'
    },
    base64Text: {
      rawText: '',
      result: ''
    },
    base64Image: {
      imageUrl: '',
      base64Result: ''
    },
    numberConverter: {
      inputValue: '',
      inputBase: '10'
    },
    timeConverter: {
      timestamp: '',
      datetime: ''
    },
    curlConverter: {
      curlCommand: '',
      targetLang: 'javascript', // 默认转换为 JavaScript
      convertedCode: ''
    },
    unicodeConverter: {
      rawText: '',
      result: '',
      mode: 'encode' // 'encode' 或 'decode'
    }
  }),
  actions: {
    createJsonEditorTab() {
      const id = Date.now().toString()
      this.jsonEditorTabs[id] = {
        code: '',
        settings: {
          autoDecodeUnicode: false,
          removeEscapes: false
        }
      }
      return id
    },
    createXmlEditorTab() {
      const id = Date.now().toString()
      this.xmlEditorTabs[id] = {
        code: ''
      }
      return id
    }
  },
  persist: true
}) 