import type {
	CoordinatetoIDMap,
	Dimensions,
	DisabledSquares,
	IPiece,
	Position,
	PiecePresentInfo,
	SquareIdx,
	SquareInfo,
	SquareNotation,
	PiecePositions
} from './types';

export const generateSquareMaps = (dimensions: Dimensions, isFlipped: boolean) => {
	const columns = [
		'a',
		'b',
		'c',
		'd',
		'e',
		'f',
		'g',
		'h',
		'i',
		'j',
		'k',
		'l',
		'm',
		'n',
		'o',
		'p'
	].slice(0, dimensions.files);
	const rows = Array.from({ length: dimensions.ranks }, (_, i) => i);

	//const sqToIdMap: { [key: string]: number } = {};
	//const idToSqMap: { [key: number]: SquareNotation } = {};
	const coordToIdMap: CoordinatetoIDMap = {};
	let squareIndex = 0;

	const squares: Record<SquareIdx, SquareInfo> = {};
	for (const row of rows) {
		columns.forEach((column, colIdx) => {
			//sqToIdMap[squareNotation] = squareIndex;
			//idToSqMap[squareIndex] = squareNotation;
			coordToIdMap[`${row}:${column}`] = squareIndex;
			squares[squareIndex] = {
				gridX: isFlipped ?  dimensions.ranks - row : row + 1,
				gridY: colIdx + 1,
				squareIndex,
				squareNotation: (isFlipped
					? `${column}${dimensions.ranks - row}`
					: `${column}${row + 1}`) as SquareNotation,
				row,
				column: colIdx
			};
			squareIndex++;
		});
	}
	//console.log(squares)
	return { /*sqToIdMap,idToSqMap,*/ coordToIdMap, squares };
};

export const updatePiecePositionsFromMaxBoard = (
	maxBoardState: PiecePresentInfo[][],
	dimensions: Dimensions
): Position => {
	const piecePositions: PiecePositions = {};
	const disabled: DisabledSquares = {};
	for (let row = 0; row < dimensions.ranks; row++) {
		for (let col = 0; col < dimensions.files; col++) {
			if (maxBoardState[row][col].isPiecePresent) {
				piecePositions[row * dimensions.files + col] = maxBoardState[row][
					col
				].piece as IPiece;
			} else if (maxBoardState[row][col].disabled) {
				disabled[row * dimensions.files + col] = true;
			}
		}
	}
	return { piecePositions, disabled };
};
