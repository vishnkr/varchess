import axios from "axios";

const actions ={
    async getPossibleToSquares(state,payload){
        let possibleMoves = await axios.post("http://localhost:5000/getPossibleToSquares",JSON.stringify(payload));
        return possibleMoves.data;
    }
}
export default actions;