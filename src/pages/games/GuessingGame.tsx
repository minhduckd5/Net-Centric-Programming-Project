// === CHANGE START: migrate to Anime.js v4 ESM modules ===
// Old: import anime from 'animejs';
import { animate, utils } from 'animejs'; // v4 named imports
// === CHANGE END ===

import { useState } from 'react';
import '../../styles/games/GuessingGame.scss';

interface Question {
  question: string;
  options: string[];
  correctAnswer: number;
}

const QUESTIONS: Question[] = [
  {
    question: "Which programming language is used for WebSocket server in this project?",
    options: ["Python", "Java", "Golang", "Node.js"],
    correctAnswer: 2
  },
  {
    question: "What animation library is used in this project?",
    options: ["GSAP", "AnimeJS", "Framer Motion", "React Spring"],
    correctAnswer: 1
  },
  {
    question: "What type of connection is used for real-time chat?",
    options: ["HTTP", "WebSocket", "Server-Sent Events", "Long Polling"],
    correctAnswer: 1
  }
];

const GuessingGame = () => {
  const [currentQuestion, setCurrentQuestion] = useState(0);
  const [score, setScore] = useState(0);
  const [gameOver, setGameOver] = useState(false);
  const [selectedAnswer, setSelectedAnswer] = useState<number | null>(null);

  const handleAnswer = (optionIndex: number) => {
    if (selectedAnswer !== null) return;
    setSelectedAnswer(optionIndex);
    const correct = optionIndex === QUESTIONS[currentQuestion].correctAnswer;

    // === CHANGE START: use animate() v4 syntax for option feedback ===
    animate(`.option-${optionIndex}`, {
      scale: [
        { to: 1.1, duration: 300, easing: 'inOutQuad' },
        { to: 1, duration: 300, easing: 'inOutQuad' }
      ],
      // color feedback using modifier utility (example)
      color: {
        to: correct ? '#10B981' : '#EF4444'
      }
    });
    // === CHANGE END ===

    setTimeout(() => {
      if (correct) {
        setScore(prev => prev + 1);
      }
      if (currentQuestion < QUESTIONS.length - 1) {
        setCurrentQuestion(prev => prev + 1);
        setSelectedAnswer(null);
      } else {
        setGameOver(true);
      }
    }, 1000);
  };

  const restartGame = () => {
    setCurrentQuestion(0);
    setScore(0);
    setGameOver(false);
    setSelectedAnswer(null);
  };

  return (
    <div className="guessing-game">
      <h1>Guessing Game</h1>
      {!gameOver ? (
        <div className="question-container">
          <div className="progress">
            Question {currentQuestion + 1} of {QUESTIONS.length}
          </div>
          <h2>{QUESTIONS[currentQuestion].question}</h2>
          <div className="options">
            {QUESTIONS[currentQuestion].options.map((option, index) => (
              <button
                key={index}
                className={`option option-${index} ${selectedAnswer === index
                    ? index === QUESTIONS[currentQuestion].correctAnswer
                      ? 'correct'
                      : 'wrong'
                    : ''
                  }`}
                onClick={() => handleAnswer(index)}
                disabled={selectedAnswer !== null}
              >
                {option}
              </button>
            ))}
          </div>
          <div className="score">
            Current Score: {score}
          </div>
        </div>
      ) : (
        <div className="game-over">
          <h2>Game Complete!</h2>
          <p>Your Final Score: {score} out of {QUESTIONS.length}</p>
          <button onClick={restartGame}>Play Again</button>
        </div>
      )}
    </div>
  );
};

export default GuessingGame;
