import { Coordinates } from "./coords";
import chessStyle from "./chess-style.css?inline"
import { Board } from "./board";

interface BoardConfig{
    dimensions: {
        width: number,
        height: number
    },
    squareColor?:{
        light: string,
        dark: string
    },
}

export class BoardElement extends HTMLElement{
    private _rankCoordinates : Coordinates;
    private _fileCoordinates: Coordinates;
    private _el_shadow: ShadowRoot;
    private _el_style: HTMLStyleElement;
    private _board_wrapper: HTMLDivElement;
    private _chessboard: Board; 
    constructor(){
        super();
        this._board_wrapper = document.createElement('div');
        this._board_wrapper.classList.add("wrapper");
        this._el_style = document.createElement('style');
        this._el_shadow = this.attachShadow({ mode: "open" });
        this._el_style.textContent = chessStyle; 
        this._el_shadow.appendChild(this._el_style);
        this._rankCoordinates = new Coordinates({coordDirection:"rank"});
        this._fileCoordinates = new Coordinates({coordDirection:"file"});

        
        this._chessboard = new Board(
            {
                shadow: this._el_shadow,
                dispatchEvent:(e) => this.dispatchEvent(e),
            }
        )
        this._board_wrapper.appendChild(this._chessboard.element);
        this._el_shadow.appendChild(this._board_wrapper);
    }

    static get observedAttributes() {
        return ["fen"];
      }
      
    connectedCallback() {

    }

    get fen(){
        return this.getAttribute("fen") || 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR';
    }

    set fen(value){
        console.log('setting fen',value,this.fen)
    }

    attributeChangedCallback(
        name: string,
        _: string | null,
        newValue: string | null
      ) {
        if (name === "fen" && newValue) {
          this.fen = newValue;
          console.log('next',this.fen)
        }
      }


    


}
declare global {
    interface HTMLElementTagNameMap {
      "board-element": BoardElement;
    }
}
customElements.get('board-element') || customElements.define('board-element', BoardElement);