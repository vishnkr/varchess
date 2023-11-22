export function handleKeyDown(event: KeyboardEvent, callback: () => void) {
	if (event.key === 'Enter') {
		callback();
	}
}
