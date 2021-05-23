import { State } from './state';
import * as cg from './types';
export declare function bindBoard(s: State, boundsUpdated: () => void): void;
export declare function bindDocument(s: State, boundsUpdated: () => void): cg.Unbind;
