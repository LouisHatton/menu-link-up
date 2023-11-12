<script lang="ts">
	import type { NewFile } from '$lib/services/FileService';
	import FileService from '$lib/services/FileService';
	import type { ApiError } from '$lib/services/NetworkService';
	import { sanitiseSlug } from '$lib/util';
	import { Button, Fileupload, Input, Modal, Select } from 'flowbite-svelte';
	import { createEventDispatcher } from 'svelte';

	export let large = false;

	const dispatch = createEventDispatcher();
	let openModal = false;

	let slug = '';
	let filename = '';
	let slugChanged = false;

	let uploadedFiles: FileList;

	let fileuploadprops = {
		id: 'menu'
	};

	$: slug = sanitiseSlug(slug);
	$: if (!slugChanged) {
		slug = sanitiseSlug(filename);
	}
	// $: checkSlug(slug);

	function handleModalClick() {
		openModal = true;
	}

	async function createNewFile() {
		let newFile: NewFile = {
			name: filename,
			slug
		};

		if (uploadedFiles.length < 1) {
			alert('no file selected');
		}

		try {
			let url = await FileService.createFile(newFile);
			if (url.url == '') {
				alert('no url returned');
			} else {
				await FileService.uploadFile(url.url, uploadedFiles[0]);
				dispatch('create');
			}
		} catch (err: unknown) {
			alert((err as ApiError).message);
		}
	}
</script>

{#if large}
	<button
		on:click={handleModalClick}
		class="w-full py-28 border-2 border-dashed border-gray-400 rounded-md hover:bg-gray-100"
		>Click here to add your first Menu!
	</button>
{:else}
	<Button on:click={handleModalClick}>Add New</Button>
{/if}
<Modal title="Add New" bind:open={openModal} class="lg:w-[80%]">
	<p class="text-base leading-relaxed text-gray-500 dark:text-gray-400">
		To add a new file, give a unique filename and upload your pdf.
	</p>
	<p>
		Name:
		<Input bind:value={filename} />
	</p>
	<p>
		URL: <span class="text-sm">(Must only contain lowercase letters and '-')</span>
		<Input
			bind:value={slug}
			on:change={() => {
				slugChanged = true;
			}}
		/>
	</p>
	<p>
		Upload file
		<Fileupload {...fileuploadprops} accept="application/pdf" bind:files={uploadedFiles} />
	</p>
	<p class="text-base text-black dark:text-gray-400">
		The file will be uploaded to: <br />
		<span class="break-before-avoid text-blue-600">https://menulinkup.com/{slug}</span>
	</p>
	<svelte:fragment slot="footer">
		<Button on:click={createNewFile}>Upload</Button>
	</svelte:fragment>
</Modal>
