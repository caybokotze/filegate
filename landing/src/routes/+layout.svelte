<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>Filegate - Instant File Sharing Over Network</title>
	<meta name="description" content="Share files instantly from your terminal. No port forwarding, no complex setup. WebDAV and DLNA support." />
</svelte:head>

{@render children()}
