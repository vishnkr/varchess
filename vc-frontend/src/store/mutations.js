const mutations ={
    updateBoardState (state,newBoardState) {
        state.boardState = newBoardState;
      },
    setUsername (state,username) {
      state.username = username;
    }
}

export default mutations;