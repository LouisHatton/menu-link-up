<script lang="ts">
	import { goto } from '$app/navigation';
	import FileService from '$lib/services/FileService';
	import type { PageData } from './$types';

	export let data: PageData;

	let error = false;

	$: getLink();

	function getLink() {
		setTimeout(async () => {
			try {
				let link = await FileService.getFileLink(data.id);
				goto(link);
			} catch {
				error = true;
			}
		}, 2000);
	}
</script>

{#if !error}
	<p>Loading '{data.id}'...</p>
{:else}
	<p>There was an error loading your file</p>
{/if}
