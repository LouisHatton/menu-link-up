<script lang="ts">
	import { goto } from '$app/navigation';
	import Card from '$lib/components/Card.svelte';
	import LoadingButton from '$lib/components/LoadingButton.svelte';
	import PageWrapper from '$lib/components/PageWrapper.svelte';
	import UrlCheck from '$lib/components/UrlCheck.svelte';
	import CheckMark from '$lib/icons/CheckMark.svelte';
	import type { ApiError } from '$lib/services/NetworkService';
	import ProjectService from '$lib/services/FileService';
	import { sanitiseSlug } from '$lib/util';
	import { Button, Input } from 'flowbite-svelte';

	let error = '';
	let projectName = '';
	let slug = '';
	let slugChanged = false;

	let checkingSlugTimeout: NodeJS.Timeout;
	let slugOk = true;
	let checkingSlug = false;
	let creating = false;

	$: if (!slugChanged) {
		slug = sanitiseSlug(projectName);
	}
	$: slug = sanitiseSlug(slug);
	$: checkSlug(slug);

	async function checkSlug(slug: string) {
		clearTimeout(checkingSlugTimeout);
		if (slug === '') {
			slugOk = false;
			checkingSlug = false;
			return;
		}
		checkingSlug = true;

		checkingSlugTimeout = setTimeout(async () => {
			console.log('checking slug %s', slug);

			try {
				slugOk = await ProjectService.checkFileSlug({
					name: projectName,
					slug: slug
				});
			} finally {
				checkingSlug = false;
			}
		}, 3000);
	}

	async function createProject() {
		creating = true;
		error = '';
		try {
			await ProjectService.createFile({
				name: projectName,
				slug
			});
		} catch (e: unknown) {
			error = (e as ApiError).message;
			return;
		} finally {
			creating = false;
		}

		goto('/');
	}
</script>

<PageWrapper>
	<div class="mt-8">
		<Card>
			<h3 class="text-2xl font-semibold mb-4">Lets finish setting up your account:</h3>
			<p class="my-6 text-lg">
				Ready to showcase your restaurant's flavors? Just enter your business name and URL, and
				start uploading your menus hassle-free.
			</p>
			<div class="flex flex-col gap-y-4 w-[50%]">
				<p>
					Business Name:
					<Input bind:value={projectName} />
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
				<UrlCheck loading={checkingSlug} available={slugOk} />
				{#if error != ''}
					<p class="text-red-500">{error}</p>
				{/if}
			</div>
			<div class="my-6">
				<LoadingButton color="blue" on:click={createProject} disabled={!slugOk} loading={creating}
					>Lets Go!</LoadingButton
				>
			</div>
		</Card>
	</div>
</PageWrapper>
