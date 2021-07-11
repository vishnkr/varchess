
const getDefaultState = () => {
  return {
    boards: {},  // roomId: {boardState} for that room
    chatMessages: {},  // roomId: [messages for that room]
    gameInfo: {}, 
  }
}

const mutations ={
    resetState (state) {
      Object.assign(state, getDefaultState())
    },
    updateBoardState (state,payload) {
        state.boards[payload.roomId] = payload.boardState;
      },
    addMessage (state,messageInfo) {
      if (state.chatMessages[messageInfo.roomId]){
      state.chatMessages[messageInfo.roomId].push(messageInfo);
      }
      else {
        state.chatMessages[messageInfo.roomId]=[messageInfo];
      }
    },
    updateGameInfo (state,payload){
      if(!state.gameInfo[payload.roomId]){state.gameInfo[payload.roomId]={}}
      state.gameInfo[payload.roomId].p1 = payload.p1;
      state.gameInfo[payload.roomId].p2 = payload.p2;
      state.gameInfo[payload.roomId].turn = payload.turn;
      state.gameInfo[payload.roomId].members  = payload.members;
    },
    updateTurn(state){
      state.gameInfo.turn = state.gameInfo.turn == 'w' ? 'b' : 'w';
    },
    setClientInfo(state,payload){
      state.clientInfo.username = payload.username
      state.clientInfo.isPlayer = payload.isPlayer
      if(payload.isPlayer){
        state.clientInfo.color = payload.color
      }
      state.clientInfo.ws = payload.ws
    },
    setSelection(state,payload){
      state.curStartPos = {piece: payload.piece, row: payload.row, col: payload.col}
    },
    undoSelection(state){
      state.curStartPos = null
    },
    performMove(state,moveInfo){
      state.currentMove = moveInfo
      //after move
      state.curStartPos = null
      state.gameInfo.turn = state.gameInfo.turn == 'w' ? 'b' : 'w';
    }
}

export default mutations;