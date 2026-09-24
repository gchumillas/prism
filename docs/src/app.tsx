import { useEffect, useRef } from 'preact/hooks'

// IMPORTANT: these values must match those defined in internal/platform/web.go
const canvasWidth = 1404
const canvasHeight = 1872
const scale = 1

const App = () => {
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const { syncCanvas, tapCanvas } = window.prism

  const handleClick = (
    e: MouseEvent & { currentTarget: HTMLCanvasElement }
  ) => {
    const canvas = e.currentTarget
    const rect = canvas.getBoundingClientRect()
    tapCanvas((e.clientX - rect.left) / scale, (e.clientY - rect.top) / scale)
  }

  useEffect(() => {
    const canvas = canvasRef.current!
    const ctx = canvas.getContext('2d')!
    const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height)

    syncCanvas(imageData.data, (data, x0, y0, x1, y1) => {
      const rowBytes = canvasWidth * 4
      for (let y = y0; y < y1; y++) {
        const start = y * rowBytes + x0 * 4
        const end = y * rowBytes + x1 * 4
        imageData.data.set(data.subarray(start, end), start)
      }
      ctx.putImageData(imageData, 0, 0, x0, y0, x1 - x0, y1 - y0)
    })
  }, [syncCanvas])

  return (
    <canvas
      ref={canvasRef}
      width={canvasWidth}
      height={canvasHeight}
      onClick={handleClick}
      style={{
        width: canvasWidth * scale,
        height: canvasHeight * scale
      }}
    />
  )
}

export default App
