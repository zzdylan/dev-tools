import { defineStore } from 'pinia'

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
    xmlEditor: {
      code: ''
    },
    jsonEditor: {
      code: '',
      settings: {
        autoDecodeUnicode: false,
        removeEscapes: false
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
    }
  }),
  persist: true
}) 