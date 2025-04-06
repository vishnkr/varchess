// src/lib/store/theme.ts
import { writable } from 'svelte/store';

export const theme = writable<'light' | 'dark'>('dark');
