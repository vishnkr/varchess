import axios from "axios";
import { ActionContext } from "vuex";
import { RootState } from "./state";

interface PossibleSquaresResponse {
  moves: number[];
}

const actions = {
  async getPossibleToSquares(
    { state }: ActionContext<RootState, RootState>,
    payload: { roomId: string; color: string; srcRow: number; srcCol: number; piece: string }
  ): Promise<PossibleSquaresResponse> {
    const { roomId, color, srcRow, srcCol, piece } = payload;
    const url = `${process.env.VUE_APP_SERVER_HOST}/possible-squares?roomid=${roomId}&color=${color}&piece=${piece}&src_row=${srcRow}&src_col=${srcCol}`;
    const possibleMoves = await axios.get(url);
    return possibleMoves.data;
  }
};

export default actions;
