<script lang="ts">
	import CheckMark from '$lib/icons/CheckMark.svelte';
	import Cross from '$lib/icons/Cross.svelte';
	import FileService, { type NewFile } from '$lib/services/FileService';
	import { classnames } from '$lib/util';
	import { Spinner } from 'flowbite-svelte';

	export let slug: string;
	export let available = false;

	let loading = false;
	let timeout: NodeJS.Timeout;

	$: checkSlug(slug);

	async function checkSlug(slug: string) {
		clearTimeout(timeout);
		if (slug == '') {
			available = false;
			return;
		}

		loading = true;
		timeout = setTimeout(async () => {
			available = await FileService.checkFileSlug(slug);
			loading = false;
		}, 2500);
	}
</script>

<div
	class={classnames(
		'flex flex-row gap-x-2 items-center',
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
