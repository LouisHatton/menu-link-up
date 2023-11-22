<script lang="ts">
	import { Input, Modal } from 'flowbite-svelte';
	import { createEventDispatcher } from 'svelte';
	import LoadingButton from './LoadingButton.svelte';
	import { authStore } from '$lib/stores/authStore';
	import AuthenticationService from '$lib/services/AuthenticationService';
	import UndrawSubscriber from '$lib/illustrations/undraw_subscriber.svelte';
	import type { ApiError } from '$lib/services/NetworkService';
	import ErrorAlert from './ErrorAlert.svelte';

	const dispatch = createEventDispatcher();
	let openModal = true;
	let loading = false;
	let sendingEmail = false;
	let apiErr: ApiError | undefined;

	$: if ($authStore.user?.emailVerified) {
		openModal = false;
	} else {
		openModal = true;
		sendAuthEmail();
	}

	async function sendAuthEmail(force: boolean = false) {
		sendingEmail = true;
		try {
			await AuthenticationService.sendVerifyEmail(force);
		} catch (err: unknown) {
			alert(JSON.stringify(err));
		} finally {
			sendingEmail = false;
		}
	}

	async function checkEmailVerified() {
		loading = true;
		apiErr = undefined;
		try {
			await AuthenticationService.reloadUser();
		} catch (err: unknown) {
			alert(JSON.stringify(err));
		} finally {
			loading = false;
		}

		if (!$authStore.user?.emailVerified) {
			apiErr = {
				message: 'Email not verified, please try again',
				status: 400
			};
		} else {
			dispatch('verified');
		}
	}
</script>

<Modal title="Verify Email" bind:open={openModal} class="lg:w-[80%] ">
	<UndrawSubscriber class="w-full h-60" />
	<ErrorAlert err={apiErr} />
	<div class="flex flex-col justify-center gap-y-4 text-center">
		<p class="text-base leading-relaxed text-gray-500 dark:text-gray-400">
			We have sent you an email with a link to verify your email address.
		</p>
		<p class="text-base leading-relaxed text-gray-500 dark:text-gray-400">
			After verifying your email, please click continue.
		</p>
		<div class="flex flex-row place-content-center gap-x-4">
			<LoadingButton
				color="light"
				loading={sendingEmail}
				on:click={() => {
					sendAuthEmail(true);
				}}>Resend Email</LoadingButton
			>
			<LoadingButton
				{loading}
				on:click={() => {
					checkEmailVerified();
				}}>Continue</LoadingButton
			>
		</div>
	</div>
</Modal>
