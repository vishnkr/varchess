
import { IDimensions, GameRole, ISquareInfo, UserInfo } from "./types";


export class Square{
    squareId?: number;
    x?: number;
    y?: number;
    disabled: boolean;
    squareInfo: ISquareInfo

    constructor({ disabled, squareInfo,squareId,x,y }:{squareInfo: ISquareInfo,disabled?:boolean, squareId?: number, x?: number, y?: number}) {
        this.squareId = squareId;
        this.disabled = disabled ?? false;
        this.x = x;
        this.y = y;
        this.squareInfo = squareInfo;
    }

    updateSquareInfo(newInfo: Partial<ISquareInfo>) {
        Object.assign(this.squareInfo, newInfo);
    }
}

interface IBoardStateProps{
    squares: Square[][];
    dimensions: IDimensions;
    castlingAvailability: string;
    enPassant: string;
    turn: string;
}

export class BoardState implements IBoardStateProps {
    constructor(props: IBoardStateProps) {
      this.squares = props.squares;
      this.dimensions = props.dimensions;
      this.castlingAvailability = props.castlingAvailability;
      this.enPassant = props.enPassant;
      this.turn = props.turn;
      this.initId()
    }
    squares: Square[][];
    dimensions: IDimensions;
    castlingAvailability: string;
    enPassant: string;
    turn: string;
  
    flattenSquares(): Square[] {
      return this.squares.flat();
    }

    initId (){
      let id = 0;
      for(let row=0; row<this.squares.length;row++){
          for(let col=0;col<this.squares[row].length;col++){
              this.squares[row][col].squareId =id;
              id+=1;
          }
      }
  }
}
  

export class User implements UserInfo{
    constructor(
        public username?: string,
        public isAuthenticated: boolean = false,
        public curGameRole: GameRole = 'member'
    ) {}

    isPlayer(){return this.curGameRole==='p1' || this.curGameRole==='p2'}
    setUsername(username:string){this.username = username}

}

  
  export class GameInfo {
    constructor(
      public players: Players,
      public members: Record<string,User>,
      public result?: string
    ) {}
  }
  
  export class Players {
    constructor(public p1: User, public p2: User) {}
    
  }

