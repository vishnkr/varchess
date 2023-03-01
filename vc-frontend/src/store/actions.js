import axios from "axios";

//{roomId:this.roomId,color:pieceInfo.pieceColor,srcRow:pieceInfo.row-1,srcCol:pieceInfo.col-1,piece:pieceInfo.pieceType}
const actions ={
    async getPossibleToSquares(state,payload){
        const { roomId, color, srcRow, srcCol, piece } = payload;
        const url = `${process.env.VUE_APP_SERVER_HOST}/possible-squares?roomid=${roomId}&color=${color}&piece=${piece}&src_row=${srcRow}&src_col=${srcCol}`
        let possibleMoves = await axios.get(url);
        return possibleMoves.data;
    }
}
export default actions;