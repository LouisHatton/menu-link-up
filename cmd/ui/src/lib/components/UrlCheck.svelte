<script>
	import CheckMark from '$lib/icons/CheckMark.svelte';
	import Cross from '$lib/icons/Cross.svelte';
	import { classnames } from '$lib/util';
	import { Spinner } from 'flowbite-svelte';

	export let loading = false;
	export let available = true;
</script>

<div
	class={classnames(
		'flex flex-row gap-x-2 text-sm',
		loading ? 'text-gray-500' : available ? 'text-green-500' : 'text-red-500'
	)}
>
	{#if loading}
		<Spinner size="5" color="gray" />
		Checking...
	{:else}
		{#if available}
			<CheckMark class="fill-current w-5" />
		{:else}
			<Cross class="fill-current w-5" />
		{/if}
		<p>URL is {available ? 'available!' : 'not available.'}</p>
	{/if}
</div>
