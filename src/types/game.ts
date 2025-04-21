export interface GameState {
  isPlaying: boolean;
  score: number;
  gameOver: boolean;
}

export interface ChatMessage {
  id: string;
  username: string;
  message: string;
  timestamp: number;
}

export interface WebSocketMessage {
  type: 'chat' | 'game_state' | 'user_presence';
  payload: any;
}
