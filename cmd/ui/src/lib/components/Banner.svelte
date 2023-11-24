<script lang="ts">
	import { classnames } from '$lib/util';

	export let bg: 'purple-to-blue' = 'purple-to-blue';
	export let text: 'light' | 'dark' = 'light';
	export let href: string | undefined = undefined;

	let styles = {
		'purple-to-blue': 'bg-gradient-to-r from-purple-400 to-blue-600'
	};

	let bannerBg = classnames(styles[bg]);
</script>

<div class={classnames('w-full flex flex-row justify-center py-2 px-2 text-center', bannerBg)}>
	<a {href} class={classnames(text == 'light' ? 'text-white' : 'text-zinc-900')}><slot /></a>
</div>
