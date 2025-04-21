import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Header from './components/layout/Header';
import Footer from './components/layout/Footer';
import Home from './pages/Home';
import HangmanGame from './pages/games/HangmanGame';
import GuessingGame from './pages/games/GuessingGame';
import Chat from './components/chat/Chat';
import './styles/App.scss';

function App() {
  return (
    <Router>
      <div className="app">
        <Header />
        <main className="main-content">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/games/hangman" element={<HangmanGame />} />
            <Route path="/games/guessing" element={<GuessingGame />} />
          </Routes>
          <Chat />
        </main>
        <Footer />
      </div>
    </Router>
  );
}

export default App;
