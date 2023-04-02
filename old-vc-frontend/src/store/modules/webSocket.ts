import { Module,ActionContext } from 'vuex';
import { RootState } from '../state';
import { MoveInfo, MoveInfoPayload, WsMessage } from '@/types';
import store from '..';
import * as MutationTypes from '@/utils/mutation_types';

const server_host = process.env.VUE_APP_SERVER_WS;

function isOpen(ws: WebSocket) {
  return ws.readyState === WebSocket.OPEN;
}

export interface WebSocketState {
  ws: WebSocket | null;
  userId: string | null;
}
interface ConnectParams {
  roomId: string;
  username: string;
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
    connect({ commit, rootState }: ActionContext<WebSocketState, RootState>,{roomId,username}:ConnectParams) {
      const ws = new WebSocket(`${server_host}/ws?roomId=${roomId}&username=${username}`);
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
             store.commit(MutationTypes.ADD_CHAT_MESSAGE,msgData)
              break;
            }
            case "error":{
              store.commit(MutationTypes.SET_SERVER_STATUS,apiMsg.data)
              break;
            }

            case "memberList":{
              let msgData = apiMsg.data;
              store.commit(MutationTypes.SET_PLAYERS,{
                p1: msgData.p1,
                p2: msgData.p2 ?? null
              });
              store.commit(MutationTypes.UPDATE_MEMBERS,{members:msgData.members})
            }
            case "performMove":{
              if(apiMsg.isValid){ //only if move is valid you perform commit
                store.commit(MutationTypes.PERFORM_MOVE,apiMsg)
              }
              if(apiMsg.result){
                store.commit(MutationTypes.SET_RESULT,apiMsg.result)
              }
                break;
              }
            case "result":{
              if(apiMsg.result){
                store.commit(MutationTypes.SET_RESULT,apiMsg.result)
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
  
    sendMessage(context: ActionContext<WebSocketState, RootState>, payload: { message: string; username: string; roomId: string }) {
      context.dispatch('sendJSONReq', { type: 'chatMessage', data: payload });
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