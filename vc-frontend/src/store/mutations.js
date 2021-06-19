
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
      console.log('new',payload.boardState)
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
    updateGameInfo (state,gameInfo){
      state.gameInfo[gameInfo.roomId] = gameInfo;
    }
}

export default mutations;