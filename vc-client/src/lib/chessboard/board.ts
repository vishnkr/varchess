import { generateSquaresMap } from "./utils";
import type { SquareNotation } from "./utils";

export enum Color{
    BLACK = "black",
    WHITE = "white",
};

export interface IPiece{
    pieceType:string,
    color: Color
}

export class Square{
    private readonly _tableDataEl: HTMLTableCellElement;
    private readonly _sqContentEl: HTMLDivElement;
    private readonly _slotWrap:HTMLElement;
    private readonly _slotEl: HTMLSlotElement;

    private _isEditable =false;
    private _piece?: IPiece
    private _sqNotation: SquareNotation;
    private _hover = false;
    constructor(props:{sqNotation:SquareNotation,containerEl:HTMLElement}){
        this._tableDataEl = document.createElement('td');
        this._slotWrap = document.createElement('div');
        this._sqContentEl = document.createElement('div');
        this._slotEl = document.createElement('slot');

        this._slotWrap.appendChild(this._slotEl)

        this._sqNotation = props.sqNotation;
        props.containerEl.appendChild(this._tableDataEl);
    }
}

export class Board{
    private readonly _squares: Square[];
    private readonly _turn: Color;
    private readonly _tableEl: HTMLElement;
    private readonly _shadow:ShadowRoot;
    private readonly _eventDispatch: <T>(e: CustomEvent<T>) => void;
    /*private _clickHandler: (e: MouseEvent) => void;
    private _rightClickHandler: (e:MouseEvent)=>void;
    private _pointerDownHandler: (e:PointerEvent)=>void;
    private _pointerUpHandler:(e:PointerEvent)=>void;*/
    
    constructor(props:{
        shadow:ShadowRoot,
        dispatchEvent: <T>(e: CustomEvent<T>) => void
    }){
        const squareMap = generateSquaresMap(8,8)
        this._tableEl = document.createElement('table');
        this._tableEl.classList.add("board");
        this._eventDispatch = dispatchEvent;
        this._squares=[];
        this._shadow = props.shadow;
        this._turn = Color.WHITE;
        for(let i=0;i<8;i++){
            const row = document.createElement("tr");
            for(let j=0;j<8;j++){
                this._squares.push(new Square({sqNotation:squareMap.idToSqMap[i],containerEl:row}));
            }
            
        }
    }

    get element() {
        return this._tableEl;
      }
}