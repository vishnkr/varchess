import axios from "axios";

const actions ={
    async getPossibleToSquares(state,payload){
        console.log('pay',payload)
        let possibleMoves = await axios.post("http://localhost:5000/getPossibleToSquares",JSON.stringify(payload));
        console.log('possible',possibleMoves);
        return possibleMoves.data;
    }
}
export default actions;