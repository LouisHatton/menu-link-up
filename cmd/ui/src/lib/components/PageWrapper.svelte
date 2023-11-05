<script lang="ts">
	import { classnames } from '$lib/util';

	export let sm = false;
</script>

{#if sm}
	<div class={classnames($$props.class, 'px-6 py-4 w-full')}><slot /></div>
{:else}
	<div class={classnames($$props.class, 'max-w-6xl mx-auto px-6 py-10 w-full')}><slot /></div>
{/if}
