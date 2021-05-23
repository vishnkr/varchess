import { State } from './state';
export declare function createElement(tagName: string): SVGElement;
export declare function renderSvg(state: State, svg: SVGElement, customSvg: SVGElement): void;
export declare function setAttributes(el: SVGElement, attrs: {
    [key: string]: any;
}): SVGElement;
