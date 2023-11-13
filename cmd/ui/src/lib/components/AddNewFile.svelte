<script lang="ts">
	import type { NewFile } from '$lib/services/FileService';
	import FileService from '$lib/services/FileService';
	import type { ApiError } from '$lib/services/NetworkService';
	import { sanitiseSlug } from '$lib/util';
	import { Button, Fileupload, Input, Modal, Select } from 'flowbite-svelte';
	import { createEventDispatcher } from 'svelte';
	import LoadingButton from './LoadingButton.svelte';
	import UrlCheck from './UrlCheck.svelte';
	import Alert from './Alert.svelte';

	export let large = false;
	export let disabled = false;

	const dispatch = createEventDispatcher();
	let openModal = false;
	let uploadDisabled = true;
	let loading = false;
	let checkingUrl = false;
	let urlIsAvailable = true;

	let slug = '';
	let filename = '';
	let slugChanged = false;
	let uploadedFiles: FileList;

	let apiError: ApiError | undefined;

	$: slug = sanitiseSlug(slug);
	$: if (!slugChanged) {
		slug = sanitiseSlug(filename);
	}

	$: if (slug != '' && filename != '' && uploadedFiles?.length && urlIsAvailable) {
		uploadDisabled = false;
	} else {
		uploadDisabled = true;
	}

	function handleModalClick() {
		if (disabled) return;
		openModal = true;
	}

	async function createNewFile() {
		loading = true;
		apiError = undefined;
		let newFile: NewFile = {
			name: filename,
			slug,
			fileSize: uploadedFiles[0].size
		};

		if (uploadedFiles == undefined || uploadedFiles?.length < 1) {
			apiError = {
				status: 400,
				message: 'No file has been selected'
			};
			loading = false;
			return;
		}

		try {
			let url = await FileService.createFile(newFile);
			let resp = await FileService.uploadFile(url.url, uploadedFiles[0], uploadedFiles[0].size);
			if (resp.ok) {
				dispatch('create');
			} else {
				apiError = {
					status: resp.status,
					message: 'Unable to upload file. If issue persists, please contact support.'
				};
			}
		} catch (err: unknown) {
			apiError = err as ApiError;
		} finally {
			loading = false;
		}
	}
</script>

{#if large && !disabled}
	<button
		on:click={handleModalClick}
		class="w-full py-28 border-2 border-dashed border-gray-400 rounded-md hover:bg-gray-100"
		>Click here to add your first Menu!
	</button>
{:else}
	<Button {disabled} on:click={handleModalClick}>Add New</Button>
{/if}
<Modal title="Add New" bind:open={openModal} class="lg:w-[80%]">
	<p class="text-base leading-relaxed text-gray-500 dark:text-gray-400">
		Enter a name and a unique URL for your pdf to be hosted from
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
		Upload file <span class="text-sm">(PDF only)</span>
		<Fileupload accept="application/pdf" bind:files={uploadedFiles} />
	</p>
	<UrlCheck available={urlIsAvailable} loading={checkingUrl} />
	{#if !checkingUrl && urlIsAvailable}
		<p class="text-base text-black dark:text-gray-400">
			The file will be uploaded to: <br />
			<span class="break-before-avoid text-blue-600">https://menulinkup.com/f/{slug}</span>
		</p>
	{/if}
	{#if apiError}
		<p class="text-red-600">Error: {apiError.message} ({apiError.status})</p>
	{/if}
	<svelte:fragment slot="footer">
		<LoadingButton {loading} on:click={createNewFile} disabled={uploadDisabled}
			>Upload</LoadingButton
		>
	</svelte:fragment>
</Modal>
