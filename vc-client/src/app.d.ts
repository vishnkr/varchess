// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare global {
	namespace App {
		type PocketBase = import('pocketbase').default;
		interface Locals {
			pb?: PocketBase;
			user?: Record<string, T> | null;
		}
		interface PageData {
			session: Session | null;
		}
		// interface Error {}
		// interface Platform {}
	}
}

export {};
