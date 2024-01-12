export function camelToSnake(obj: any): any {
  if (obj !== null && typeof obj === 'object') {
    if (Array.isArray(obj)) {
      return obj.map((item) => camelToSnake(item));
    } else {
      return Object.fromEntries(
        Object.entries(obj).map(([key, value]) => [
          key.replace(/[A-Z]/g, (match) => `_${match.toLowerCase()}`),
          camelToSnake(value),
        ])
      );
    }
  } else {
    return obj;
  }
}