<script lang="ts">
	import AuthenticationService from '$lib/services/AuthenticationService';
	import { Button, Label, Input, Spinner } from 'flowbite-svelte';
	import SignInWithGoogle from '../SignInWithGoogle.svelte';
	import ErrorAlert from '$lib/components/ErrorAlert.svelte';

	let email: string;
	let password: string;
	let loading = false;
	let errorMessage = '';

	async function handleLogin(e: MouseEvent) {
		e.preventDefault();
		loading = true;
		errorMessage = '';
		try {
			await AuthenticationService.signInWithPassword(email, password);
		} catch (err: unknown) {
			errorMessage =
				'There was an issue with your login, please check your email and password and try again';
		} finally {
			loading = false;
		}
	}
</script>

<div class="min-h-screen flex justify-center items-center">
	<div class="p-6 w-4/5 lg:w-2/5 space-y-4 md:space-y-6 sm:p-8">
		<div class="flex flex-row justify-center items-center gap-x-2 my-8">
			<!-- <img src="/static/insight-wave-800x800.png" alt="logo" class="w-10" /> -->
			<h1 class="text-2xl text-blue-600 font-semibold">MenuLink-Up</h1>
		</div>
		<form class="flex flex-col space-y-6" action="/">
			<h3 class="text-xl font-medium text-gray-900 dark:text-white p-0">Login</h3>
			{#if errorMessage != ''}
				<ErrorAlert>{errorMessage}</ErrorAlert>
			{/if}
			<Label class="space-y-2">
				<span>Your email</span>
				<Input
					bind:value={email}
					type="email"
					name="email"
					placeholder="name@company.com"
					required
				/>
			</Label>
			<Label class="space-y-2">
				<span>Your password</span>
				<Input bind:value={password} type="password" name="password" placeholder="•••••" required />
			</Label>
			<div class="flex items-start">
				<a href="/register" class="ml-auto text-sm text-blue-500 hover:underline dark:text-blue-400"
					>Forgot password?</a
				>
			</div>
			<Button disabled={loading} on:click={handleLogin} type="submit" class="w-full" color="blue"
				>{#if loading}<Spinner class="mr-3" size="4" color="white" />{/if}Sign in</Button
			>
			<SignInWithGoogle />
			<p class="text-sm font-light text-gray-500 dark:text-gray-400">
				Don’t have an account yet? <a
					href="/register"
					class="font-medium text-primary-600 hover:underline dark:text-primary-500">Sign up</a
				>
			</p>
		</form>
	</div>
</div>
