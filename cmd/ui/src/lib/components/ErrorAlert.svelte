<script lang="ts">
	import type { ApiError } from '$lib/services/NetworkService';

	export let title = '';
	export let err: ApiError | undefined = undefined;
</script>

{#if err || $$slots.default}
	<div role="alert" class="rounded border-s-4 border-red-500 bg-red-50 p-4">
		{#if title != ''}
			<strong class="block font-medium text-red-800 mb-2"> {title} </strong>
		{/if}

		<p class="text-sm text-red-700">
			{#if err}
				{err?.message} ({err?.status})
			{/if}
			<slot />
		</p>
	</div>
{/if}
