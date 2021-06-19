const state = {
    boards: {},  // roomId: {boardState} for that room
    chatMessages: {},  // roomId: [messages for that room]
    gameInfo: {},  // roomId: {game info}
    curStartPosClick: [],
    curDestPosClick: [],
}

export default state