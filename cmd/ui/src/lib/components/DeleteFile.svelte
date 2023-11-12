<script lang="ts">
	import Cross from '$lib/icons/Cross.svelte';
	import type { File } from '$lib/services/FileService';
	import FileService from '$lib/services/FileService';
	import type { ApiError } from '$lib/services/NetworkService';
	import { Button, Modal } from 'flowbite-svelte';
	import { createEventDispatcher } from 'svelte';

	export let file: File;
	const dispatch = createEventDispatcher();

	let popupModal = false;

	async function handleDelete() {
		try {
			await FileService.deleteFile(file.id);
			dispatch('delete');
			popupModal = false;
		} catch (err: unknown) {
			alert((err as ApiError).message);
		}
	}
</script>

<Button color="red" on:click={() => (popupModal = true)}>Delete</Button>

<Modal bind:open={popupModal} size="xs">
	<div class="text-center text-black">
		<Cross class="mx-auto mb-4 fill-current  w-12 h-12 " />
		<h3 class="mb-5 text-lg font-normal">
			Are you sure you want to delete the file "{file.name}"? <br /> This action can not be undone
		</h3>
		<Button color="red" class="mr-2" on:click={handleDelete}>Yes, I'm sure</Button>
		<Button color="alternative" on:click={() => (popupModal = false)}>No, cancel</Button>
	</div>
</Modal>
