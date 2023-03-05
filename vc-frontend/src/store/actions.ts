import axios from "axios";
import { ActionContext } from "vuex";
import { RootState } from "./state";

interface PossibleSquaresResponse {
  moves: number[];
}

const BASE_URL = process.env.VUE_APP_SERVER_HOST;

const actions = {
  async getPossibleToSquares(
    { state }: ActionContext<RootState, RootState>,
    payload: { roomId: string; color: string; srcRow: number; srcCol: number; piece: string }
  ): Promise<PossibleSquaresResponse> {
    const { roomId, color, srcRow, srcCol, piece } = payload;
    const url = `${BASE_URL}/possible-squares?roomid=${roomId}&color=${color}&piece=${piece}&src_row=${srcRow}&src_col=${srcCol}`;
    const possibleMoves = await axios.get(url);
    return possibleMoves.data;
  },

  async checkServerStatus({commit}:ActionContext<RootState,RootState>):Promise<void>{
    try{
      const response = await axios.get(`${BASE_URL}/server-status`);
      if (response.status==200){
        commit('setServerStatus',{isOnline:true,errorMessage:null})
      } else {
        commit('setServerStatus',{isOnline:false,errorMessage:"Can't connect to server"})
      }
    } catch(error) {
      commit('setServerStatus',{isOnline:false,errorMessage:"Connection to the server cannot be established at the moment, Please try again later."})
    }
  }
};

export default actions;
