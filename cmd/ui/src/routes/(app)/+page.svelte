<script lang="ts">
	import Card from '$lib/components/Card.svelte';
	import DataBox from '$lib/components/DataBox.svelte';
	import PageWrapper from '$lib/components/PageWrapper.svelte';
	import CheckMark from '$lib/icons/CheckMark.svelte';
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
</script>

<PageWrapper>
	<h2 class="text-4xl font-semibold">Welcome back {$authStore.user?.displayName ?? ''}!</h2>
	<div class="mt-14">
		<Card>
			<h3 class="text-2xl font-semibold mb-4">
				Your Menus: <span class="text-sm text-gray-500">(1/3)</span>
			</h3>
			<Table>
				<TableHead>
					<TableHeadCell>Name</TableHeadCell>
					<TableHeadCell>Url</TableHeadCell>
					<TableHeadCell>Status</TableHeadCell>
					<TableHeadCell>Edit</TableHeadCell>
					<TableHeadCell />
				</TableHead>
				<TableBody>
					<TableBodyRow>
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
					</TableBodyRow>
				</TableBody>
			</Table>
			<div class="mt-2">
				<Button>Add New</Button>
			</div>
		</Card>
	</div>
</PageWrapper>
