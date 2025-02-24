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

interface MenuItem {
  id: string
  path: string
  icon: string
  title: string
  visible: boolean
  order: number
  description: string
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
    },
    menuConfig: {
      items: <MenuItem[]>[
        { id: 'json', path: '/json-editor', icon: '{ }', title: 'JSON 编辑器', visible: true, order: 0, description: 'JSON 格式化、验证、转换工具' },
        { id: 'xml', path: '/xml-editor', icon: '📄', title: 'XML 编辑器', visible: true, order: 1, description: 'XML 格式化、验证、转换工具' },
        { id: 'time', path: '/time-converter', icon: '⏰', title: '时间戳转换', visible: true, order: 2, description: '将时间戳转换为日期时间或反之' },
        { id: 'url-encode', path: '/url-converter', icon: '🔗', title: 'URL 编解码', visible: true, order: 3, description: 'URL 编码和解码工具' },
        { id: 'url-parse', path: '/url-parser', icon: '🔍', title: 'URL 解析', visible: true, order: 4, description: '解析和提取URL中的信息' },
        { id: 'qrcode', path: '/qrcode', icon: '📱', title: '二维码工具', visible: true, order: 5, description: '生成和扫描二维码' },
        { id: 'base64-img', path: '/base64-image', icon: '🖼️', title: 'Base64 图像', visible: true, order: 6, description: '将图像转换为Base64编码' },
        { id: 'base64-text', path: '/base64-text', icon: '📝', title: 'Base64 文本', visible: true, order: 7, description: '将文本转换为Base64编码' },
        { id: 'number', path: '/number-converter', icon: '🔢', title: '进制转换', visible: true, order: 8, description: '不同进制之间的数值转换' },
        { id: 'diff', path: '/text-diff', icon: '📋', title: '文本对比', visible: true, order: 9, description: '比较两个文本的差异' },
        { id: 'curl', path: '/curl-converter', icon: '🔄', title: 'cURL 转换', visible: true, order: 10, description: '将cURL命令转换为其他语言' },
        { id: 'unicode', path: '/unicode-converter', icon: '🔤', title: 'Unicode 转换', visible: true, order: 11, description: 'Unicode编码和解码' }
      ]
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
    },
    updateMenuItemVisibility(id: string, visible: boolean) {
      const item = this.menuConfig.items.find(item => item.id === id)
      if (item) {
        item.visible = visible
      }
    },
    updateMenuOrder(items: MenuItem[]) {
      items.forEach((item, index) => {
        const existingItem = this.menuConfig.items.find(i => i.id === item.id)
        if (existingItem) {
          existingItem.order = index
        }
      })
    }
  },
  persist: true
}) 