declare module 'jinrishici' {
  interface JinrishiciResult {
    data: {
      content: string
      origin: {
        title: string
        dynasty: string
        author: string
        content: string[]
      }
    }
  }

  export function load(
    success: (result: JinrishiciResult) => void,
    error?: (err: Error) => void
  ): void
}
