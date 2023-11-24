<script lang="ts">
	import PageLoader from '$lib/components/PageLoader.svelte';
	import VerifyEmail from '$lib/components/VerifyEmail.svelte';
	import type { ApiError } from '$lib/services/NetworkService';
	import UserService, { type DbUser } from '$lib/services/UserService';
	import { authStore } from '$lib/stores/authStore';
	import Navbar from './Navbar.svelte';

	let getDbUserError: string | undefined = undefined;
	let dbUser: DbUser | undefined = undefined;
	let loading = false;

	$: getDbUser();

	async function getDbUser() {
		if (!$authStore.user?.emailVerified) return;

		dbUser = undefined;
		getDbUserError = undefined;
		loading = true;
		try {
			dbUser = await UserService.getUser();
			$authStore.dbUser = dbUser;
		} catch (err: unknown) {
			getDbUserError = (err as ApiError).message;
		} finally {
			loading = false;
		}
	}
</script>

<section
	class="min-h-screen max-w-full flex flex-col w-full bg-zinc-50 text-zinc-800 dark:bg-gray-700 dark:text-zinc-50"
>
	<Navbar />
	{#if $authStore.user?.emailVerified}
		<PageLoader {loading} size="2xl">
			<div class="w-full relative">
				{#if getDbUserError}
					<div class="bg-red-900 text-white flex flex-row py-3 px-8">
						<p>Encountered an error fetching user profile: {getDbUserError}</p>
					</div>
				{/if}
				<slot />
			</div>
		</PageLoader>
	{/if}
	<VerifyEmail on:verified={getDbUser} />
</section>
