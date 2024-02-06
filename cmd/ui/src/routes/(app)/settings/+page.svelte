<script lang="ts">
	import Card from '$lib/components/Card.svelte';
	import type { User } from 'firebase/auth';

	import PageWrapper from '$lib/components/PageWrapper.svelte';
	import { authStore } from '$lib/stores/authStore';
	import { Button, ChevronLeft, Input, Label } from 'flowbite-svelte';
	import UserService, { type Billing } from '$lib/services/UserService';
	import LoadingButton from '$lib/components/LoadingButton.svelte';
	import AuthenticationService from '$lib/services/AuthenticationService';
	import Arrow from '$lib/icons/Arrow.svelte';
	import type { ApiError } from '$lib/services/NetworkService';
	import TrialExpiresSoon from '$lib/components/TrialExpiresSoon.svelte';
	import PageLoader from '$lib/components/PageLoader.svelte';
	import dayjs from 'dayjs';
	import { dateNumMonthYear } from '$lib/util';

	let deleting = false;
	let loading = true;
	let billingInfo: Billing | undefined;

	$: getBillingInfo();

	async function getBillingInfo() {
		loading = true;
		try {
			billingInfo = await UserService.getUserBilling();
		} catch (err: unknown) {
		} finally {
			loading = false;
		}
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
	<PageLoader {loading} size="2xl">
		<div class="flex flex-col gap-6 mt-14">
			<Card class="col-span-2">
				<h3 class="text-2xl font-semibold">Your Current Plan</h3>
				<TrialExpiresSoon />
				<div class="mt-6 flex flex-row justify-between items-center">
					<div>
						<p class="text-lg">
							{billingInfo?.planName} (£{(billingInfo?.price ?? 0) /
								100}/{billingInfo?.billingInterval})
						</p>
						<p class="">Renews: {dateNumMonthYear(billingInfo?.currentPeriodEnd ?? '')}</p>
					</div>
					<Button color="dark">Change Plan</Button>
				</div>
				<div class="my-8 border-t border-gray-200" />
				<h3 class="text-2xl font-semibold">Billing Information</h3>
				<div class="mt-6 flex flex-row justify-between items-center">
					<div>
						<p class="text-lg">
							{#if billingInfo?.defaultPayment?.brand && billingInfo?.defaultPayment.lastFour}{billingInfo
									.defaultPayment.brand} •••• {billingInfo.defaultPayment.lastFour}{:else}No card on
								file{/if}
						</p>
						<p class="">
							{#if billingInfo?.defaultPayment?.expiresMonth && billingInfo?.defaultPayment.expiresYear}Expires
								{billingInfo?.defaultPayment.expiresMonth}/{billingInfo?.defaultPayment
									.expiresYear}{:else}Unknown Expiry{/if}
						</p>
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
	</PageLoader>
</PageWrapper>
