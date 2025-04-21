import { useState } from 'react';
import { Link } from 'react-router-dom';
import '../../styles/layout/Header.scss';

const Header = () => {
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);

  const games = [
    { name: 'Hangman', path: '/games/hangman' },
    { name: 'Guessing Game', path: '/games/guessing' },
  ];

  return (
    <header className="header">
      <nav className="navbar">
        <Link to="/" className="logo">
          Game Hub
        </Link>
        
        <div className="nav-items">
          <Link to="/" className="nav-item">Home</Link>
          
          <div 
            className="dropdown"
            onMouseEnter={() => setIsDropdownOpen(true)}
            onMouseLeave={() => setIsDropdownOpen(false)}
          >
            <button className="dropdown-trigger">Games</button>
            {isDropdownOpen && (
              <div className="dropdown-menu">
                {games.map((game) => (
                  <Link 
                    key={game.path}
                    to={game.path}
                    className="dropdown-item"
                    onClick={() => setIsDropdownOpen(false)}
                  >
                    {game.name}
                  </Link>
                ))}
              </div>
            )}
          </div>
        </div>
      </nav>
    </header>
  );
};

export default Header;
