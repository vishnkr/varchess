const state = {
    boards: {},  // roomId: {boardState} for that room
    chatMessages: {},  // roomId: [messages for that room]
    gameInfo: {},  // roomId: {game info}
    curStartPos: null, //store the piece and square info of clicked src square
    curDestPos: null,
    currentMove: null,
    clientInfo: {},
    roomClients:{},
    turn: 'w',
    errorMessage:null,
    movePatterns:null,
}

export default state