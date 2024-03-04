export function camelToSnake(obj: any): any {
	if (obj !== null && typeof obj === 'object') {
		if (Array.isArray(obj)) {
			return obj.map((item) => camelToSnake(item));
		} else {
			return Object.fromEntries(
				Object.entries(obj).map(([key, value]) => [
					key.replace(/[A-Z]/g, (match) => `_${match.toLowerCase()}`),
					camelToSnake(value)
				])
			);
		}
	} else {
		return obj;
	}
}

export const COLOR_THEMES: Record<string, { lightColor: string; darkColor: string }> = {
	Default: { lightColor: 'hsl(51deg 24% 84%)', darkColor: 'hsl(145deg 32% 44%)' },
	Brown: { lightColor: 'hsl(36, 81%, 84%)', darkColor: 'hsl(25, 31%, 51%)' },
	Aqua: { lightColor: 'hsl(197, 34%, 83%)', darkColor: 'hsl(217, 68%, 52%)' },
	Classic: { lightColor: 'hsl(0, 0%, 100%)', darkColor: 'hsl(0, 0%, 45%)' },
	Candy: { lightColor: 'hsl(314, 100%, 90%)', darkColor: 'hsl(328, 100%, 55%)' }
};