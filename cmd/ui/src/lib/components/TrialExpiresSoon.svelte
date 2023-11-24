<script>
	import { authStore } from '$lib/stores/authStore';

	import { timeFromNow } from '$lib/util';
	import LoadingButton from './LoadingButton.svelte';
</script>

{#if $authStore.dbUser?.subscriptionStatus == 'trialing'}
	<div
		class="border-purple-500 border-2 my-8 py-4 px-24 gap-y-4 flex flex-col justify-center text-center items-center rounded-2xl"
	>
		<p>
			You are currently on a free trial which expires
			<strong>
				{#if $authStore.dbUser.trialEnd}
					{timeFromNow($authStore.dbUser.trialEnd)}
				{:else}
					soon
				{/if}
				🔒</strong
			>
		</p>
		<p>
			If you do not have a valid card when the trial ends you will loose access to MenuLink-Up
			<br />and any uploaded files will be removed.
		</p>
		<LoadingButton color="purple">Update Card Information</LoadingButton>
	</div>
{/if}
