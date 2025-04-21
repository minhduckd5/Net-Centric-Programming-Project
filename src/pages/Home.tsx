import { Link } from 'react-router-dom';
import { useEffect } from 'react';
import * as anime from 'animejs';
import '../styles/pages/Home.scss';

const Home = () => {
  return (
    <div className="home">
      <h1>Welcome to Game Hub</h1>
      <p>Choose a game to play and chat with other players in real-time!</p>
      
      <div className="game-grid">
        <Link to="/games/hangman" className="game-card">
          <h2>Hangman</h2>
          <p>Try to guess the word one letter at a time!</p>
        </Link>
        
        <Link to="/games/guessing" className="game-card">
          <h2>Guessing Game</h2>
          <p>Test your knowledge with multiple choice questions!</p>
        </Link>
      </div>
    </div>
  );
};

export default Home;
