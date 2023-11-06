<script lang="ts">
	import { goto } from '$app/navigation';
	import AddNewFile from '$lib/components/AddNewFile.svelte';
	import Card from '$lib/components/Card.svelte';
	import LoadingScreen from '$lib/components/LoadingScreen.svelte';
	import PageLoader from '$lib/components/PageLoader.svelte';
	import PageWrapper from '$lib/components/PageWrapper.svelte';
	import CheckMark from '$lib/icons/CheckMark.svelte';
	import type { Project } from '$lib/services/ProjectService';
	import ProjectService from '$lib/services/ProjectService';
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
	let projects: Project[] = [];

	$: getProjects();
	async function getProjects() {
		loading = true;
		try {
			let response = await ProjectService.listProjects();
			if (response.length == 0) {
				goto('/create-project');
			}
			projects = response;
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
						<!-- <TableBodyRow>
						<TableBodyCell>A La Carte Menu</TableBodyCell>
						<TableBodyCell
							><a
								class="text-blue-600 hover:underline"
								target="_blank"
								href="https://menulinkup.com/pizza-place/a-la-carte-menu"
								>https://menulinkup.com/pizza-place/a-la-carte-menu</a
							></TableBodyCell
						>
						<TableBodyCell><CheckMark class="fill-current text-green-500 w-7" /></TableBodyCell>
						<TableBodyCell><Button>Edit</Button></TableBodyCell>
					</TableBodyRow> -->
					</TableBody>
				</Table>
				{#if projects.length > 0}
					<div class="mt-4">
						<AddNewFile large project={projects[0]} />
					</div>
				{/if}
				<!-- <pre>{JSON.stringify(projects, undefined, 2)}</pre> -->
			</PageLoader>
		</Card>
	</div>
</PageWrapper>
