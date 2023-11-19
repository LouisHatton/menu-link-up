<script>
	import { classnames } from '$lib/util';
</script>

<div
	class={classnames(
		'bg-white dark:bg-gray-800 rounded-lg w-full px-6 py-8 shadow-md',
		$$props.class
	)}
>
	<slot />
</div>
