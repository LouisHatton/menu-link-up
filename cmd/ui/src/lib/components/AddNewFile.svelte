<script lang="ts">
	import type { Project } from '$lib/services/ProjectService';
	import { sanitiseSlug } from '$lib/util';
	import { Button, Fileupload, Input, Modal, Select } from 'flowbite-svelte';

	export let large = false;
	export let project: Project;

	let openModal = false;

	let filename = '';

	let fileuploadprops = {
		id: 'user_avatar'
	};

	$: filename = sanitiseSlug(filename);

	function handleModalClick() {
		openModal = true;
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
<Modal title="Add New" bind:open={openModal} class="lg:w-[80%]" autoclose>
	<p class="text-base leading-relaxed text-gray-500 dark:text-gray-400">
		To add a new file, give a unique filename and upload your pdf.
	</p>
	<!-- <p>
		Project:
		<Select bind:value={projectName} items={[{ name: 'pizza', value: 'pizza' }]} />
	</p> -->
	<p>
		File Name: <span class="text-sm">(Must only contain lowercase letters and '-')</span>
		<Input bind:value={filename} />
	</p>
	<p>
		Upload file
		<Fileupload {...fileuploadprops} accept="application/pdf" />
	</p>
	<p class="text-base text-black dark:text-gray-400">
		The file will be uploaded to: <br />
		<span class="break-before-avoid text-blue-600"
			>https://menulinkup.com/{project.name}/{filename}</span
		>
	</p>
	<svelte:fragment slot="footer">
		<Button on:click={() => alert('Handle "success"')}>Upload</Button>
	</svelte:fragment>
</Modal>
