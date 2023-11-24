<script lang="ts">
	import Card from '$lib/components/Card.svelte';
	import type { User } from 'firebase/auth';

	import PageWrapper from '$lib/components/PageWrapper.svelte';
	import { authStore } from '$lib/stores/authStore';
	import { Button, ChevronLeft, Input, Label } from 'flowbite-svelte';
	import UserService from '$lib/services/UserService';
	import LoadingButton from '$lib/components/LoadingButton.svelte';
	import AuthenticationService from '$lib/services/AuthenticationService';
	import Arrow from '$lib/icons/Arrow.svelte';
	import type { ApiError } from '$lib/services/NetworkService';
	import TrialExpiresSoon from '$lib/components/TrialExpiresSoon.svelte';

	let emailAddress = '';
	let name = '';
	let saving = false;
	let deleting = false;
	let loadingData = true;

	$: $authStore.user ? updateVars($authStore.user) : (loadingData = true);

	function updateVars(user: User) {
		loadingData = false;
		if (user.email) {
			emailAddress = user.email;
		}
		if (user.displayName) {
			name = user.displayName;
		}
	}

	async function handleClick() {
		saving = true;
		await UserService.updateUserDisplayName(name);
		saving = false;
	}

	async function handleResetPassword() {
		console.log(await AuthenticationService.getToken());
	}

	async function handleDeleteAccount() {
		deleting = true;
		try {
			await UserService.deleteUser();
			await AuthenticationService.logOut();
		} catch (err: unknown) {
			alert(
				`There was an issue, if this continues please contact support (${
					(err as ApiError).message
				} - ${(err as ApiError).status})`
			);
		} finally {
			deleting = false;
		}
	}
</script>

<PageWrapper>
	<a href="/" class="text-blue-600 hover:underline flex flex-row gap-x-2 mb-2">
		<Arrow class="w-4 rotate-180 fill-current" />
		Go back to the Dashboard
	</a>
	<h2 class="text-4xl font-semibold">Settings</h2>
	<div class="flex flex-col gap-6 mt-14">
		<Card class="col-span-2">
			<h3 class="text-2xl font-semibold">Your Current Plan</h3>
			<TrialExpiresSoon />
			<div class="mt-6 flex flex-row justify-between items-center">
				<div>
					<p class="text-lg">Basic Plan (£5/month)</p>
					<p class="">Renews on: 4th December</p>
				</div>
				<Button color="dark">Change Plan</Button>
			</div>
			<div class="my-8 border-t border-gray-200" />
			<h3 class="text-2xl font-semibold">Billing Information</h3>
			<div class="mt-6 flex flex-row justify-between items-center">
				<div>
					<p class="text-lg">No card on file</p>
					<p class="">Unknown Expiry</p>
				</div>
				<Button color="dark">Update Billing Information</Button>
			</div>
		</Card>
		<div>
			<Card>
				<h4 class="text-lg font-semibold mb-4">Danger Zone</h4>
				<LoadingButton color="red" on:click={handleDeleteAccount} loading={deleting}
					>Delete Account</LoadingButton
				>
			</Card>
		</div>
	</div>
</PageWrapper>
