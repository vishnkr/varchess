import { CreateRoomResponse, MovePattern, PossibleSquaresResponse, RoomState } from "@/types";
import axios from "axios";
import { ActionContext } from "vuex";
import { RootState } from "./state";
import { SET_SERVER_STATUS } from "../utils/action_mutation_types";


export async function makeHttpRequest<T>(
  url: string,
  method: string = 'get',
  data: any = null,
  config: any = {}
):Promise<T>{
  try {
    const response = await axios({
      method,
      url,
      data,
      ...config
    });
    return response.data;
  } catch (error) {
    console.error(error);
    throw error;
  }
}

const BASE_URL = import.meta.env.VITE_SERVER_HOST;

const actions = {

  async getPossibleToSquares(
    { state }: ActionContext<RootState, RootState>,
    payload: { roomId: string; color: string; srcRow: number; srcCol: number; piece: string }
  ): Promise<PossibleSquaresResponse> {
    const { roomId, color, srcRow, srcCol, piece } = payload;
    const url = `${BASE_URL}/possible-squares?roomid=${roomId}&color=${color}&piece=${piece}&src_row=${srcRow}&src_col=${srcCol}`;
    try{
      const possibleMoves = await makeHttpRequest<PossibleSquaresResponse>(url);
      return possibleMoves;
    } catch (error){
      console.log(error)
      throw new Error('Error getting possible moves')
    }
  },

  async checkServerStatus({commit}:ActionContext<RootState,RootState>):Promise<void>{
    try{
      const response = await axios.get(`${BASE_URL}/server-status`);
      if (response.status==200){
        commit(SET_SERVER_STATUS,{isOnline:true,errorMessage:null})
      } else {
        commit(SET_SERVER_STATUS,{isOnline:false,errorMessage:"Can't connect to server"})
      }
    } catch(error) {
      commit(SET_SERVER_STATUS,{isOnline:false,errorMessage:"Connection to the server cannot be established at the moment, Please try again later."})
    }
  },
  async getRoomState(
    { state }: ActionContext<RootState,RootState>,
    payload:{roomId:string}
    ):Promise<RoomState>{
      try{
        const response = await makeHttpRequest<RoomState>(`${BASE_URL}/room-state?roomid=${payload.roomId}`);
        return response;
      } catch (error) {
        console.error(error);
        throw new Error('Error creating room');
      }

  },
  async createRoom(
    { state }: ActionContext<RootState,RootState>,
    payload:{ fen: string, movePatterns?: MovePattern[] }
    ):Promise<string>{
      try{
        const response = await makeHttpRequest<CreateRoomResponse>(`${BASE_URL}/create-room`,'post',JSON.stringify(payload));
        return response.roomId;
      } catch (error) {
        console.error(error);
        throw new Error('Error creating room');
      }
    },

    async deleteRoom(
      { state }: ActionContext<RootState,RootState>,
      payload:{ roomId: string}
      ):Promise<void>{
        try{
          const response = await makeHttpRequest<CreateRoomResponse>(`${BASE_URL}/delete-room`,'post',JSON.stringify(payload));
          return;
        } catch (error) {
          console.error(error);
          throw new Error('Error creating room');
        }
      }
};

export default actions;
