import { ChatMessage, WebSocketMessage } from '../types/game';

class WebSocketService {
  private ws: WebSocket | null = null;
  private messageCallbacks: ((message: WebSocketMessage) => void)[] = [];

  connect() {
    // Use relative path when using proxy, or full URL in production
    const wsUrl = process.env.NODE_ENV === 'production' 
      ? 'wss://your-production-url/ws'
      : `ws://${window.location.hostname}:8080/ws`;
    this.ws = new WebSocket(wsUrl);

    this.ws.onopen = () => {
      console.log('Connected to WebSocket');
    };

    this.ws.onmessage = (event) => {
      try {
      const message: WebSocketMessage = JSON.parse(event.data);
      this.messageCallbacks.forEach(callback => callback(message));
      } catch (error) {
        console.error('Error parsing WebSocket message:', error);
      }
    };

    this.ws.onclose = () => {
      console.log('Disconnected from WebSocket');
      setTimeout(() => this.connect(), 5000); // Reconnect after 5 seconds
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
  }

  sendMessage(message: WebSocketMessage) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message));
    } else {
      console.warn('WebSocket is not connected');
  }
}

  onMessage(callback: (message: WebSocketMessage) => void) {
    this.messageCallbacks.push(callback);
    return () => {
      this.messageCallbacks = this.messageCallbacks.filter(cb => cb !== callback);
    };
  }
}

export const wsService = new WebSocketService();
