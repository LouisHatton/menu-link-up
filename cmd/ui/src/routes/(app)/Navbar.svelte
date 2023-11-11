<script>
	import BarChart from '$lib/icons/BarChart.svelte';
	import DashboardIcon from '$lib/icons/DashboardIcon.svelte';
	import LifeBuoy from '$lib/icons/LifeBuoy.svelte';
	import LogOut from '$lib/icons/LogOut.svelte';
	import Megaphone from '$lib/icons/Megaphone.svelte';
	import Moon from '$lib/icons/Moon.svelte';
	import Settings from '$lib/icons/Settings.svelte';
	import Speedometer from '$lib/icons/Speedometer.svelte';
	import Sun from '$lib/icons/Sun.svelte';
	import AuthenticationService from '$lib/services/AuthenticationService';
	import { darkMode } from '$lib/stores/darkMode';
	import SidebarButton from './SidebarButton.svelte';
	import SidebarIcon from './SidebarIcon.svelte';
	import { Button, Dropdown, DropdownItem, Chevron, Avatar } from 'flowbite-svelte';

	function toggleDarkMode() {
		darkMode.toggle();
	}
</script>

<div
	class="border-r border-b-zinc-200 bg-white dark:bg-gray-800 dark:border-r-zinc-600 flex-shrink-0"
>
	<div class="flex flex-row justify-between items-center px-6 h-full">
		<div class="flex flex-row justify-center items-center gap-x-2 my-8">
			<!-- <img
				src="/static/insight-wave-800x800.png"
				alt="logo"
				class="w-10"
				on:keydown={toggleDarkMode}
				on:click={toggleDarkMode}
			/> -->
			<h1 class="text-2xl text-blue-600 font-semibold">MenuLink-Up</h1>
		</div>
		<!-- <div class="my-2 2xl:my-2">
			<div class="flex flex-row gap-x-2 my-2 2xl:gap-y-4 2xl:my-4">
				<SidebarButton href="/"><DashboardIcon filled slot="svg" />Home</SidebarButton>
				<SidebarButton href="/settings"><Settings slot="svg" />Settings</SidebarButton>
			</div>
		</div> -->
		<div class="">
			<div class="flex flex-row justify-center gap-x-6 items-center my-2">
				<SidebarIcon href="/settings">
					<Settings />
				</SidebarIcon>
				<SidebarIcon
					on:click={() => {
						AuthenticationService.logOut();
					}}
				>
					<LogOut />
				</SidebarIcon>
			</div>
		</div>
	</div>
</div>
