<script lang="ts">
	import type { ApiError } from '$lib/services/NetworkService';
	import UserService from '$lib/services/UserService';
	import { authStore } from '$lib/stores/authStore';
	import Navbar from './Navbar.svelte';

	let getDbUserError: string | undefined = undefined;

	$: getDbUser();

	async function getDbUser() {
		getDbUserError = undefined;
		try {
			let dbUser = await UserService.getUser();
			$authStore.dbUser = dbUser;
		} catch (err: unknown) {
			getDbUserError = (err as ApiError).message;
		}
	}
</script>

<section
	class="min-h-screen max-w-full flex flex-col w-full bg-zinc-50 text-zinc-800 dark:bg-gray-700 dark:text-zinc-50"
>
	<Navbar />

	<div class="w-full relative">
		{#if getDbUserError}
			<div class="bg-red-900 text-white flex flex-row py-3 px-8">
				<p>Encountered an error fetching user profile: {getDbUserError}</p>
			</div>
		{/if}
		<slot />
	</div>
</section>
