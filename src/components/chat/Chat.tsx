import { useState, useEffect, useRef } from 'react';
import { wsService } from '../../services/websocket';
import { ChatMessage, WebSocketMessage } from '../../types/game';
import '../../styles/components/Chat.scss';

const Chat = () => {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [newMessage, setNewMessage] = useState('');
  const [username, setUsername] = useState('');
  const [isUsernameSet, setIsUsernameSet] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    wsService.connect();
    
    const unsubscribe = wsService.onMessage((message: WebSocketMessage) => {
      if (message.type === 'chat') {
        setMessages(prev => [...prev, message.payload as ChatMessage]);
      }
    });

    return () => unsubscribe();
  }, []);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleSubmitUsername = (e: React.FormEvent) => {
    e.preventDefault();
    if (username.trim()) {
      setIsUsernameSet(true);
      wsService.sendMessage({
        type: 'user_presence',
        payload: {
          type: 'join',
          username: username.trim()
    }
      });
    }
  };

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    if (newMessage.trim() && username) {
      const chatMessage: ChatMessage = {
        id: Date.now().toString(),
        username,
        message: newMessage.trim(),
        timestamp: Date.now(),
};

      wsService.sendMessage({
        type: 'chat',
        payload: chatMessage,
      });

      setNewMessage('');
    }
  };

  if (!isUsernameSet) {
    return (
      <div className="chat-container">
        <div className="username-form">
          <h3>Enter your username to chat</h3>
          <form onSubmit={handleSubmitUsername}>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="Enter username"
              maxLength={20}
              required
            />
            <button type="submit">Join Chat</button>
          </form>
        </div>
      </div>
    );
  }

  return (
    <div className="chat-container">
      <div className="chat-messages">
        {messages.map((msg) => (
          <div 
            key={msg.id} 
            className={`message ${msg.username === username ? 'own-message' : ''}`}
          >
            <span className="username">{msg.username}</span>
            <p>{msg.message}</p>
            <span className="timestamp">
              {new Date(msg.timestamp).toLocaleTimeString()}
            </span>
          </div>
        ))}
        <div ref={messagesEndRef} />
      </div>

      <form className="chat-input" onSubmit={handleSendMessage}>
        <input
          type="text"
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
          placeholder="Type a message..."
          maxLength={200}
        />
        <button type="submit">Send</button>
      </form>
    </div>
  );
};

export default Chat;
