import { createRoot } from 'react-dom/client'
import { registerSW } from "virtual:pwa-register"

import { StrictMode } from 'react'
import { Provider } from 'react-redux'
import { App } from './App'
import store from './store'


createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <Provider store={store}>
      <App />
    </Provider>
  </StrictMode>
)

if ("serviceWorker" in navigator) {
  registerSW()
}