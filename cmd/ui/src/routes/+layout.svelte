<script lang="ts">
	import { goto } from '$app/navigation';
	import { navigating, page } from '$app/stores';
	import AuthenticationService from '$lib/services/AuthenticationService';
	import { authStore } from '$lib/stores/authStore';
	import { darkMode } from '$lib/stores/darkMode';
	import { Spinner } from 'flowbite-svelte';
	import '../app.css';
	import Navigating from './Navigating.svelte';
	import LoadingScreen from '$lib/components/LoadingScreen.svelte';

	let loggingIn = true;
	$: $page.route && checkIfLoggedIn($authStore.initialised);
	let loading = true;

	async function checkIfLoggedIn(init: boolean) {
		loading = true;
		loggingIn = localStorage.getItem('loggingIn') != null;
		if (!init) return;
		try {
			const isInAuthRoute = $page.route.id?.includes('(authentication)');
			const token = await AuthenticationService.getToken();

			globalThis.jwt = token; // Set console JWT for access in dev tools

			if (!token && !isInAuthRoute) {
				await goto('/login');
			} else if (token && isInAuthRoute) {
				await goto('/');
				localStorage.removeItem('loggingIn');
			}
		} finally {
			loading = false;
		}
	}
</script>

<main class={`${$darkMode ? 'dark' : ''}`}>
	<Navigating />
	{#if $authStore.isLoggedIn || ($page.route.id?.includes('(authentication)') && !loggingIn && !$navigating)}
		<slot />
	{:else if loading || $navigating}
		<LoadingScreen />
	{/if}
</main>
