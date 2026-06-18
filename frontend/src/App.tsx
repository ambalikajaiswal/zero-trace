import { Routes, Route } from 'react-router-dom'
import Header from './components/Header'
import CreateSecret from './pages/CreateSecret'
import ViewSecret from './pages/ViewSecret'
import ShareLink from './pages/ShareLink'
import NotFound from './pages/NotFound'

function App() {
  return (
    <div className="app">
      <Header />
      <main className="container">
        <Routes>
          <Route path="/" element={<CreateSecret />} />
          <Route path="/shared" element={<ShareLink />} />
          <Route path="/secret/:id" element={<ViewSecret />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </main>
    </div>
  )
}

export default App
