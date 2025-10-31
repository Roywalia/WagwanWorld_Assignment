<script lang="ts">
	import { enhance } from '$app/forms';
	import { invalidateAll } from '$app/navigation';
	import type { ActionResult } from '@sveltejs/kit';
	import Swal from 'sweetalert2';

	const { data } = $props<{ data: any }>();

	let showEvents = $state(data.view === 'events');
	let guests = $state<any[]>([]);
	let events = $state<any[]>([]);
	let showEventForm = $state(false);
	let eventTitle = $state('');
	let eventDesc = $state('');
	let eventDate = $state('');
	let eventLocation = $state('');
	let eventError = $state<string | null>(null);
	let eventSubmitting = $state(false);
	let statusFilter = $state('');
	let searchQuery = $state('');

	//LOAD DATA BASED ON CURRENT URL STATE
	async function loadData() {
		const url = new URL(window.location.href);
		const currentView = url.searchParams.get('view') === 'events' ? 'events' : 'guests';
		const currentStatus = url.searchParams.get('status') || '';
		const currentSearch = url.searchParams.get('search') || '';

		try {
			if (currentView === 'events') {
				const res = await fetch('http://localhost:8080/api/v1/events');
				if (res.ok) {
					events = await res.json();
				}
			} else {
				const guestsUrl = new URL('http://localhost:8080/api/v1/guests');
				if (currentStatus) guestsUrl.searchParams.set('status', currentStatus);
				if (currentSearch) guestsUrl.searchParams.set('search', currentSearch);
				
				const res = await fetch(guestsUrl);
				if (res.ok) {
					guests = await res.json();
				}
			}
		} catch (error) {
			console.error('Failed to load data:', error);
		}
	}

	//INITIAL DATA LOAD AND REACTIVE UPDATES
	$effect(() => {
		// Initial load
		loadData();

		// Set up a more reliable way to detect URL changes
		const handleUrlChange = () => {
			loadData();
		};

		// Listen to popstate (browser back/forward) and custom events
		window.addEventListener('popstate', handleUrlChange);
		
		// Also reload when invalidateAll is called
		const interval = setInterval(() => {
			// This ensures data is reloaded after actions
			loadData();
		}, 100);

		return () => {
			window.removeEventListener('popstate', handleUrlChange);
			clearInterval(interval);
		};
	});

	//URL Sync for Filters and View
	$effect(() => {
		const url = new URL(window.location.href);
		
		if (showEvents) {
			url.searchParams.set('view', 'events');
			// Remove guest filters when switching to events
			url.searchParams.delete('status');
			url.searchParams.delete('search');
			statusFilter = '';
			searchQuery = '';
		} else {
			url.searchParams.delete('view');
			// Set current filters in URL
			if (statusFilter) {
				url.searchParams.set('status', statusFilter);
			} else {
				url.searchParams.delete('status');
			}
			if (searchQuery) {
				url.searchParams.set('search', searchQuery);
			} else {
				url.searchParams.delete('search');
			}
		}
		
		history.replaceState(null, '', url.toString());
		// Force data reload after URL change
		setTimeout(() => loadData(), 50);
	});

	//Initialize state from URL on load
	$effect(() => {
		const url = new URL(window.location.href);
		showEvents = url.searchParams.get('view') === 'events';
		if (!showEvents) {
			statusFilter = url.searchParams.get('status') || '';
			searchQuery = url.searchParams.get('search') || '';
		}
	});

	//Auto-reload when filters change
	$effect(() => {
		if (!showEvents) {
			// Debounce filter changes to avoid too many requests
			const timeoutId = setTimeout(() => {
				const url = new URL(window.location.href);
				if (statusFilter) {
					url.searchParams.set('status', statusFilter);
				} else {
					url.searchParams.delete('status');
				}
				if (searchQuery) {
					url.searchParams.set('search', searchQuery);
				} else {
					url.searchParams.delete('search');
				}
				history.replaceState(null, '', url.toString());
				loadData();
			}, 300);

			return () => clearTimeout(timeoutId);
		}
	});

	//HELPERS
	function resetEventForm() {
		eventTitle = '';
		eventDesc = '';
		eventDate = '';
		eventLocation = '';
		eventError = null;
		showEventForm = false;
	}

	function getBadge(status: string) {
		const map: Record<string, string> = {
			attending: 'bg-green-100 text-green-800',
			declined: 'bg-red-100 text-red-800',
			pending: 'bg-yellow-100 text-yellow-800'
		};
		return map[status] ?? map.pending;
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

	//FORM HANDLERS
	function handleEventFormSubmit() {
		eventSubmitting = true;
		eventError = null;
	}

	function handleEventFormResult(result: ActionResult) {
		eventSubmitting = false;
		if (result.type === 'success') {
			resetEventForm();
			// Reload events data after successful creation
			loadData();
			Swal.fire('Success', 'Event created!', 'success');
		} else if (result.type === 'failure') {
			eventError = result.data?.error ?? 'Failed to create event';
		} else if (result.type === 'error') {
			eventError = result.error?.message ?? 'An error occurred';
		} else {
			eventError = 'Failed to create event';
		}
	}

	// Fixed enhance configuration
	function enhanceFunction({ formData, action, cancel }: { 
		formData: FormData; 
		action: URL; 
		cancel: () => void 
	}) {
		handleEventFormSubmit();
		
		return async ({ result }: { result: ActionResult }) => {
			handleEventFormResult(result);
		};
	}
</script>

<div class="min-h-screen bg-gray-50 py-8">
	<div class="max-w-6xl mx-auto px-4">
		<div class="mb-8">
			<h1 class="text-3xl font-bold text-gray-900 mb-2">Event Guest List Manager</h1>
			<p class="text-gray-600">Manage events and guest RSVPs</p>
		</div>

		<!-- Filters Bar -->
		<div class="bg-white rounded-lg shadow-sm p-4 mb-6 flex flex-wrap gap-4 items-center justify-between">
			<!-- Toggle -->
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

			<!-- Filters + Add Button -->
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
				<form 
					method="POST" 
					action="?/createEvent" 
					use:enhance={enhanceFunction}
					class="grid md:grid-cols-2 gap-4"
				>
					<input name="title" bind:value={eventTitle} placeholder="Event Title *" required class="customInput" />
					<input name="location" bind:value={eventLocation} placeholder="Location" class="customInput" />
					<input name="event_date" bind:value={eventDate} type="datetime-local" required min={new Date().toISOString().slice(0, 16)} class="customInput" />
					<textarea name="description" bind:value={eventDesc} placeholder="Description (optional)" rows="2" class="customInput md:col-span-2"></textarea>
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
				{#if events.length === 0}
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
				{#if guests.length === 0}
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
										<form method="POST" action="?/deleteGuest" use:enhance={({ action }) => {
											return async ({ result }) => {
												await loadData();
											};
										}} class="inline">
											<input type="hidden" name="id" value={g.id} />
											<button type="submit" class="text-red-600 hover:text-red-800">
												Delete
											</button>
										</form>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				{/if}
			</div>

			<!-- Stats -->
			{#if guests.length > 0}
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