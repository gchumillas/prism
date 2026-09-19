import { render } from 'preact'
import App from './app.tsx'

window.wasmReady
  .then(() => {
    render(<App />, document.getElementById('app')!)
  })
  .catch((err) => {
    console.error('Failed to load WASM:', err)
  })
