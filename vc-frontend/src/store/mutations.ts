import { BoardState, ChatMessage, MoveInfo, MovePattern, PiecePosition, Players } from "@/types";
import * as MutationTypes from "../utils/action_mutation_types";
import { RootState } from "./state";


type PayloadWithRoomId<T> = T & {roomId: string};

const mutations ={
    [MutationTypes.UPDATE_BOARD_STATE](state:RootState,payload:{roomId:string,boardState:BoardState}) {
        state.board = payload.boardState;
    },

    [MutationTypes.ADD_CHAT_MESSAGE](state:RootState,messageInfo:PayloadWithRoomId<ChatMessage>) {
      if (state.chatMessages[messageInfo.roomId]){
        state.chatMessages[messageInfo.roomId].push(messageInfo);
      }
      else {
        state.chatMessages[messageInfo.roomId]=[messageInfo];
      }
    },
    [MutationTypes.SET_PLAYERS](state:RootState,payload:Players){
      if (!state.gameInfo){
        state.gameInfo = {players:payload, members: []}
        return
      }
      state.gameInfo.players = payload
    },

    [MutationTypes.UPDATE_MEMBERS](state:RootState,payload:{members:string[]}){
      if (state.gameInfo){
        state.gameInfo.members = payload.members;
      }
    },

    [MutationTypes.SET_SRC_SELECTION](state:RootState,payload:PiecePosition){
      state.curStartPos = {piece: payload.piece, row: payload.row, col: payload.col}
    },
    [MutationTypes.UNDO_SRC_SELECTION](state:RootState){
      state.curStartPos = null
    },
    [MutationTypes.PERFORM_MOVE](state:RootState,moveInfo: MoveInfo){
      state.currentMove = moveInfo
      //after move
      state.curStartPos = null
      if (state.board){
        state.board.turn = state.board.turn == 'w' ? 'b' : 'w';
      }
    },
    [MutationTypes.SET_SERVER_STATUS](state:RootState,payload:{isOnline:boolean,errorMessage:string|null}){
      state.serverStatus = payload;
    },

    [MutationTypes.SET_MOVE_PATTERNS](state:RootState,payload:{movePatterns: MovePattern[]}){
      state.movePatterns = payload.movePatterns;
    },

    [MutationTypes.SET_RESULT](state:RootState,result:string){
      if (state.gameInfo){
        state.gameInfo.result = result
      }
    }
}

export default mutations;