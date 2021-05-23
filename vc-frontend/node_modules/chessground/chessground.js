import { start } from './api';
import { configure } from './config';
import { defaults } from './state';
import { renderWrap } from './wrap';
import * as events from './events';
import { render, updateBounds } from './render';
import * as svg from './svg';
import * as util from './util';
export function Chessground(element, config) {
    const maybeState = defaults();
    configure(maybeState, config || {});
    function redrawAll() {
        const prevUnbind = 'dom' in maybeState ? maybeState.dom.unbind : undefined;
        // compute bounds from existing board element if possible
        // this allows non-square boards from CSS to be handled (for 3D)
        const relative = maybeState.viewOnly && !maybeState.drawable.visible, elements = renderWrap(element, maybeState, relative), bounds = util.memo(() => elements.board.getBoundingClientRect()), redrawNow = (skipSvg) => {
            render(state);
            if (!skipSvg && elements.svg)
                svg.renderSvg(state, elements.svg, elements.customSvg);
        }, boundsUpdated = () => {
            bounds.clear();
            updateBounds(state);
            if (elements.svg)
                svg.renderSvg(state, elements.svg, elements.customSvg);
        };
        const state = maybeState;
        state.dom = {
            elements,
            bounds,
            redraw: debounceRedraw(redrawNow),
            redrawNow,
            unbind: prevUnbind,
            relative,
        };
        state.drawable.prevSvgHash = '';
        redrawNow(false);
        events.bindBoard(state, boundsUpdated);
        if (!prevUnbind)
            state.dom.unbind = events.bindDocument(state, boundsUpdated);
        state.events.insert && state.events.insert(elements);
        return state;
    }
    return start(redrawAll(), redrawAll);
}
function debounceRedraw(redrawNow) {
    let redrawing = false;
    return () => {
        if (redrawing)
            return;
        redrawing = true;
        requestAnimationFrame(() => {
            redrawNow();
            redrawing = false;
        });
    };
}
//# sourceMappingURL=chessground.js.map