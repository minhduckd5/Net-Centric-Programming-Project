// === CHANGE START: migrate to Anime.js v4 ESM modules ===
// Old: import anime from 'animejs';
import { animate, utils } from 'animejs'; // v4 named imports; use createSpring if spring easing is needed
// === CHANGE END ===

import { useState, useEffect } from 'react';
import '../../styles/games/HangmanGame.scss';

const WORDS = ['TYPESCRIPT', 'JAVASCRIPT', 'REACT', 'GOLANG', 'WEBSOCKET'];

const HangmanGame = () => {
  const [word, setWord] = useState('');
  const [guessedLetters, setGuessedLetters] = useState<Set<string>>(new Set());
  const [wrongGuesses, setWrongGuesses] = useState(0);
  const [gameOver, setGameOver] = useState(false);
  const [won, setWon] = useState(false);

  useEffect(() => {
    startNewGame();
  }, []);

  const startNewGame = () => {
    const randomWord = WORDS[Math.floor(Math.random() * WORDS.length)];
    setWord(randomWord);
    setGuessedLetters(new Set());
    setWrongGuesses(0);
    setGameOver(false);
    setWon(false);
  };

  const handleGuess = (letter: string) => {
    if (gameOver || guessedLetters.has(letter)) return;

    const newGuessedLetters = new Set(guessedLetters).add(letter);
    setGuessedLetters(newGuessedLetters);

    if (!word.includes(letter)) {
      const newWrong = wrongGuesses + 1;
      setWrongGuesses(newWrong);

      // === CHANGE START: use animate() and v4 options syntax ===
      animate('.hangman-figure', {
        translateX: [
          { to: -10, duration: 100, easing: 'inOutQuad' },
          { to: 10, duration: 100, easing: 'inOutQuad' },
          { to: 0, duration: 100, easing: 'inOutQuad' }
        ],
      });
      // === CHANGE END ===

      if (newWrong >= 6) {
        setGameOver(true);
      }
    } else {
      // === CHANGE START: letter reveal animation with v4 syntax ===
      animate(`.letter-${letter}`, {
        scale: [
          { to: 1.2, duration: 300 },
          { to: 1, duration: 200 }
        ],
      });
      // === CHANGE END ===
    }

    const isWon = word.split('').every(ch => newGuessedLetters.has(ch));
    if (isWon) {
      setWon(true);
      setGameOver(true);

      // === CHANGE START: win animation with v4 syntax ===
      animate('.word-container', {
        scale: [
          { to: 1.1, duration: 400 },
          { to: 1, duration: 400 }
        ],
      });
      // === CHANGE END ===
    }
  };

  const renderHangmanFigure = () => {
    const parts = [
      <line key="base" x1="20" y1="180" x2="180" y2="180" stroke="#2c3e50" strokeWidth="4" />,
      <line key="pole" x1="60" y1="20" x2="60" y2="180" stroke="#2c3e50" strokeWidth="4" />,
      <line key="top" x1="60" y1="20" x2="140" y2="20" stroke="#2c3e50" strokeWidth="4" />,
      <line key="rope" x1="140" y1="20" x2="140" y2="40" stroke="#2c3e50" strokeWidth="4" />,
      <circle key="head" cx="140" cy="60" r="20" stroke="#2c3e50" strokeWidth="4" fill="none" />,
      <line key="body" x1="140" y1="80" x2="140" y2="120" stroke="#2c3e50" strokeWidth="4" />,
      <g key="arms">
        <line x1="140" y1="90" x2="120" y2="100" stroke="#2c3e50" strokeWidth="4" />
        <line x1="140" y1="90" x2="160" y2="100" stroke="#2c3e50" strokeWidth="4" />
      </g>,
      <g key="legs">
        <line x1="140" y1="120" x2="120" y2="140" stroke="#2c3e50" strokeWidth="4" />
        <line x1="140" y1="120" x2="160" y2="140" stroke="#2c3e50" strokeWidth="4" />
      </g>
    ];

    return parts.slice(0, wrongGuesses);
  };

  return (
    <div className="hangman-game">
      <h1>Hangman</h1>

      <div className="game-container">
        <div className="hangman-figure">
          <svg width="200" height="200" viewBox="0 0 200 200">
            {renderHangmanFigure()}
          </svg>
        </div>

        <div className="word-container">
          {Array.from(word).map((letter, idx) => (
            <span
              key={idx}
              className={`letter letter-${letter}`}
            >
              {guessedLetters.has(letter) ? letter : '_'}
            </span>
          ))}
        </div>

        <div className="keyboard">
          {Array.from("ABCDEFGHIJKLMNOPQRSTUVWXYZ").map(letter => (
            <button
              key={letter}
              onClick={() => handleGuess(letter)}
              disabled={guessedLetters.has(letter) || gameOver}
              className={`keyboard-button ${guessedLetters.has(letter) ? 'guessed' : ''} ${word.includes(letter) && guessedLetters.has(letter) ? 'correct' : ''}`}
            >
              {letter}
            </button>
          ))}
        </div>

        {gameOver && (
          <div className="game-over">
            <h2>{won ? 'Congratulations!' : 'Game Over!'}</h2>
            <p>The word was: {word}</p>
            <button onClick={startNewGame}>Play Again</button>
          </div>
        )}
      </div>
    </div>
  );
};

export default HangmanGame;
