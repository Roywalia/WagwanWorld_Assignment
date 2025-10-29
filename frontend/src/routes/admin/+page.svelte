<script lang="ts">
	import { onMount } from 'svelte';
	import { getGuests, deleteGuest, createEvent, getEvents, type Guest } from '$lib/api';
	import Swal from 'sweetalert2';

	/* -------------------------------------------------
	 *  VIEW TOGGLE – boolean proxy
	 * ------------------------------------------------- */
	let showEvents = $state(false);

	// === Guests ===
	let guests: Guest[] = $state([]);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let statusFilter = $state('');
	let searchQuery = $state('');

	// === Events ===
	let events: any[] = $state([]);
	let loadingEvents = $state(true);
	let errorEvents = $state<string | null>(null);

	// === Event Form ===
	let showEventForm = $state(false);
	let eventTitle = $state('');
	let eventDesc = $state('');
	let eventDate = $state('');
	let eventLocation = $state('');
	let eventError = $state<string | null>(null);
	let eventSubmitting = $state(false);

	onMount(() => {
		if (!showEvents) loadGuests();
		else loadEvents();
	});

	let searchTimeout: any;
	$effect(() => {
		clearTimeout(searchTimeout);
		searchTimeout = setTimeout(() => {
			if (!showEvents) loadGuests();
		}, 300);
	});

	// Reload when view toggles
	$effect(() => {
		if (showEvents) loadEvents();
		else loadGuests();
	});

	// Load guests
	async function loadGuests() {
		try {
			loading = true;
			error = null;
			guests = await getGuests(statusFilter || undefined, searchQuery || undefined);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load guests';
		} finally {
			loading = false;
		}
	}

	// Load events
	async function loadEvents() {
		try {
			loadingEvents = true;
			errorEvents = null;
			events = await getEvents();
		} catch (e) {
			errorEvents = e instanceof Error ? e.message : 'Failed to load events';
		} finally {
			loadingEvents = false;
		}
	}

	// Create event
	async function handleEventSubmit(e: Event) {
		e.preventDefault();
		eventError = null;

		if (!eventTitle.trim()) return eventError = 'Title is required';
		if (!eventDate) return eventError = 'Date is required';

		try {
			eventSubmitting = true;
			await createEvent({
				title: eventTitle,
				description: eventDesc,
				event_date: new Date(eventDate).toISOString(),
				location: eventLocation
			});
			resetEventForm();
			await loadEvents();
			Swal.fire('Success', 'Event created!', 'success');
		} catch (e: any) {
			eventError = e.message || 'Failed to create event';
		} finally {
			eventSubmitting = false;
		}
	}

	function resetEventForm() {
		eventTitle = '';
		eventDesc = '';
		eventDate = '';
		eventLocation = '';
		showEventForm = false;
	}

	async function handleDelete(id: number) {
		if (!confirm('Delete this guest?')) return;
		try {
			await deleteGuest(id);
			await loadGuests();
		} catch (e) {
			alert(e instanceof Error ? e.message : 'Failed to delete');
		}
	}

	function getBadge(status: string) {
		const map: Record<string, string> = {
			attending: 'bg-green-100 text-green-800',
			declined: 'bg-red-100 text-red-800',
			pending: 'bg-yellow-100 text-yellow-800'
		};
		return map[status] || map.pending;
	}

	function formatDate(iso: string) {
		return new Date(iso).toLocaleString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric',
			hour: 'numeric',
			minute: '2-digit'
		});
	}

	function safeString(value: any): string {
		return value && typeof value === 'object' && value.String ? value.String : (value || '-');
	}
</script>

<div class="min-h-screen bg-gray-50 py-8">
	<div class="max-w-6xl mx-auto px-4">
		<div class="mb-8">
			<h1 class="text-3xl font-bold text-gray-900 mb-2">Event Guest List Manager</h1>
			<p class="text-gray-600">Manage events and guest RSVPs</p>
		</div>

		<!-- Filters Bar: Toggle (left) | Dropdown + Search + Button (right) -->
		<div class="bg-white rounded-lg shadow-sm p-4 mb-6 flex flex-wrap gap-4 items-center justify-between">
			<!-- LEFT: Toggle -->
			<div class="flex items-center gap-3">
				<span class="text-sm font-medium text-gray-700">Guests</span>
				<label class="relative inline-flex items-center cursor-pointer">
					<input type="checkbox" bind:checked={showEvents} class="sr-only peer" />
					<div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-purple-300 rounded-full peer 
						peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] 
						after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all 
						peer-checked:bg-purple-600"></div>
				</label>
				<span class="text-sm font-medium text-gray-700">Events</span>
			</div>

			<!-- RIGHT: Filters + Add Event Button -->
			<div class="flex gap-4 items-center flex-1 justify-end">
				{#if !showEvents}
					<select bind:value={statusFilter} class="border rounded-md px-3 py-2">
						<option value="">All</option>
						<option value="pending">Pending</option>
						<option value="attending">Attending</option>
						<option value="declined">Declined</option>
					</select>

					<div class="relative max-w-md w-full min-w-[200px]">
						<input 
							type="text" 
							bind:value={searchQuery}
							placeholder="Search by name or email..."
							class="customInput w-full"
						/>
						{#if searchQuery}
							<button 
								onclick={() => searchQuery = ''}
								class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
							>
								Clear
							</button>
						{/if}
					</div>
				{/if}

				{#if showEvents}
					<button 
						onclick={() => showEventForm = !showEventForm}
						class="bg-purple-600 text-white px-4 py-2 rounded-md hover:bg-purple-700 whitespace-nowrap"
					>
						{showEventForm ? 'Cancel' : '+ Add Event'}
					</button>
				{/if}
			</div>
		</div>

		<!-- Add Event Form -->
		{#if showEventForm && showEvents}
			<div class="bg-white p-6 rounded-lg shadow-sm mb-6">
				<h2 class="text-xl font-semibold mb-4">Add New Event</h2>
				{#if eventError}
					<div class="bg-red-50 text-red-700 p-3 rounded mb-4">{eventError}</div>
				{/if}
				<form onsubmit={handleEventSubmit} class="grid md:grid-cols-2 gap-4">
					<input bind:value={eventTitle} placeholder="Event Title *" required class="customInput" />
					<input bind:value={eventLocation} placeholder="Location" class="customInput" />
					<input bind:value={eventDate} type="datetime-local" required class="customInput" min={new Date().toISOString().slice(0, 16)}/>
					<textarea bind:value={eventDesc} placeholder="Description (optional)" rows="2" class="customInput md:col-span-2"></textarea>
					<div class="md:col-span-2 flex gap-2">
						<button type="submit" disabled={eventSubmitting} class="bg-purple-600 text-white px-6 py-2 rounded-md">
							{eventSubmitting ? 'Creating...' : 'Create Event'}
						</button>
						<button type="button" onclick={() => showEventForm = false} class="bg-gray-200 px-6 py-2 rounded-md">
							Cancel
						</button>
					</div>
				</form>
			</div>
		{/if}

		<!-- Events Table -->
		{#if showEvents}
			<div class="bg-white rounded-lg shadow-sm overflow-hidden">
				{#if loadingEvents}
					<p class="p-8 text-center text-gray-500">Loading events...</p>
				{:else if errorEvents}
					<p class="p-8 text-center text-red-600">{errorEvents}</p>
				{:else if events.length === 0}
					<p class="p-8 text-center text-gray-500">No events yet.</p>
				{:else}
					<table class="w-full">
						<thead class="bg-gray-50">
							<tr>
								<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
								<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Title</th>
								<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Date</th>
								<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Location</th>
								<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">RSVPs</th>
							</tr>
						</thead>
						<tbody class="divide-y">
							{#each events as e (e.id)}
								<tr class="hover:bg-gray-50">
									<td class="px-6 py-4 text-sm font-medium">#{e.id}</td>
									<td class="px-6 py-4 text-sm font-medium">{e.title}</td>
									<td class="px-6 py-4 text-sm text-gray-600">{formatDate(e.event_date)}</td>
									<td class="px-6 py-4 text-sm text-gray-600">{safeString(e.location)}</td>
									<td class="px-6 py-4">
										<span class="px-2 py-1 text-xs font-semibold rounded-full bg-blue-100 text-blue-800">
											{e.rsvps}
										</span>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				{/if}
			</div>
		{/if}

		<!-- Guest Table -->
		{#if !showEvents}
			<div class="bg-white rounded-lg shadow-sm overflow-hidden">
				{#if loading}
					<p class="p-8 text-center text-gray-500">Loading...</p>
				{:else if error}
					<p class="p-8 text-center text-red-600">{error}</p>
				{:else if guests.length === 0}
					<p class="p-8 text-center text-gray-500">
						{searchQuery || statusFilter ? 'No guests match your search.' : 'No guests yet.'}
					</p>
				{:else}
					<table class="w-full">
						<thead class="bg-gray-50">
							<tr>
								<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
								<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Email</th>
								<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Phone</th>
								<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
								<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
							</tr>
						</thead>
						<tbody class="divide-y">
							{#each guests as g (g.id)}
								<tr class="hover:bg-gray-50">
									<td class="px-6 py-4 text-sm font-medium">{g.name}</td>
									<td class="px-6 py-4 text-sm text-gray-600">{g.email}</td>
									<td class="px-6 py-4 text-sm text-gray-600">{g.phone || '-'}</td>
									<td class="px-6 py-4">
										<span class={`px-2 py-1 text-xs font-semibold rounded-full ${getBadge(g.status)}`}>
											{g.status}
										</span>
									</td>
									<td class="px-6 py-4">
										<button onclick={() => handleDelete(g.id)} class="text-red-600 hover:text-red-800">Delete</button>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				{/if}
			</div>

			<!-- Stats -->
			{#if !loading && guests.length > 0}
				<div class="mt-6 bg-white p-4 rounded-lg shadow-sm text-sm">
					<span>Total: <strong>{guests.length}</strong></span> |
					<span>Attending: <strong class="text-green-600">{guests.filter(g => g.status === 'attending').length}</strong></span> |
					<span>Pending: <strong class="text-yellow-600">{guests.filter(g => g.status === 'pending').length}</strong></span> |
					<span>Declined: <strong class="text-red-600">{guests.filter(g => g.status === 'declined').length}</strong></span>
				</div>
			{/if}
		{/if}
	</div>
</div>

<style>
  :global(.swal2-popup) { 
    font-family: system-ui, -apple-system, sans-serif; 
  }

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

  :global(.customInput::before) {
    content: '';
    @apply absolute inset-0 rounded-lg opacity-0 transition-opacity duration-300;
    background: linear-gradient(45deg, #6366f1, #8b5cf6, #ec4899);
    z-index: -1;
  }

  :global(.customInput:focus::before) {
    @apply opacity-100;
  }
</style>