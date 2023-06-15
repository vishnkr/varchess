import { Color, type Dimensions, type PiecePresentInfo } from './types';

export const convertFenToPosition = (
	fen: string
):
	| {
			dimensions: Dimensions;
			maxBoardState: PiecePresentInfo[][];
	  }
	| undefined => {
	const dimensions = { ranks: 8, files: 8 };
	const maxBoardState = createEmptyMaxBoardState();
	const fenSplit = fen.split(' ');
	const ranks = fenSplit[0].split('/');
	const rankCount = ranks.length;
	dimensions.ranks = rankCount;
	if (rankCount > 16) {
		return undefined;
	}
	let secDigit,
		colEnd = 0;
	let char;
	let colCount = 0;
	let row = 0;
	for (let i = 0; i < ranks.length; i++) {
		secDigit = 0;
		colCount = 0;
		let j = 0;
		let col = 0;
		while (j < ranks[i].length) {
			char = ranks[i].charAt(j);
			if (char === '.') {
				maxBoardState[i][col] = { isPiecePresent: false, disabled: true };
				colCount += 1;
				col += 1;
			} else if (/\d/.test(char)) {
				if (j + 1 < ranks[i].length && /\d/.test(ranks[i].charAt(j + 1))) {
					secDigit = parseInt(char);
				} else {
					if (secDigit != 0) {
						colEnd = secDigit * 10 + parseInt(char);
					} else {
						colEnd = parseInt(char);
					}
					for (let empty = 0; empty < colEnd; empty++) {
						maxBoardState[row][col] = { isPiecePresent: false };
						col += 1;
					}
					colCount += colEnd;
				}
			} else {
				const piece = {
					color: char.toLowerCase() === char ? Color.BLACK : Color.WHITE,
					pieceType: char
				};
				maxBoardState[row][col] = { isPiecePresent: true, piece };
				colCount += 1;
				col += 1;
			}
			j += 1;
		}
		row += 1;
	}
	dimensions.files = colCount;

	return { dimensions, maxBoardState };
};

export const createEmptyMaxBoardState = (): PiecePresentInfo[][] => {
	const maxBoard: PiecePresentInfo[][] = [];
	for (let row = 0; row < 16; row++) {
		const boardRow = [];
		for (let col = 0; col < 16; col++) {
			boardRow.push({ isPiecePresent: false });
		}
		maxBoard.push(boardRow);
	}
	return maxBoard;
};
