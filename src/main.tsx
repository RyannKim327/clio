import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import '@/index.css'
import { BrowserRouter, HashRouter, Route, Routes } from 'react-router'
import Reader from '@/routes/reader/index.tsx'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <HashRouter>
      <Routes>
        <Route path='' element={<Reader />} />
      </Routes>
    </HashRouter>
  </StrictMode>,
)
