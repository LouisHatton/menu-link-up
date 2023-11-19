<script lang="ts">
	import { Button, Spinner } from 'flowbite-svelte';
	import type { ButtonColorType } from 'flowbite-svelte/dist/types';

	export let disabled = false;
	export let loading = false;
	export let btnClass = '';
	export let color: ButtonColorType = 'blue';
</script>

<Button disabled={disabled || loading} {color} on:click class={btnClass}>
	{#if loading}
		<Spinner class="mr-3" size="4" color="white" />
	{/if}
	<slot />
</Button>
