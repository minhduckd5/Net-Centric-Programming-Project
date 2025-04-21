import '../../styles/layout/Footer.scss';

const Footer = () => {
  return (
    <footer className="footer">
      <div className="footer-content">
        <div className="footer-section">
          <h3>Game Hub</h3>
          <p>A fun collection of web games with real-time chat</p>
        </div>
        
        <div className="footer-section">
          <h3>Credits</h3>
          <p>Created with React, TypeScript, and Go</p>
          <p>Powered by WebSocket and AnimeJS</p>
        </div>
        
        <div className="footer-section">
          <h3>Connect</h3>
          <a href="https://github.com/yourusername" target="_blank" rel="noopener noreferrer">
            GitHub
          </a>
        </div>
      </div>
      
      <div className="footer-bottom">
        <p>&copy; {new Date().getFullYear()} Game Hub. All rights reserved.</p>
      </div>
    </footer>
  );
};

export default Footer;
