import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import '@/index.css'
import { BrowserRouter, HashRouter, Route, Routes } from 'react-router'
import Reader from '@/routes/reader/index.tsx'
import Read from './routes/reader/read'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <div className='w-dvw h-dvh bg-bg text-fg'>
      <HashRouter>
        <Routes>
          <Route path='' element={<Reader />} />
          <Route path='read' element={<Read />} />
        </Routes>
      </HashRouter>
    </div>
  </StrictMode>,
)
