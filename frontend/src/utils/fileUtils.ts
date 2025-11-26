import { SaveFile } from '../../wailsjs/go/app/App'
import { file } from '../../wailsjs/go/models'

// 使用 Wails 生成的类型
export type FileFilter = file.Filter
export type SaveOptions = file.SaveOptions

// 常用文件过滤器
export const FILE_FILTERS = {
  IMAGE_PNG: {
    displayName: 'PNG图片 (*.png)',
    pattern: '*.png'
  },
  IMAGE_JPG: {
    displayName: 'JPEG图片 (*.jpg)',
    pattern: '*.jpg'
  },
  IMAGE_ALL: {
    displayName: '图片文件 (*.png;*.jpg;*.gif)',
    pattern: '*.png;*.jpg;*.gif'
  },
  TEXT: {
    displayName: '文本文件 (*.txt)',
    pattern: '*.txt'
  },
  JSON: {
    displayName: 'JSON文件 (*.json)',
    pattern: '*.json'
  },
  XML: {
    displayName: 'XML文件 (*.xml)',
    pattern: '*.xml'
  },
  CSV: {
    displayName: 'CSV文件 (*.csv)',
    pattern: '*.csv'
  },
  ALL: {
    displayName: '所有文件 (*.*)',
    pattern: '*.*'
  }
}

// 保存图片（base64格式）
export async function saveImage(
  base64Data: string,
  filename: string = 'image.png',
  title: string = '保存图片'
): Promise<string> {
  const filters = [
    new file.Filter(FILE_FILTERS.IMAGE_PNG),
    new file.Filter(FILE_FILTERS.IMAGE_JPG),
    new file.Filter(FILE_FILTERS.IMAGE_ALL),
    new file.Filter(FILE_FILTERS.ALL)
  ]

  const options = new file.SaveOptions({
    title,
    defaultFilename: filename,
    filters
  })

  try {
    const savedPath = await SaveFile(base64Data, options, true)
    return savedPath
  } catch (error) {
    console.error('保存图片失败:', error)
    throw error
  }
}

// 保存文本文件
export async function saveText(
  content: string,
  filename: string = 'file.txt',
  title: string = '保存文件',
  customFilters: any[] = [FILE_FILTERS.TEXT, FILE_FILTERS.ALL]
): Promise<string> {
  const filters = customFilters.map(filter => new file.Filter(filter))

  const options = new file.SaveOptions({
    title,
    defaultFilename: filename,
    filters
  })

  try {
    const savedPath = await SaveFile(content, options, false)
    return savedPath
  } catch (error) {
    console.error('保存文本失败:', error)
    throw error
  }
}

// 保存JSON文件
export async function saveJSON(
  data: any,
  filename: string = 'data.json',
  title: string = '保存JSON文件'
): Promise<string> {
  const content = JSON.stringify(data, null, 2)
  return saveText(content, filename, title, [
    FILE_FILTERS.JSON,
    FILE_FILTERS.TEXT,
    FILE_FILTERS.ALL
  ])
}

// 保存XML文件
export async function saveXML(
  content: string,
  filename: string = 'data.xml',
  title: string = '保存XML文件'
): Promise<string> {
  return saveText(content, filename, title, [
    FILE_FILTERS.XML,
    FILE_FILTERS.TEXT,
    FILE_FILTERS.ALL
  ])
}

// 保存CSV文件
export async function saveCSV(
  content: string,
  filename: string = 'data.csv',
  title: string = '保存CSV文件'
): Promise<string> {
  return saveText(content, filename, title, [
    FILE_FILTERS.CSV,
    FILE_FILTERS.TEXT,
    FILE_FILTERS.ALL
  ])
}