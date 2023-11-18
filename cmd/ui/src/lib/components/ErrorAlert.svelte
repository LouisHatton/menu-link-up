<script lang="ts">
	import type { ApiError } from '$lib/services/NetworkService';

	export let title = 'Something went wrong';
	export let err: ApiError | undefined = undefined;
</script>

<div role="alert" class="rounded border-s-4 border-red-500 bg-red-50 p-4">
	<strong class="block font-medium text-red-800"> {title} </strong>

	<p class="mt-2 text-sm text-red-700">
		{#if err}
			{err?.message} ({err?.status})
		{/if}
		<slot />
	</p>
</div>
