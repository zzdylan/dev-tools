declare module 'json-to-go' {
  function jsonToGo(json: string): { go: string };
  export default jsonToGo;
}

declare module 'gofmt.js' {
  function gofmt(code: string): string;
  export default gofmt;
} 