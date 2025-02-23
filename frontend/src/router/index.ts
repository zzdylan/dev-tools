import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'
import MainLayout from '../layouts/MainLayout.vue'
import TimeConverter from '../views/TimeConverter.vue'
import UrlConverter from '../views/UrlConverter.vue'
import UrlParser from '../views/UrlParser.vue'
import QrCode from '../views/QrCode.vue'
import Base64Image from '../views/Base64Image.vue'
import Base64Text from '../views/Base64Text.vue'
import NumberConverter from '../views/NumberConverter.vue'
import TextDiff from '../views/TextDiff.vue'
import CurlConverter from '../views/CurlConverter.vue'
import UnicodeConverter from '../views/UnicodeConverter.vue'

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
      // {
      //   path: 'tools',
      //   name: 'Tools',
      //   component: () => import('../views/Tools.vue')
      // },
      // {
      //   path: 'settings',
      //   name: 'Settings',
      //   component: () => import('../views/Settings.vue')
      // },
      {
        path: '/time-converter',
        name: 'TimeConverter',
        component: TimeConverter
      },
      {
        path: '/url-converter',
        name: 'UrlConverter',
        component: UrlConverter
      },
      {
        path: '/url-parser',
        name: 'UrlParser',
        component: UrlParser
      },
      {
        path: '/qrcode',
        name: 'QrCode',
        component: QrCode
      },
      {
        path: '/base64-image',
        name: 'Base64Image',
        component: Base64Image
      },
      {
        path: '/base64-text',
        name: 'Base64Text',
        component: Base64Text
      },
      {
        path: '/number-converter',
        name: 'NumberConverter',
        component: NumberConverter
      },
      {
        path: '/text-diff',
        name: 'TextDiff',
        component: TextDiff
      },
      {
        path: '/curl-converter',
        name: 'CurlConverter',
        component: CurlConverter
      },
      {
        path: '/unicode-converter',
        name: 'UnicodeConverter',
        component: UnicodeConverter
      }
    ]
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router 