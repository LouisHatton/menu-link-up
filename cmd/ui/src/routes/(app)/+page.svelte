<script lang="ts">
	import AddNewFile from '$lib/components/AddNewFile.svelte';
	import Banner from '$lib/components/Banner.svelte';
	import Card from '$lib/components/Card.svelte';
	import DeleteFile from '$lib/components/DeleteFile.svelte';
	import ErrorAlert from '$lib/components/ErrorAlert.svelte';
	import PageLoader from '$lib/components/PageLoader.svelte';
	import PageWrapper from '$lib/components/PageWrapper.svelte';
	import CheckMark from '$lib/icons/CheckMark.svelte';
	import type { File } from '$lib/services/FileService';
	import FileService from '$lib/services/FileService';
	import type { ApiError } from '$lib/services/NetworkService';
	import { authStore } from '$lib/stores/authStore';
	import {
		Button,
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell
	} from 'flowbite-svelte';
	import dayjs from 'dayjs';
	import { timeFromNow } from '$lib/util';

	let loading = false;
	let filesError: ApiError | undefined = undefined;
	let files: File[] = [];

	$: getProjects();
	async function getProjects() {
		loading = true;
		filesError = undefined;
		try {
			files = await FileService.listFiles();
		} catch (err: unknown) {
			filesError = err as ApiError;
		} finally {
			loading = false;
		}
	}
</script>

{#if $authStore.dbUser?.subscriptionStatus === 'trialing'}
	<Banner href="/settings">
		<p>
			Your trial expires
			<strong>
				{#if $authStore.dbUser.trialEnd}
					{timeFromNow($authStore.dbUser.trialEnd)}
				{:else}
					soon
				{/if}
				🔒</strong
			>
		</p>
		<p>
			To continue using MenuLink-Up, please add your billing info <strong>here</strong> before your trial
			finishes.
		</p>
	</Banner>
{/if}
<PageWrapper>
	<h2 class="text-4xl font-semibold text-center">Your Dashboard</h2>
	<div class="mt-10">
		<Card>
			<h3 class="text-3xl font-semibold mb-6">
				Menus <span class="text-sm font-medium text-gray-500">({files.length})</span>
			</h3>
			<PageLoader {loading}>
				{#if filesError}
					<ErrorAlert title="Something went wrong fetching your files" err={filesError} />
				{:else}
					<Table>
						<TableHead>
							<TableHeadCell>Name</TableHeadCell>
							<TableHeadCell>Url</TableHeadCell>
							<TableHeadCell>Status</TableHeadCell>
							<TableHeadCell>Edit</TableHeadCell>
							<TableHeadCell />
						</TableHead>
						<TableBody>
							{#each files as file}
								<TableBodyRow>
									<TableBodyCell>{file.name}</TableBodyCell>
									<TableBodyCell
										><a
											class="text-blue-600 hover:underline"
											target="_blank"
											href={location.origin + '/f/' + file.slug}>{location.origin}/f/{file.slug}</a
										></TableBodyCell
									>
									<TableBodyCell
										><CheckMark class="fill-current text-green-500 w-7" /></TableBodyCell
									>
									<TableBodyCell>
										<div class="flex flex-row gap-x-4">
											<Button>Edit</Button>
											<DeleteFile {file} on:delete={getProjects} />
										</div>
									</TableBodyCell>
								</TableBodyRow>
							{/each}
						</TableBody>
					</Table>
					<div class="mt-4">
						<AddNewFile large={files.length < 1} on:create={getProjects} disabled={loading} />
					</div>
				{/if}
			</PageLoader>
		</Card>
	</div>
</PageWrapper>
