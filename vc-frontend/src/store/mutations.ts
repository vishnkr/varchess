import { IChatMessage, IMovePatterns, IMoveInfo, IMovePattern, IPiecePosition } from "@/types";
import { BoardState, User } from "../classes";
import {Players, UserInfo} from "../types"
import * as MutationTypes from "./mutation_types";
import { RootState } from "./state";


type PayloadWithRoomId<T> = T & {roomId: string};

const mutations ={
    [MutationTypes.UPDATE_BOARD_STATE](state:RootState,payload:{roomId:string,boardState:BoardState}) {
        state.board = payload.boardState;
    },
    [MutationTypes.SET_USER_INFO](state:RootState,payload:UserInfo){
      state.userInfo.username = payload.username,
      state.userInfo.isAuthenticated = payload.isAuthenticated
      state.userInfo.curGameRole = payload.curGameRole
    },

    [MutationTypes.ADD_CHAT_MESSAGE](state:RootState,messageInfo:PayloadWithRoomId<IChatMessage>) {
      if (state.chatMessages[messageInfo.roomId]){
        state.chatMessages[messageInfo.roomId].push(messageInfo);
      }
      else {
        state.chatMessages[messageInfo.roomId]=[messageInfo];
      }
    },
    [MutationTypes.SET_PLAYERS](state:RootState,payload:Players){
      if (!state.gameInfo){
        state.gameInfo = {players:payload, members: {}}
        return
      }
      state.gameInfo.players = payload
    },

    [MutationTypes.UPDATE_MEMBERS](state:RootState,payload:{members:string[]}){
      if (state.gameInfo){
        for(let member of payload.members){
          state.gameInfo.members[member] = new User()
        }
        
      }
    },

    [MutationTypes.SET_SRC_SELECTION](state:RootState,payload:IPiecePosition|null){
      state.curStartPos = payload
    },
    [MutationTypes.UNDO_SRC_SELECTION](state:RootState){
      state.curStartPos = null
    },
    [MutationTypes.PERFORM_MOVE](state:RootState,moveInfo: IMoveInfo){
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

    [MutationTypes.SET_MOVE_PATTERNS](state:RootState,movePatterns: IMovePatterns){
      state.movePatterns = movePatterns;
    },

    [MutationTypes.SET_RESULT](state:RootState,result:string){
      if (state.gameInfo){
        state.gameInfo.result = result
      }
    }
}

export default mutations;