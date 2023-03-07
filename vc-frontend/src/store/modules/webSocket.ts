import { Module } from 'vuex';
import { RootState } from '../state';
import { ActionContext } from 'vuex';
import { MoveInfo, MoveInfoPayload } from '@/types';
import store from '..';

const server_host = process.env.VUE_APP_SERVER_WS;

function isOpen(ws: WebSocket) {
  return ws.readyState === WebSocket.OPEN;
}

export interface WebSocketState {
  ws: WebSocket | null;
  userId: string | null;
}

const webSocketModule: Module<WebSocketState, RootState> = {
  namespaced: true,
  state: {
    ws: null,
    userId: null,
  },
  mutations: {
    setWS(state, ws: WebSocket) {
      state.ws = ws;
    },
    setUserId(state,userId:string){
      state.userId = userId;
    }
  },
  actions: {
    connect({ commit, rootState }) {
      const ws = new WebSocket(`${server_host}/ws`);
      ws.onopen = function(){
        console.log("Socket opened");
      }

      ws.onerror = function(){
        commit('websocketError','Connection to server could not be established! Try again soon!')
      }

      ws.onmessage = function(msg){
        let apiMsg = JSON.parse(msg.data);
        switch(apiMsg.type){
            case "chatMessage": {
              let msgData = JSON.parse(apiMsg.data);
              if (rootState.chatMessages[msgData.roomId]==undefined){
                   msgData.id=1
               } else {
                  msgData.id = (rootState.chatMessages[msgData.roomId]).length+1;
              }
             store.commit('addMessage',msgData)
              break;
            }
            case "error":{
              store.commit('setServerStatus',apiMsg.data)
              break;
            }
            case "gameInfo":{
              store.commit('updateGameInfo',apiMsg)
                if(apiMsg.result){
                  commit("setResult",apiMsg.result)
                }
                break;
            }

            case "performMove":{
              if(apiMsg.isValid){ //only if move is valid you perform commit
                store.commit('performMove',apiMsg)
              }
              if(apiMsg.result){
                store.commit('setResult',apiMsg.result)
              }
                break;
              }
            case "result":{
              if(apiMsg.result){
                store.commit('setResult',apiMsg.result)
              }
              break;
            }
            default:
              break;
            }
      }
    commit('setWS',ws);
    },

    close({ state }) {
      if (state.ws !== null) {
        console.log('closing ws');
        state.ws.close();
        state.ws = null;
      }
    },

    async sendJSONReq(context: ActionContext<WebSocketState, RootState>,payload:{type:string,msg:any}){
      console.log('in send json',context.state.ws,isOpen(context.state.ws!))
      await new Promise<void>((resolve, reject) => {
        if (!context.state.ws) {
          reject(new Error('WebSocket not available'));
          return;
        }

        if (isOpen(context.state.ws)) {
          resolve();
        } else {
          context.state.ws.addEventListener('open', () => {
          resolve();
          });
        }
      });
      if (!context.state.ws || !isOpen(context.state.ws)) return;
      console.log('sending jsonreq',JSON.stringify(payload))
      context.state.ws.send(JSON.stringify(payload));
      
    },

    createRoom(context: ActionContext<WebSocketState, RootState>, payload: { roomId: string; username: string; fen: string; customMovePatterns?: any }) {
      context.dispatch('sendJSONReq', { type: 'createRoom', data: payload });
    },
  
    joinRoom(context: ActionContext<WebSocketState, RootState>, payload: { roomId: string; username: string }) {
      context.dispatch('sendJSONReq', { type: 'joinRoom', data: payload });
    },
  
    sendMessage(context: ActionContext<WebSocketState, RootState>, payload: { message: string; username: string; roomId: string }) {
      context.dispatch('sendJSONReq', { type: 'chatMessage', data: payload });
    },
  
    requestGameinfo(context: ActionContext<WebSocketState, RootState>, payload: { roomId: string }) {
      context.dispatch('sendJSONReq', { type: 'reqGameInfo', data:payload });
    },
  
    sendResign(context: ActionContext<WebSocketState, RootState>, payload: { roomId: string; color: string }) {
      context.dispatch('sendJSONReq', { type: 'resign', data: payload });
    },
  
    sendDrawOffer(context: ActionContext<WebSocketState, RootState>, payload: { roomId: string; color: string }) {
      context.dispatch('sendJSONReq', { type: 'draw', data: payload});
    },

    sendMoveInfo(context: ActionContext<WebSocketState, RootState>, payload: MoveInfoPayload) {
      if (!context.rootState.gameInfo?.result) {
        const data: MoveInfo = {
          piece: payload.piece,
          roomId: payload.roomId,
          srcRow: payload.srcRow - 1,
          srcCol: payload.srcCol - 1,
          destRow: payload.destRow - 1,
          destCol: payload.destCol - 1,
          color: payload.color,
          castle: !!payload.castle
        };
        context.dispatch('sendJSONReq', { type: 'performMove', data });
      }
    },
  },
};

export default webSocketModule;