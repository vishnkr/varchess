import { writable, type Writable } from 'svelte/store';
import {
  type RuleEditorState,
  VariantType,
  type MovePattern,
  EventUserConnect,
  EventJoinGame,
  EventStartGame,
  EventChatMessage,
} from './types';
import { camelToSnake } from '$lib/utils';
import type { IPiece, Move} from '$lib/board/types';

export enum Role {
  Viewer,
  Black,
  White
}

export enum MessageType {
  UserJoin,
  UserLeave,
  ChatMessage,
  ResultUpdate
}

export interface ChatMessage {
  messageType: MessageType;
  content?: string;
  username?: string;
}

export interface Member {
  id: number;
  username: string;
  role?: Role;
  isHost?: boolean;
}

export interface User {
  username?: string;
}

export interface Config {
  variantType: VariantType;
  dimensions: {
    ranks: number;
    files: number;
  };
  fen: string;
  pieceProps: Record<string, MovePattern>;
  additionalData?: object;
}

const roomId = writable<string | null>(null);
const members = writable<Member[]>([]);

function newConfigStore() {
  const { subscribe, set, update } = writable<Config | null>(null);
  const setConfig = (config: Config) => {
    set(config);
  };

  const removeConfig = () => {
    set(null);
  };
  return {
    setConfig,
    removeConfig,
    update,
    subscribe
  };
}

function newChatStore() {
  const { subscribe, set, update } = writable<ChatMessage[]>([]);
  const userJoin = (username: string) => {
    update((chats) => [
      ...chats,
      {
        messageType: MessageType.UserJoin,
        username: username,
        content: `${username} has joined the Room`
      }
    ]);
  };
  const userLeave = (username: string) => {
    update((chats) => [
      ...chats,
      {
        messageType: MessageType.UserLeave,
        username: username,
        content: `${username} has left the Room`
      }
    ]);
  };
  const updateChat = (username: string, content: string) => {
    update((chats) => [
      ...chats,
      { messageType: MessageType.ChatMessage, username: username, content: content }
    ]);
  };

  return {
    subscribe,
    set,
    update,
    userJoin,
    userLeave,
    updateChat
  };
}

const configStore = newConfigStore();
const chats = newChatStore();

export interface ConnectParams {
  sessionId: string;
  gameId?: string;
  username: string;
}

export interface CreateParams extends ConnectParams {
  gameConfig?: Config;
  color: string;
}
export type ConnectType = 'create' | 'join';

function createWebSocketStore(ws: WebSocket | null) {
  const { subscribe, set, update } = writable<WebSocket | null>(ws);

  const newWebSocketConnection = async (
    url: string,
    params: ConnectParams,
    connectType: ConnectType = 'join'
  ) => {
    const ws = await new WebSocket(url);

    ws.onopen = () => {
      const wsMessage = { event: `game.${connectType}_game`, params: params };
      const json = JSON.stringify(camelToSnake(wsMessage));
      ws.send(json);
      console.log('WebSocket connection success');
    };
    ws.onclose = () => {
      console.log('close called')
      set(null);
      gameId.set(null);
      configStore.removeConfig();
    };
    ws.onerror = (e) => {
      console.error('WebSocket connection error:', e);
      return;
    };
    ws.onmessage = (e) => {
      const data = JSON.parse(e.data);
      console.log('got message', data);
      if (data.success && data.result) {
        switch (data?.event) {
          case EventUserConnect:
            gameId.set(data.result.game_id);
            break;
          case EventJoinGame:
            //gameState.setGameConfig(data.result.game_config)

            /*gameId.set(data.result.game_id);
            gameState.setGameConfig(data.result.game_config)*/
            break;
          case EventStartGame:
            gameState.update((oldState) => {
              return {
                ...oldState,
                players: {
                  playerBlack: data.result.players.black,
                  playerWhite: data.result.players.white
                },
                status: 'InProgress'
              };
            });
            configStore.setConfig(data.result.game_config);
            console.log('updating state', gameState);
            break;
          case EventChatMessage:
            chats.updateChat(data.result.username, data.result.message)
            break;
          default:
            console.log('invalid msg type');
        }
      }
    };
    set(ws);
  };
  const close = ()=>{ if(ws){ws.close()}}
  return {
    newWebSocketConnection,
    subscribe,
    set,
    update,
    close,
  };
}



export interface MoveSelector{
	src: Writable<number|null>;
	dest: Writable<number|null>;
	piece: Writable<IPiece|null>;
  recentMove: Writable<Move | null>;
	legalMoves: Writable<Move[]>;
}

const moveSelector: MoveSelector  = {
  src: writable<number|null>(null),
  dest: writable<number|null>(null),
  piece: writable<IPiece|null>(null),
  legalMoves: writable([]),
  recentMove : writable(null)
}

const wsStore = createWebSocketStore(null);
const gameId = writable<string | null>(null);

export interface Players {
  playerWhite: string;
  playerBlack: string;
}


type Status = 'Waiting' | 'InProgress';
export interface GameState {
  status: Status;
  //chesscore : ChessCoreLib,
  turn: 'w' | 'b';
  players?: Players;
}

function newGameState() {
  const { subscribe, update } = writable<GameState>({
    status: 'Waiting',
    turn: 'w'
  });

  const updateStatus = (newStatus: Status) => {
    update((gameState) => {
      return {
        ...gameState,
        status: newStatus
      };
    });
  };
  const changeTurn = () => {
    update((gameState) => {
      return {
        ...gameState,
        turn: gameState.turn === 'w' ? 'b' : 'w'
      };
    });
  };

  const setPlayers = (playerBlack: string, playerWhite: string) => {
    update((gameState) => {
      return {
        ...gameState,
        players: { playerBlack, playerWhite }
      };
    });
  };

  return {
    updateStatus,
    setPlayers,
    changeTurn,
    update,
    subscribe
  };
}

const gameState = newGameState();

const ruleEditor = writable<RuleEditorState>({
  variantType: VariantType.Checkmate
});

export {
  createWebSocketStore,
  members,
  roomId,
  chats,
  wsStore,
  ruleEditor,
  configStore,
  gameId,
  gameState,
  moveSelector
};
