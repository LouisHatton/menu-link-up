<script lang="ts">
	import AddNewFile from '$lib/components/AddNewFile.svelte';
	import Card from '$lib/components/Card.svelte';
	import DeleteFile from '$lib/components/DeleteFile.svelte';
	import PageLoader from '$lib/components/PageLoader.svelte';
	import PageWrapper from '$lib/components/PageWrapper.svelte';
	import CheckMark from '$lib/icons/CheckMark.svelte';
	import type { File } from '$lib/services/FileService';
	import FileService from '$lib/services/FileService';
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

	let loading = false;
	let files: File[] = [];

	$: getProjects();
	async function getProjects() {
		loading = true;
		try {
			files = await FileService.listFiles();
		} finally {
			loading = false;
		}
	}
</script>

<PageWrapper>
	<h2 class="text-4xl font-semibold">Welcome back {$authStore.user?.displayName ?? ''}!</h2>
	<div class="mt-14">
		<Card>
			<h3 class="text-2xl font-semibold mb-4">
				Your Menus: <span class="text-sm text-gray-500">(1/3)</span>
			</h3>
			<PageLoader {loading}>
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
										href={'https://menulinkup.com/' + file.slug}
										>https://menulinkup.com/{file.slug}</a
									></TableBodyCell
								>
								<TableBodyCell><CheckMark class="fill-current text-green-500 w-7" /></TableBodyCell>
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
					<AddNewFile large={files.length < 1} on:create={getProjects} />
				</div>
				<!-- <pre>{JSON.stringify(files, undefined, 2)}</pre> -->
			</PageLoader>
		</Card>
	</div>
</PageWrapper>
