<script lang="ts">
	import { onMount } from 'svelte';
	import Swal from 'sweetalert2';

	/* ---------- STATE ---------- */
	let events: any[] = [];
	let eventId = 0;
	let totalRsvps = 0;

	let form = {
		name: '',
		email: '',
		phone: '',
		status: 'attending',   // attending | maybe | declined
		notes: '',
		plus_ones: 0,
		dietary: ''
	};

	let loading = true;
	let submitting = false;

	/* ---------- LOAD EVENTS ---------- */
	onMount(async () => {
		await loadEvents();
	});

	async function loadEvents() {
		try {
			const res = await fetch('http://localhost:8080/api/v1/events');
			if (!res.ok) throw new Error('Failed to load events');
			events = await res.json();
			totalRsvps = events.reduce((s: number, e: any) => s + (e.rsvps ?? 0), 0);
		} catch (e) {
			Swal.fire('Error', 'Could not load events', 'error');
		} finally {
			loading = false;
		}
	}

	/* ---------- SUBMIT RSVP ---------- */
	async function submit(e: Event) {
		e.preventDefault();

		if (!form.name.trim() || !form.email.trim() || !eventId) {
			Swal.fire('Oops!', 'Name, email and event are required', 'warning');
			return;
		}

		submitting = true;

		const payload = {
			name: form.name.trim(),
			email: form.email.trim(),
			phone: form.phone.trim(),
			rsvp_status: form.status,
			notes: form.notes.trim(),
			plus_ones: form.plus_ones,
			dietary_restrictions: form.dietary.trim()
		};

		try {
			const res = await fetch(`http://localhost:8080/api/v1/events/${eventId}/rsvp`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(payload)
			});

			const data = await res.json();

			if (res.ok) {
				Swal.fire({
					icon: 'success',
					title: 'Thank You!',
					text: data.message,
					confirmButtonColor: '#6366f1'
				});

				// Reset form
				form = {
					name: '',
					email: '',
					phone: '',
					status: 'attending',
					notes: '',
					plus_ones: 0,
					dietary: ''
				};
				eventId = 0;

				// Refresh total RSVPs
				await loadEvents();
			} else {
				Swal.fire({
					icon: 'error',
					title: 'Error',
					text: data.message || 'Something went wrong',
					confirmButtonColor: '#ef4444'
				});
			}
		} catch {
			Swal.fire('Network Error', 'Please check your connection', 'error');
		} finally {
			submitting = false;
		}
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-indigo-50 to-purple-50 p-6">
	<div class="max-w-2xl mx-auto bg-white rounded-2xl shadow-xl p-8">
		<div class="text-center mb-8">
			<h1 class="text-3xl font-bold text-indigo-800 mb-2">RSVP to Our Event</h1>
		</div>

		{#if loading}
			<p class="text-center text-gray-600">Loading events…</p>
		{:else}
			<form on:submit={submit} class="space-y-5">
				<!-- Event selector -->
				<select bind:value={eventId} class="w-full p-3 border rounded-lg customInput" required>
					<option value={0}>Select Event</option>
					{#each events as e}
						<option value={e.id}>{e.display}</option>
					{/each}
				</select>

				<!-- Name -->
				<div class="relative">
					<input bind:value={form.name} placeholder=" " class="customInput" required />
					<label>Full Name *</label>
				</div>

				<!-- Email -->
				<div class="relative">
					<input bind:value={form.email} type="email" placeholder=" " class="customInput" required />
					<label>Email *</label>
				</div>

				<!-- Phone -->
				<div class="relative">
					<input bind:value={form.phone} placeholder=" " class="customInput" />
					<label>Phone</label>
				</div>

				<!-- RSVP Status -->
				<div>
					<label class="block font-medium mb-1">RSVP Status *</label>
					<select bind:value={form.status} class="w-full p-3 border rounded-lg customInput">
						<option value="attending">Yes, I’ll be there!</option>
						<option value="maybe">Maybe</option>
						<option value="declined">Sorry, can’t make it</option>
					</select>
				</div>

				<!-- Plus Ones -->
				<div class="relative">
					<input bind:value={form.plus_ones} type="number" min="0" placeholder=" " class="customInput" />
					<label>Plus Ones</label>
				</div>

				<!-- Dietary -->
				<textarea
					bind:value={form.dietary}
					placeholder="Dietary Restrictions"
					class="w-full p-3 border rounded-lg customInput"
					rows="2"
				></textarea>

				<!-- Notes -->
				<textarea
					bind:value={form.notes}
					placeholder="Notes"
					class="w-full p-3 border rounded-lg customInput"
					rows="3"
				></textarea>

				<!-- Submit -->
				<button
					type="submit"
					disabled={submitting}
					class="w-full bg-indigo-600 text-white py-3 rounded-lg font-bold hover:bg-indigo-700 transition disabled:opacity-70"
				>
					{submitting ? 'Submitting…' : 'Submit RSVP'}
				</button>

				<!-- Total RSVPs -->
				<div class="text-center mt-6">
					<p class="text-indigo-600 font-semibold">Total RSVPs: {totalRsvps}</p>
				</div>
			</form>
		{/if}
	</div>
</div>

<style>
	/* SweetAlert2 font */
	:global(.swal2-popup) {
		font-family: system-ui, -apple-system, sans-serif;
	}

	/* Floating label inputs */
	:global(.customInput) {
		@apply w-full p-2 border rounded-lg transition-all duration-300 ease-in-out;
		border: 2px solid #e2e8f0;
		background: #ffffff;
		font-size: 1rem;
		color: #1e293b;
		outline: none;
		position: relative;
		z-index: 1;
	}

	:global(.customInput:focus) {
		@apply ring-4 ring-indigo-200 ring-opacity-50;
		border-color: #6366f1;
		transform: scale(1.02);
		box-shadow: 0 0 0 4px rgba(99, 102, 241, 0.2),
			0 10px 25px -5px rgba(99, 102, 241, 0.15);
	}

	:global(.customInput + label) {
		@apply absolute left-3 -top-2.5 px-1 bg-white text-sm text-gray-600 transition-all duration-300;
		pointer-events: none;
		z-index: 10;
	}

	:global(.customInput:focus + label),
	:global(.customInput:not(:placeholder-shown) + label) {
		@apply -top-2.5 text-indigo-600;
	}

	/* Gradient background on focus */
	:global(.customInput::before) {
		content: '';
		@apply absolute inset-0 rounded-lg opacity-0 transition-opacity duration-300;
		background: linear-gradient(45deg, #6366f1, #8b5cf6, #ec4899);
		z-index: -1;
	}
	:global(.customInput:focus::before) {
		@apply opacity-100;
	}

	/* Hide number spinner */
	:global(input[type='number']) {
		-moz-appearance: textfield;
	}
	:global(input[type='number'])::-webkit-outer-spin-button,
	:global(input[type='number'])::-webkit-inner-spin-button {
		-webkit-appearance: none;
		margin: 0;
	}
</style>