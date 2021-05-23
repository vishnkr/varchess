import { State } from './state';
import * as cg from './types';
export interface DragCurrent {
    orig: cg.Key;
    piece: cg.Piece;
    origPos: cg.NumberPair;
    pos: cg.NumberPair;
    started: boolean;
    element: cg.PieceNode | (() => cg.PieceNode | undefined);
    newPiece?: boolean;
    force?: boolean;
    previouslySelected?: cg.Key;
    originTarget: EventTarget | null;
}
export declare function start(s: State, e: cg.MouchEvent): void;
export declare function dragNewPiece(s: State, piece: cg.Piece, e: cg.MouchEvent, force?: boolean): void;
export declare function move(s: State, e: cg.MouchEvent): void;
export declare function end(s: State, e: cg.MouchEvent): void;
export declare function cancel(s: State): void;
