import { BrowserRouter as Router, Route, Routes } from 'react-router-dom'
import { Toaster } from 'react-hot-toast'
import './App.css'
import Home from './pages/home'
import SignIn from './pages/signin'

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/auth" element={<SignIn />} />
      </Routes>
      <Toaster />
    </Router>
  )
}

export default App
