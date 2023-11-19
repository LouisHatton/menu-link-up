<script lang="ts">
	import type { NewFile } from '$lib/services/FileService';
	import FileService from '$lib/services/FileService';
	import type { ApiError } from '$lib/services/NetworkService';
	import { sanitiseSlug } from '$lib/util';
	import { Button, Fileupload, Input, Modal, Select } from 'flowbite-svelte';
	import { createEventDispatcher } from 'svelte';
	import LoadingButton from './LoadingButton.svelte';
	import UrlCheck from './UrlCheck.svelte';
	import Alert from './ErrorAlert.svelte';
	import ErrorAlert from './ErrorAlert.svelte';

	export let large = false;
	export let disabled = false;

	const dispatch = createEventDispatcher();
	let openModal = false;
	let uploadDisabled = true;
	let loading = false;
	let checkingUrl = false;
	let urlIsAvailable = false;

	let slug = '';
	let filename = '';
	let slugChanged = false;
	let uploadedFiles: FileList;

	let apiError: ApiError | undefined;

	$: uploadedFile = uploadedFiles?.length > 0 ? uploadedFiles[0] : undefined;
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

		if (!uploadedFile) {
			apiError = {
				status: 400,
				message: 'No file has been selected'
			};
			loading = false;
			return;
		}

		let newFile: NewFile = {
			name: filename,
			slug,
			fileSize: uploadedFiles[0].size
		};

		try {
			let url = await FileService.createFile(newFile);
			let resp = await FileService.uploadFile(url.url, uploadedFile, uploadedFile.size);
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
		class="w-full py-28 text-2xl font-medium text-zinc-700 border-2 bg-stone-100 border-gray-300 rounded-3xl hover:bg-blue-100 hover:border-blue-600"
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
	{#if openModal}
		<p>
			Upload file <span class="text-sm">(PDF only)</span>
			<Fileupload accept="application/pdf" bind:files={uploadedFiles} />
		</p>
	{/if}
	<UrlCheck bind:available={urlIsAvailable} {slug} />
	{#if !checkingUrl && urlIsAvailable && uploadedFile}
		<p class="text-base text-black dark:text-gray-400">
			The file will be uploaded to: <br />
			<span class="break-before-avoid text-blue-600">https://menulinkup.com/f/{slug}</span>
		</p>
	{/if}
	{#if apiError}
		<ErrorAlert err={apiError} />
	{/if}
	<LoadingButton {loading} on:click={createNewFile} disabled={uploadDisabled}>Upload</LoadingButton>
</Modal>
