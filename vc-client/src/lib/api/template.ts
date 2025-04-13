// $lib/api/templates.ts
import { CORE_URL_PROTECTED } from './config';
import type { Template } from '$lib/types';
import { customFetch } from './fetch';

export async function fetchTemplates(page = 1, pageSize = 10): Promise<{
	items: Template[];
	totalCount: number;
	page: number;
	pageSize: number;
	totalPages: number;
}> {
	const res = await customFetch(`${CORE_URL_PROTECTED}/templates?page=${page}&pageSize=${pageSize}`);
	if (!res.ok) throw new Error('Failed to customFetch templates');
	return res.json();
}

export async function getTemplate(id: string): Promise<Template> {
	const res = await customFetch(`${CORE_URL_PROTECTED}/templates/${id}`);
	if (!res.ok) throw new Error('Failed to customFetch template');
	return res.json();
}

export async function createTemplate(template: Omit<Template, '_id'>): Promise<Template> {
	const res = await customFetch(`${CORE_URL_PROTECTED}/templates`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(template)
	});
	if (!res.ok) throw new Error('Failed to create template');
	return res.json();
}

export async function updateTemplate(id: string, updates: Partial<Template>): Promise<void> {
	const res = await customFetch(`${CORE_URL_PROTECTED}/templates/${id}`, {
		method: 'PUT',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(updates)
	});
	if (!res.ok) throw new Error('Failed to update template');
	return res.json();
}

export async function deleteTemplate(id: string): Promise<void> {
	const res = await customFetch(`${CORE_URL_PROTECTED}/templates/${id}`, {
		method: 'DELETE'
	});
	if (!res.ok) throw new Error('Failed to delete template');
}
