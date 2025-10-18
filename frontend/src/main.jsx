import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './app/App.jsx'
import './styles/tailwind.css'   // если подключали — оставь
import './styles/tailwind.css';
import './styles/chat-anim.css';


ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode><App /></React.StrictMode>
)
