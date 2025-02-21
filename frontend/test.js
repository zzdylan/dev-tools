// 测试数据
const testJson = {
  user: {
    id: 123,
    name: "John Doe",
    email: "john.doe@example.com"
  },
  products: [
    {
      id: "p1",
      name: "Product A",
      price: 19.99
    },
    {
      id: "p2",
      name: "Product B",
      price: 29.99
    },
    {
      id: "p3",
      name: "Product \u4E2D\u6587",
      price: 39.99
    }
  ],
  order: {
    orderId: "abc123",
    date: "2023-08-18",
    items: [
      {
        productId: "p1",
        quantity: 2
      },
      {
        productId: "p3",
        quantity: 1
      }
    ]
  }
}

// 辅助函数：将字符串中的非 ASCII 字符转换为 Unicode 转义序列
function encodeUnicode(str) {
  return str.replace(/[\u007f-\uffff]/g, char => {
    const code = char.charCodeAt(0).toString(16).toUpperCase()
    return '\\\\u' + ('0000' + code).slice(-4)  // 使用双反斜杠
  })
}

// 模拟格式化实现
function formatJson(json, autoDecodeUnicode = false) {
  if (autoDecodeUnicode) {
    // 解码时，直接使用 JSON.stringify
    return JSON.stringify(json, null, 4)
  } else {
    // 不解码时，先转成字符串再解析，这样可以保持原有的转义序列
    const str = JSON.stringify(json)
    const parsed = JSON.parse(str)
    return JSON.stringify(parsed, null, 4)
  }
}

// 模拟压缩实现
function compressJson(str, autoDecodeUnicode = false) {
  let parsed = JSON.parse(str)
  
  if (autoDecodeUnicode) {
    // 解码时，直接使用 JSON.stringify
    return JSON.stringify(parsed)
  } else {
    // 不解码时，先转成字符串再解析，这样可以保持原有的转义序列
    const temp = JSON.stringify(parsed)
    return JSON.stringify(JSON.parse(temp))
  }
}

// 运行测试
console.log('=== 不解码 Unicode ===')
console.log(formatJson(testJson, false))

console.log('\n=== 解码 Unicode ===')
console.log(formatJson(testJson, true))

console.log('\n=== 压缩测试 ===')
console.log(compressJson(formatJson(testJson, false)))
console.log(compressJson(formatJson(testJson, true))) 