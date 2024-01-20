import { Role, chats, members, wsStore } from './store/stores';

interface UserJoin {
	id: number;
	username: string;
	isHost: boolean;
	role: string;
}
function UserJoinHandler(data: UserJoin) {
	chats.userJoin(data.username);
	members.update((value) => [
		...value,
		{
			id: data.id,
			username: data.username,
			isHost: data.isHost,
			role: Role[data.role as keyof typeof Role]
		}
	]);
}

export const createWebSocket = (wsServerUrl: string) => {
	return new Promise((resolve, reject) => {
		const ws = new WebSocket(wsServerUrl);

		ws.onmessage = function (event) {
			const { type, data } = JSON.parse(event.data);
			switch (type) {
				case 'UserJoin':
					UserJoinHandler(data);
					break;
				case 'UserLeave':
					chats.userLeave(data.username);
					members.update((value) => value.filter((member) => data.username !== member.username));
					break;
				case 'ChatMessage':
					chats.updateChat(data.username, data.content);
					break;
			}
		};

		ws.onerror = function (error) {
			console.error('WebSocket connection error:', error);
			reject(error);
		};

		ws.onopen = function () {
			wsStore.set(ws);
			resolve(ws);
		};
		ws.onclose = function () {
			wsStore.set(null);
			resolve(ws);
		};
	});
};
