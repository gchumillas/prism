export {}

declare global {
  interface Window {
    wasmReady: Promise<void>
    prism: {
      syncCanvas: (
        buf: ImageDataArray,
        refreshCb: (buf: Uint8Array, x0: int, y0: int, x1: int, y1: int) => void
      ) => void
      tapCanvas: (x: number, y: number) => void
    }
  }
}
