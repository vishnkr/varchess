
type CoordDirection = "rank"|"file";

export class Coordinates{
    private readonly _coordDirection: CoordDirection;

    constructor(props:{coordDirection:CoordDirection}){
        this._coordDirection = props.coordDirection

    }
}