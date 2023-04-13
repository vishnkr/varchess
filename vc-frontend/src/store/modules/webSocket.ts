import { Module,ActionContext } from 'vuex';
import { RootState } from '../state';
import { IMoveInfo, IMoveInfoPayload, IWsMessage } from '../../types';
import store from '..';
import * as MutationTypes from '../mutation_types';
import * as ActionTypes from "../action_types";

const server_host = import.meta.env.VITE_SERVER_WS;

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
              let msgData = apiMsg.data;
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
        state.ws.close();
        state.ws = null;
      }
    },

    async sendJSONReq(context: ActionContext<WebSocketState, RootState>,payload:{type:string,msg:any}){                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        
      await new Promise<void>((resolve, reject) => {
        if (!context.state.ws) {
          reject(new Error('WebSocket not available'));
          store.commit(MutationTypes.SET_SERVER_STATUS,{isOnline:false,errorMessage:'WebSocket connection lost, please head to the homepage'})
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

    validateMove(context: ActionContext<WebSocketState, RootState>, payload: IMoveInfoPayload) {
      if (!context.rootState.gameInfo?.result) {
        const data: IMoveInfo = {
          piece: payload.piece,
          roomId: payload.roomId,
          srcRow: payload.srcRow,
          srcCol: payload.srcCol,
          destRow: payload.destRow,
          destCol: payload.destCol,
          castle: payload.castle
        };
        context.dispatch('sendJSONReq', { type: 'performMove', data });
      }
    },
  },
};

export default webSocketModule;