import axios from "axios";

const actions ={
    async getPossibleToSquares(state,payload){
        let possibleMoves = await axios.post(`${process.env.VUE_APP_SERVER_HOST}/getPossibleToSquares`,JSON.stringify(payload));
        return possibleMoves.data;
    }
}
export default actions;