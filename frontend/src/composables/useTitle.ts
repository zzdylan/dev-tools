import { inject } from 'vue'

export function useTitle(title: string) {
  const setTitle = inject('setTitle') as (title: string) => void
  setTitle(title)
} 