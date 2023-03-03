import { BoardState, ChatMessage, ClientInfo, GameInfo, MoveInfo, MovePattern, PiecePosition } from "@/types";
import { RootState } from "./state";
const getDefaultState = () => {
  return {
    boards: {},  // roomId: {boardState} for that room
    chatMessages: {},  // roomId: [messages for that room]
    gameInfo: {}, 
  }
}

type PayloadWithRoomId<T> = T & {roomId: string};

const mutations ={
    resetState (state:RootState) {
      Object.assign(state, getDefaultState())
    },
    updateBoardState (state:RootState,payload:{roomId:string,boardState:BoardState}) {
        state.boards[payload.roomId] = payload.boardState;
      },
    addMessage (state:RootState,messageInfo:PayloadWithRoomId<ChatMessage>) {
      if (state.chatMessages[messageInfo.roomId]){
      state.chatMessages[messageInfo.roomId].push(messageInfo);
      }
      else {
        state.chatMessages[messageInfo.roomId]=[messageInfo];
      }
    },
    updateGameInfo (state:RootState,payload:PayloadWithRoomId<GameInfo>){
      if(!state.gameInfo){state.gameInfo={p1:payload.p1,turn:payload.turn}}
      state.gameInfo.p1 = payload.p1;
      state.gameInfo.p2 = payload.p2;
      state.gameInfo.turn = payload.turn;
      state.gameInfo.members  = payload.members;
    },

    setClientInfo(state:RootState,payload:ClientInfo){
      state.clientInfo = payload;
    },
    setSelection(state:RootState,payload:PiecePosition){
      state.curStartPos = {piece: payload.piece, row: payload.row, col: payload.col}
    },
    undoSelection(state:RootState){
      state.curStartPos = null
    },
    performMove(state:RootState,moveInfo: MoveInfo){
      state.currentMove = moveInfo
      //after move
      state.curStartPos = null
      if (state.gameInfo){
        state.gameInfo.turn = state.gameInfo.turn == 'w' ? 'b' : 'w';
      }
    },
    websocketError(state:RootState,errorMessage:string){
      state.errorMessage = errorMessage;
    },
    storeMovePatterns(state:RootState,payload:{movePatterns: MovePattern[]}){
      state.movePatterns = payload.movePatterns;
    },
    setResult(state:RootState,result:string){
      if (state.gameInfo){
        state.gameInfo.result = result
      }
    }
}

export default mutations;