import { ENVIRONMENT, POSTGRES_PORT, SERVER_HOST, SERVER_PORT } from '$env/static/private';

export const environmentType = ENVIRONMENT;
export const serverBase = SERVER_HOST;
export const serverPort = SERVER_PORT;
export const dbPort = POSTGRES_PORT;
export const wsServerUrl = `${
	environmentType === 'production' ? serverBase : 'localhost'
}:${serverPort}`;
export const apiServerUrl = `${
	environmentType === 'production' ? serverBase : 'localhost'
}:${serverPort}`;
