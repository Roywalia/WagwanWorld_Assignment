// // src/lib/api.ts
// const API_BASE_URL = 'http://localhost:8080/api/v1';

// export interface Guest {
// 	id: number;
// 	name: string;
// 	email: string;
// 	phone: string;
// 	status: 'pending' | 'attending' | 'declined';
// 	created_at: string;
// }

// export interface CreateGuestData {
// 	name: string;
// 	email: string;
// 	phone: string;
// 	status?: 'pending' | 'attending' | 'declined';
// }

// // Fixed: use ?status= not ?filter=
// export async function getGuests(status?: string, search?: string): Promise<Guest[]> {
// 	let url = `${API_BASE_URL}/guests`;
// 	const params = new URLSearchParams();
// 	if (status) params.append('status', status);
// 	if (search) params.append('search', search);
// 	if (params.toString()) url += `?${params.toString()}`;

// 	const res = await fetch(url);
// 	if (!res.ok) throw new Error('Failed to fetch guests');
// 	return res.json();
// }

// export async function createGuest(data: CreateGuestData): Promise<Guest> {
// 	const res = await fetch(`${API_BASE_URL}/guests`, {
// 		method: 'POST',
// 		headers: { 'Content-Type': 'application/json' },
// 		body: JSON.stringify(data)
// 	});
// 	if (!res.ok) {
// 		const err = await res.json();
// 		throw new Error(err.error || err.message || 'Failed to create guest');
// 	}
// 	return res.json();
// }

// export async function deleteGuest(id: number): Promise<void> {
// 	const res = await fetch(`${API_BASE_URL}/guests/${id}`, { method: 'DELETE' });
// 	if (!res.ok) throw new Error('Failed to delete guest');
// }

// // function to create a new event
// export async function createEvent(data: {
// 	title: string;
// 	description?: string;
// 	event_date: string;
// 	location?: string;
// }): Promise<any> {
// 	const res = await fetch(`${API_BASE_URL}/events`, {
// 		method: 'POST',
// 		headers: { 'Content-Type': 'application/json' },
// 		body: JSON.stringify(data),
// 	});
// 	if (!res.ok) {
// 		const err = await res.json();
// 		throw new Error(err.error || 'Failed to create event');
// 	}
// 	return res.json();
// }

// export async function getEvents(): Promise<any[]> {
// 	const res = await fetch(`${API_BASE_URL}/events`);
// 	if (!res.ok) throw new Error('Failed to fetch events');
// 	return res.json();
// }