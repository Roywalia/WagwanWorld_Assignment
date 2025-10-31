// src/routes/admin/+page.server.ts
import type { PageServerLoad, Actions } from './$types';
import { fail } from '@sveltejs/kit';
import { z } from 'zod';


// Load data (guests OR events)
export const load = (async ({ url, fetch }: { url: URL; fetch: typeof globalThis.fetch }) => {
	const showEvents = url.searchParams.get('view') === 'events';

	if (showEvents) {
		const res = await fetch('http://localhost:8080/api/v1/events');
		if (!res.ok) throw new Error('Failed to load events');
		const events = await res.json();
		return { view: 'events' as const, events };
	}

	// ---- GUESTS ----
	const status = url.searchParams.get('status') ?? undefined;
	const search = url.searchParams.get('search') ?? undefined;

	const guestsUrl = new URL('http://localhost:8080/api/v1/guests');
	if (status) guestsUrl.searchParams.set('status', status);
	if (search) guestsUrl.searchParams.set('search', search);

	const res = await fetch(guestsUrl);
	if (!res.ok) throw new Error('Failed to load guests');
	const guests = await res.json();

	return { view: 'guests' as const, guests };
}) satisfies PageServerLoad;

export const actions = {
	// ---------- CREATE EVENT ----------
	// src/routes/admin/+page.server.ts
	createEvent: async ({ request }) => {
	const data = await request.formData();

	const title = data.get('title') as string;
	const description = data.get('description') as string;
	const eventDate = data.get('event_date') as string;
	const location = data.get('location') as string;

	if (!title?.trim()) return fail(400, { error: 'Title is required' });
	if (!eventDate) return fail(400, { error: 'Date is required' });

	const dateObj = new Date(eventDate);
	if (isNaN(dateObj.getTime())) return fail(400, { error: 'Invalid date' });

	const payload = {
		title: title.trim(),
		description: description?.trim() || undefined,
		event_date: dateObj.toISOString(),
		location: location?.trim() || undefined
	};

	const res = await fetch('http://localhost:8080/api/v1/events', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify(payload),
	});

	if (!res.ok) {
		const err = await res.json().catch(() => ({}));
		return fail(res.status, { error: err.error ?? 'Failed to create event' });
	}

	return { success: true };
},

	// ---------- DELETE GUEST ----------
	deleteGuest: async ({ request }: { request: Request }) => {
		const data = await request.formData();
		const idStr = data.get('id') as string ?? '';

		if (!idStr) {
			return fail(400, { error: 'Guest ID missing' });
		}

		const id = parseInt(idStr, 10);
		if (isNaN(id) || id <= 0) {
			return fail(400, { error: 'Invalid Guest ID' });
		}

		const res = await fetch(`http://localhost:8080/api/v1/guests/${id}`, {
			method: 'DELETE',
		});

		if (!res.ok) {
			return fail(res.status, { error: 'Failed to delete guest' });
		}

		return { success: true };
	},
} satisfies Actions;