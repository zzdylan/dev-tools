import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'
import MainLayout from '../layouts/MainLayout.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: MainLayout,
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('../views/Home.vue')
      },
      {
        path: 'json-editor',
        name: 'JsonEditor',
        component: () => import('../views/JsonEditor.vue')
      },
      {
        path: 'xml-editor',
        name: 'XmlEditor',
        component: () => import('../views/XmlEditor.vue')
      },
      {
        path: 'tools',
        name: 'Tools',
        component: () => import('../views/Tools.vue')
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('../views/Settings.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router 