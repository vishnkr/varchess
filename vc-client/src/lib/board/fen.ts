import type { configDefaults } from 'vitest/config';
import {
	Color,
	type Dimensions,
	type PiecePositions,
	type PiecePresentInfo,
	type Position
} from './types';

export const convertFenToPosition = (
	fen: string
):
	| {
			dimensions: Dimensions;
			position: Position;
			maxBoardState: PiecePresentInfo[][];
	  }
	| undefined => {
	const dimensions = { ranks: 8, files: 8 };
	const maxBoardState = createEmptyMaxBoardState();
	const fenSplit = fen.split(' ');
	const ranks = fenSplit[0].split('/');
	const rankCount = ranks.length;
	const position: Position = { piecePositions: {}, walls: {} };
	dimensions.ranks = rankCount;
	if (rankCount > 16) {
		return undefined;
	}
	let secDigit,
		colEnd = 0;
	let char;
	let colCount = 0;
	let row = 0;
	let idx = 0;
	for (let i = 0; i < ranks.length; i++) {
		secDigit = 0;
		colCount = 0;
		let j = 0;
		let col = 0;
		while (j < ranks[i].length) {
			char = ranks[i].charAt(j);
			if (char === '.') {
				maxBoardState[i][col] = { isPiecePresent: false, wall: true };
				position.walls[idx] = true;
				colCount += 1;
				idx += 1;
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
						idx += 1;
					}
					colCount += colEnd;
				}
			} else {
				const piece = {
					color: char.toLowerCase() === char ? Color.BLACK : Color.WHITE,
					pieceType: char
				};
				maxBoardState[row][col] = { isPiecePresent: true, piece };
				position.piecePositions[idx] = piece;
				colCount += 1;
				col += 1;
				idx += 1;
			}
			j += 1;
		}
		row += 1;
	}
	dimensions.files = colCount;
	const newPiecePositions: PiecePositions = {};
	const total = dimensions.files * dimensions.ranks;
	/*for (let idx = 0; idx < total; idx++) {
		newPiecePositions[total - 1 - idx] = position.piecePositions[idx];
	}
	
	position.piecePositions = newPiecePositions;*/
	return { dimensions, position, maxBoardState };
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
