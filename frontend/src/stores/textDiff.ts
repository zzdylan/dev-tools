import { defineStore } from 'pinia'

export const useTextDiffStore = defineStore('textDiff', {
  state: () => ({
    oldText: '',
    newText: '',
    ignoreWhitespace: false
  }),
  
  persist: {
    enabled: true,
    strategies: [
      {
        key: 'text-diff-store',
        storage: localStorage,
      },
    ],
  },
}) 