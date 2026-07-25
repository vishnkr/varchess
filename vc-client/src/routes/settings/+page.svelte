<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { COLOR_THEMES } from '$lib/utils/index';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Dialog from '$lib/components/ui/dialog';
	import { settings, type UserSettings, DEFAULT_SETTINGS } from '$lib/store/settings';
	import { toast } from '$lib/store/alert';
	import { authStore, logout } from '$lib/store/auth';
	import { deleteAccount } from '$lib/api/account';
	import { playSfx, unlockSfx } from '$lib/utils/sfx';

	let draft: UserSettings = { ...DEFAULT_SETTINGS };
	let saving = false;
	let loaded = false;

	let deleteOpen = false;
	let deleteConfirm = '';
	let deleting = false;

	$: username = $authStore.username ?? '';
	$: canDelete = deleteConfirm === username && username.length > 0;

	onMount(async () => {
		const synced = await settings.load({ syncRemote: !!$authStore.accessToken });
		draft = { ...synced };
		prevSoundEnabled = draft.soundEnabled;
		loaded = true;
	});

	function onBoardThemeChange() {
		settings.patch({ boardTheme: draft.boardTheme });
	}

	function onAppThemeChange() {
		settings.patch({ appTheme: draft.appTheme });
	}

	// Apply toggles to the live store immediately (Save still syncs to the server).
	$: if (loaded) {
		settings.patch({
			showCoordinates: draft.showCoordinates,
			highlightLastMove: draft.highlightLastMove,
			showPossibleMoves: draft.showPossibleMoves,
			enablePremove: draft.enablePremove,
			confirmResign: draft.confirmResign,
			hideChat: draft.hideChat,
			soundEnabled: draft.soundEnabled,
			boardTheme: draft.boardTheme,
			appTheme: draft.appTheme
		});
	}

	let prevSoundEnabled = false;
	$: if (loaded && draft.soundEnabled !== prevSoundEnabled) {
		const turningOn = draft.soundEnabled && !prevSoundEnabled;
		prevSoundEnabled = draft.soundEnabled;
		if (turningOn) {
			unlockSfx();
			playSfx('move');
		}
	}

	async function save() {
		saving = true;
		try {
			await settings.save(draft);
			toast.success('Settings saved');
		} catch {
			toast.error('Saved locally; server sync failed');
		} finally {
			saving = false;
		}
	}

	function resetDefaults() {
		draft = { ...DEFAULT_SETTINGS };
		settings.patch(draft);
	}

	function openDeleteDialog() {
		deleteConfirm = '';
		deleteOpen = true;
	}

	async function confirmDeleteAccount() {
		if (!canDelete || deleting) return;
		deleting = true;
		try {
			await deleteAccount(deleteConfirm);
			deleteOpen = false;
			logout();
			toast.success('Your account has been deleted');
			await goto('/login');
		} catch (e) {
			toast.error(e instanceof Error ? e.message : 'Could not delete account');
		} finally {
			deleting = false;
		}
	}
</script>

<svelte:head>
	<title>Settings - Varchess</title>
</svelte:head>

<div class="text-gray-700 dark:text-gray-300 px-4 py-6 max-w-xl mx-auto">
	<h1 class="text-2xl font-bold text-center mb-6">Settings</h1>

	{#if !loaded}
		<p class="text-center text-gray-500">Loading…</p>
	{:else}
		<div class="bg-white dark:bg-gray-800 shadow-md rounded-lg p-6 space-y-8">
			<section class="space-y-4">
				<h2 class="text-lg font-semibold border-b border-gray-200 dark:border-gray-600 pb-2">Board</h2>
				<div class="flex flex-wrap items-center gap-3">
					<label class="text-sm font-medium" for="board-theme">Board theme</label>
					<select
						id="board-theme"
						class="bg-white dark:bg-gray-900 text-black dark:text-white appearance-none cursor-pointer border rounded-md py-1.5 px-3 leading-tight focus:outline-none focus:ring"
						bind:value={draft.boardTheme}
						on:change={onBoardThemeChange}
					>
						{#each Object.keys(COLOR_THEMES) as themeName}
							<option value={themeName}>{themeName}</option>
						{/each}
					</select>
					{#if draft.boardTheme && COLOR_THEMES[draft.boardTheme]}
						<div class="flex items-center gap-2">
							<div
								class="w-6 h-6 rounded-sm border"
								style="background-color: {COLOR_THEMES[draft.boardTheme]
									.lightColor}; border-color: rgba(0, 0, 0, 0.3);"
							/>
							<div
								class="w-6 h-6 rounded-sm border"
								style="background-color: {COLOR_THEMES[draft.boardTheme]
									.darkColor}; border-color: rgba(0, 0, 0, 0.3);"
							/>
						</div>
					{/if}
				</div>

				<div class="flex items-center space-x-2">
					<Checkbox id="show-coords" bind:checked={draft.showCoordinates} />
					<label for="show-coords" class="text-sm">Show board coordinates</label>
				</div>

				<div class="flex items-center space-x-2">
					<Checkbox id="highlight-last" bind:checked={draft.highlightLastMove} />
					<label for="highlight-last" class="text-sm">Highlight last move</label>
				</div>
			</section>

			<section class="space-y-4">
				<h2 class="text-lg font-semibold border-b border-gray-200 dark:border-gray-600 pb-2">Play</h2>
				<div class="flex items-center space-x-2">
					<Checkbox id="show-possible-moves" bind:checked={draft.showPossibleMoves} />
					<label for="show-possible-moves" class="text-sm"
						>Show possible moves after clicking a piece</label
					>
				</div>
				<div class="flex items-center space-x-2">
					<Checkbox id="enable-premove" bind:checked={draft.enablePremove} />
					<label for="enable-premove" class="text-sm">Enable premoves</label>
				</div>
				<div class="flex items-center space-x-2">
					<Checkbox id="confirm-resign" bind:checked={draft.confirmResign} />
					<label for="confirm-resign" class="text-sm">Confirm before resigning</label>
				</div>
			</section>

			<section class="space-y-4">
				<h2 class="text-lg font-semibold border-b border-gray-200 dark:border-gray-600 pb-2">
					Notifications & chat
				</h2>
				<div class="flex items-center space-x-2">
					<Checkbox id="sound-enabled" bind:checked={draft.soundEnabled} />
					<label for="sound-enabled" class="text-sm">Sound effects</label>
				</div>
				<div class="flex items-center space-x-2">
					<Checkbox id="hide-chat" bind:checked={draft.hideChat} />
					<label for="hide-chat" class="text-sm">Hide in-game chat</label>
				</div>
			</section>

			<section class="space-y-4">
				<h2 class="text-lg font-semibold border-b border-gray-200 dark:border-gray-600 pb-2">App</h2>
				<div class="flex flex-wrap items-center gap-3">
					<label class="text-sm font-medium" for="app-theme">Color mode</label>
					<select
						id="app-theme"
						class="bg-white dark:bg-gray-900 text-black dark:text-white appearance-none cursor-pointer border rounded-md py-1.5 px-3 leading-tight focus:outline-none focus:ring"
						bind:value={draft.appTheme}
						on:change={onAppThemeChange}
					>
						<option value="light">Light</option>
						<option value="dark">Dark</option>
					</select>
				</div>
			</section>

			<div class="flex flex-col sm:flex-row gap-2 pt-2">
				<Button on:click={save} disabled={saving} class="flex-1">
					{saving ? 'Saving…' : 'Save Settings'}
				</Button>
				<Button variant="outline" on:click={resetDefaults} class="flex-1">Reset defaults</Button>
			</div>
		</div>

		{#if $authStore.accessToken}
			<section
				class="mt-8 rounded-lg border border-red-300 dark:border-red-900/60 bg-red-50/80 dark:bg-red-950/30 p-6 space-y-3"
			>
				<h2 class="text-lg font-semibold text-red-700 dark:text-red-300">Danger zone</h2>
				<p class="text-sm text-red-800/80 dark:text-red-200/80">
					Permanently delete your account, settings, and templates. Past game history may remain
					anonymized. This cannot be undone.
				</p>
				<Button variant="destructive" on:click={openDeleteDialog}>Delete account</Button>
			</section>
		{/if}
	{/if}
</div>

<Dialog.Root bind:open={deleteOpen}>
	<Dialog.Content class="max-w-md dark:bg-darkbg bg-white text-gray-900 dark:text-white">
		<Dialog.Header>
			<Dialog.Title>Delete account?</Dialog.Title>
			<Dialog.Description class="text-gray-600 dark:text-gray-300">
				This permanently removes <span class="font-semibold">{username}</span>. Type your username
				to confirm.
			</Dialog.Description>
		</Dialog.Header>
		<div class="space-y-2 py-2">
			<label for="delete-confirm" class="text-sm font-medium">Username</label>
			<Input
				id="delete-confirm"
				bind:value={deleteConfirm}
				placeholder={username}
				autocomplete="off"
				disabled={deleting}
			/>
		</div>
		<Dialog.Footer class="flex gap-2 justify-end">
			<Button variant="outline" disabled={deleting} on:click={() => (deleteOpen = false)}
				>Cancel</Button
			>
			<Button variant="destructive" disabled={!canDelete || deleting} on:click={confirmDeleteAccount}>
				{deleting ? 'Deleting…' : 'Delete forever'}
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
