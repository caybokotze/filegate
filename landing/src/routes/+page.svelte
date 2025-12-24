<script lang="ts">
	const features = [
		{ icon: 'link', title: 'Public URLs Without Port Forwarding', description: 'Share files instantly with anyone via a generated URL like brave-tiger.filegate.app. No router configuration needed.' },
		{ icon: 'folder', title: 'WebDAV Protocol', description: 'Native file manager integration on Windows, macOS, and Linux. Mount as a network drive.' },
		{ icon: 'play', title: 'DLNA Media Server', description: 'Stream media files directly to Smart TVs, game consoles, and media players.' },
		{ icon: 'lock', title: 'Built-in Authentication', description: 'Secure your shares with automatic or custom username/password protection.' },
		{ icon: 'devices', title: 'Cross-Platform', description: 'Works on Windows, macOS, and Linux. Single binary with zero dependencies.' },
		{ icon: 'home', title: 'Local Network Mode', description: 'Run entirely on your LAN without any external connections.' }
	];
	const steps = [
		{ num: 1, title: 'Run Filegate', description: 'Execute filegate in any directory', code: 'filegate' },
		{ num: 2, title: 'Get Your URL', description: 'Receive a unique public URL and credentials', code: 'URL: https://brave-tiger.filegate.app' },
		{ num: 3, title: 'Share and Access', description: 'Share the URL or mount as network drive', code: '' }
	];
</script>

<div class="min-h-screen bg-gradient-to-b from-slate-900 via-slate-800 to-slate-900 text-white">
	<nav class="fixed top-0 left-0 right-0 z-50 bg-slate-900/80 backdrop-blur-md border-b border-slate-700/50">
		<div class="max-w-6xl mx-auto px-6 py-4 flex items-center justify-between">
			<div class="flex items-center gap-2">
				<div class="w-8 h-8 bg-gradient-to-br from-emerald-400 to-cyan-500 rounded-lg flex items-center justify-center font-bold text-slate-900">F</div>
				<span class="text-xl font-semibold">filegate</span>
			</div>
			<div class="flex items-center gap-6">
				<a href="#features" class="text-slate-300 hover:text-white transition-colors">Features</a>
				<a href="#how-it-works" class="text-slate-300 hover:text-white transition-colors">How It Works</a>
				<a href="#install" class="text-slate-300 hover:text-white transition-colors">Install</a>
				<a href="https://github.com/caybokotze/filegate" target="_blank" rel="noopener noreferrer" class="text-slate-300 hover:text-white transition-colors">GitHub</a>
			</div>
		</div>
	</nav>

	<section class="pt-32 pb-20 px-6">
		<div class="max-w-4xl mx-auto text-center">
			<div class="inline-flex items-center gap-2 bg-emerald-500/10 border border-emerald-500/20 rounded-full px-4 py-2 mb-8">
				<span class="w-2 h-2 bg-emerald-400 rounded-full animate-pulse"></span>
				<span class="text-emerald-400 text-sm font-medium">Open Source CLI Tool</span>
			</div>
			<h1 class="text-5xl md:text-6xl font-bold mb-6 bg-gradient-to-r from-white via-slate-200 to-slate-400 bg-clip-text text-transparent leading-tight">Instant File Sharing<br />From Your Terminal</h1>
			<p class="text-xl text-slate-400 mb-10 max-w-2xl mx-auto leading-relaxed">Expose any directory over the network with a single command. Get a public URL instantly - no port forwarding, no configuration, no hassle.</p>
			<div class="flex flex-col sm:flex-row items-center justify-center gap-4 mb-12">
				<a href="https://github.com/caybokotze/filegate/releases" target="_blank" rel="noopener noreferrer" class="group bg-gradient-to-r from-emerald-500 to-cyan-500 hover:from-emerald-400 hover:to-cyan-400 text-slate-900 font-semibold px-8 py-4 rounded-xl transition-all duration-300 shadow-lg shadow-emerald-500/25 hover:shadow-emerald-500/40 flex items-center gap-2">
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" /></svg>
					Download for Free
					<svg class="w-4 h-4 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" /></svg>
				</a>
				<a href="https://github.com/caybokotze/filegate" target="_blank" rel="noopener noreferrer" class="border border-slate-600 hover:border-slate-500 text-slate-300 hover:text-white font-semibold px-8 py-4 rounded-xl transition-all duration-300 flex items-center gap-2">
					<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
					View on GitHub
				</a>
			</div>
			<div class="bg-slate-800/50 border border-slate-700/50 rounded-2xl overflow-hidden shadow-2xl max-w-2xl mx-auto">
				<div class="flex items-center gap-2 px-4 py-3 bg-slate-800/80 border-b border-slate-700/50">
					<div class="w-3 h-3 rounded-full bg-red-500/80"></div>
					<div class="w-3 h-3 rounded-full bg-yellow-500/80"></div>
					<div class="w-3 h-3 rounded-full bg-green-500/80"></div>
					<span class="ml-2 text-slate-500 text-sm">terminal</span>
				</div>
				<div class="p-6 font-mono text-sm text-left">
					<div class="text-slate-400">$ <span class="text-emerald-400">filegate</span></div>
					<div class="mt-3 text-slate-300"><span class="text-cyan-400">INFO</span> Starting filegate v1.0.0</div>
					<div class="text-slate-300"><span class="text-cyan-400">INFO</span> Connected to relay server</div>
					<div class="mt-3 text-slate-300"><span class="text-emerald-400">URL:</span> https://brave-tiger.filegate.app</div>
					<div class="text-slate-300"><span class="text-emerald-400">User:</span> admin</div>
					<div class="text-slate-300"><span class="text-emerald-400">Pass:</span> aB3kL9mN2pQr</div>
					<div class="mt-3 text-slate-500">Press Ctrl+C to stop</div>
				</div>
			</div>
		</div>
	</section>

	<section id="features" class="py-20 px-6 border-t border-slate-700/50">
		<div class="max-w-6xl mx-auto">
			<div class="text-center mb-16">
				<h2 class="text-3xl md:text-4xl font-bold mb-4">Features</h2>
				<p class="text-slate-400 text-lg max-w-2xl mx-auto">A powerful yet simple tool for sharing files across networks</p>
			</div>
			<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each features as feature}
					<div class="group bg-slate-800/30 hover:bg-slate-800/50 border border-slate-700/50 hover:border-slate-600/50 rounded-2xl p-6 transition-all duration-300 text-center">
						<div class="flex justify-center mb-4">
						{#if feature.icon === 'link'}
							<svg class="w-8 h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" /></svg>
						{:else if feature.icon === 'folder'}
							<svg class="w-8 h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" /></svg>
						{:else if feature.icon === 'play'}
							<svg class="w-8 h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
						{:else if feature.icon === 'lock'}
							<svg class="w-8 h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" /></svg>
						{:else if feature.icon === 'devices'}
							<svg class="w-8 h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" /></svg>
						{:else if feature.icon === 'home'}
							<svg class="w-8 h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" /></svg>
						{/if}
					</div>
						<h3 class="text-xl font-semibold mb-2 text-white">{feature.title}</h3>
						<p class="text-slate-400 leading-relaxed">{feature.description}</p>
					</div>
				{/each}
			</div>
		</div>
	</section>

	<section id="how-it-works" class="py-20 px-6 border-t border-slate-700/50 bg-slate-800/20">
		<div class="max-w-4xl mx-auto">
			<div class="text-center mb-16">
				<h2 class="text-3xl md:text-4xl font-bold mb-4">How It Works</h2>
				<p class="text-slate-400 text-lg max-w-2xl mx-auto">Get started in seconds with just one command</p>
			</div>
			<div class="space-y-8">
				{#each steps as step}
					<div class="flex gap-6 items-start">
						<div class="flex-shrink-0 w-12 h-12 bg-gradient-to-br from-emerald-500 to-cyan-500 rounded-xl flex items-center justify-center text-xl font-bold text-slate-900">{step.num}</div>
						<div class="flex-1">
							<h3 class="text-xl font-semibold mb-2">{step.title}</h3>
							<p class="text-slate-400 mb-4">{step.description}</p>
							{#if step.code}
								<div class="bg-slate-900/80 border border-slate-700/50 rounded-xl p-4 font-mono text-sm">
									<pre class="text-emerald-400 whitespace-pre-wrap">{step.code}</pre>
								</div>
							{/if}
						</div>
					</div>
				{/each}
			</div>
			<div class="mt-16 bg-slate-800/30 border border-slate-700/50 rounded-2xl p-8">
				<h3 class="text-xl font-semibold mb-6 text-center">Architecture</h3>
				<div class="flex flex-col md:flex-row items-center justify-center gap-4 text-center">
					<div class="bg-slate-700/50 rounded-xl p-4 min-w-[140px]">
						<div class="text-2xl mb-2">🌐</div>
						<div class="font-medium">Browser/Client</div>
						<div class="text-slate-400 text-sm">HTTPS</div>
					</div>
					<div class="text-slate-500 text-2xl hidden md:block">→</div>
					<div class="text-slate-500 text-2xl md:hidden">↓</div>
					<div class="bg-emerald-500/20 border border-emerald-500/30 rounded-xl p-4 min-w-[140px]">
						<div class="text-2xl mb-2">☁️</div>
						<div class="font-medium">Relay Server</div>
						<div class="text-slate-400 text-sm">filegate.app</div>
					</div>
					<div class="text-slate-500 text-2xl hidden md:block">→</div>
					<div class="text-slate-500 text-2xl md:hidden">↓</div>
					<div class="bg-slate-700/50 rounded-xl p-4 min-w-[140px]">
						<div class="text-2xl mb-2">💻</div>
						<div class="font-medium">Your Machine</div>
						<div class="text-slate-400 text-sm">WebSocket</div>
					</div>
				</div>
				<p class="text-slate-400 text-sm text-center mt-6">Requests are tunneled through our relay server to your local machine - no port forwarding required</p>
			</div>
		</div>
	</section>

	<section id="install" class="py-20 px-6 border-t border-slate-700/50">
		<div class="max-w-4xl mx-auto">
			<div class="text-center mb-16">
				<h2 class="text-3xl md:text-4xl font-bold mb-4">Installation</h2>
				<p class="text-slate-400 text-lg max-w-2xl mx-auto">Choose your platform and get started</p>
			</div>
			<div class="grid md:grid-cols-3 gap-6 mb-12">
				<div class="bg-slate-800/30 border border-slate-700/50 rounded-2xl p-6 text-center">
					<div class="flex justify-center mb-4">
						<svg class="w-7 h-7 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5"><rect x="3" y="3" width="8" height="8" rx="1" /><rect x="13" y="3" width="8" height="8" rx="1" /><rect x="3" y="13" width="8" height="8" rx="1" /><rect x="13" y="13" width="8" height="8" rx="1" /></svg>
					</div>
					<h3 class="font-semibold mb-2">Windows</h3>
					<p class="text-slate-400 text-sm">Download the .exe from releases</p>
				</div>
				<div class="bg-slate-800/30 border border-slate-700/50 rounded-2xl p-6 text-center">
					<div class="flex justify-center mb-4">
						<svg class="w-7 h-7 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5"><circle cx="12" cy="12" r="9" /><circle cx="12" cy="12" r="3" /></svg>
					</div>
					<h3 class="font-semibold mb-2">macOS</h3>
					<p class="text-slate-400 text-sm">Download the darwin binary</p>
				</div>
				<div class="bg-slate-800/30 border border-slate-700/50 rounded-2xl p-6 text-center">
					<div class="flex justify-center mb-4">
						<svg class="w-7 h-7 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5"><path d="M12 3L4 9v12h16V9l-8-6z" /><path d="M9 21v-6h6v6" /></svg>
					</div>
					<h3 class="font-semibold mb-2">Linux</h3>
					<p class="text-slate-400 text-sm">Download or build from source</p>
				</div>
			</div>
			<div class="space-y-6">
				<h3 class="text-xl font-semibold text-center mb-8">Usage Examples</h3>
				<div class="bg-slate-800/50 border border-slate-700/50 rounded-xl overflow-hidden">
					<div class="px-4 py-2 bg-slate-800/80 border-b border-slate-700/50 text-slate-400 text-sm">WebDAV (Remote - Public URL)</div>
					<div class="p-4 font-mono text-sm"><span class="text-slate-400">$</span> <span class="text-emerald-400">filegate</span></div>
				</div>
				<div class="bg-slate-800/50 border border-slate-700/50 rounded-xl overflow-hidden">
					<div class="px-4 py-2 bg-slate-800/80 border-b border-slate-700/50 text-slate-400 text-sm">WebDAV (Local Network Only)</div>
					<div class="p-4 font-mono text-sm"><span class="text-slate-400">$</span> <span class="text-emerald-400">filegate webdav</span> <span class="text-cyan-400">--local --port</span> <span class="text-amber-400">8080</span></div>
				</div>
				<div class="bg-slate-800/50 border border-slate-700/50 rounded-xl overflow-hidden">
					<div class="px-4 py-2 bg-slate-800/80 border-b border-slate-700/50 text-slate-400 text-sm">DLNA Media Server</div>
					<div class="p-4 font-mono text-sm"><span class="text-slate-400">$</span> <span class="text-emerald-400">filegate dlna</span> <span class="text-cyan-400">--name</span> <span class="text-amber-400">"My Media Server"</span></div>
				</div>
			</div>
			<div class="mt-12 overflow-x-auto">
				<h3 class="text-xl font-semibold mb-6">CLI Options</h3>
				<table class="w-full text-left">
					<thead>
						<tr class="border-b border-slate-700/50">
							<th class="py-3 px-4 text-slate-300 font-medium">Flag</th>
							<th class="py-3 px-4 text-slate-300 font-medium">Description</th>
							<th class="py-3 px-4 text-slate-300 font-medium">Default</th>
						</tr>
					</thead>
					<tbody class="text-slate-400">
						<tr class="border-b border-slate-700/30"><td class="py-3 px-4 font-mono text-cyan-400">--local, -l</td><td class="py-3 px-4">Run in local mode (LAN only)</td><td class="py-3 px-4">false</td></tr>
						<tr class="border-b border-slate-700/30"><td class="py-3 px-4 font-mono text-cyan-400">--port, -p</td><td class="py-3 px-4">Port to listen on</td><td class="py-3 px-4">8080</td></tr>
						<tr class="border-b border-slate-700/30"><td class="py-3 px-4 font-mono text-cyan-400">--user, -u</td><td class="py-3 px-4">Username for Basic Auth</td><td class="py-3 px-4">admin</td></tr>
						<tr class="border-b border-slate-700/30"><td class="py-3 px-4 font-mono text-cyan-400">--pass</td><td class="py-3 px-4">Password for Basic Auth</td><td class="py-3 px-4">auto-generated</td></tr>
						<tr><td class="py-3 px-4 font-mono text-cyan-400">--relay</td><td class="py-3 px-4">Custom relay server URL</td><td class="py-3 px-4">wss://filegate.app/tunnel</td></tr>
					</tbody>
				</table>
			</div>
			<div class="mt-16 text-center">
				<a href="https://github.com/caybokotze/filegate/releases" target="_blank" rel="noopener noreferrer" class="inline-flex items-center gap-2 bg-gradient-to-r from-emerald-500 to-cyan-500 hover:from-emerald-400 hover:to-cyan-400 text-slate-900 font-semibold px-8 py-4 rounded-xl transition-all duration-300 shadow-lg shadow-emerald-500/25 hover:shadow-emerald-500/40">
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" /></svg>
					Download Latest Release
				</a>
			</div>
		</div>
	</section>

	<!-- Support Section -->
	<section class="py-16 px-6 border-t border-slate-700/50 bg-slate-800/20">
		<div class="max-w-2xl mx-auto text-center">
			<div class="flex justify-center mb-6">
				<svg class="w-10 h-10 text-amber-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
					<path d="M4 4h13a3 3 0 013 3v1a3 3 0 01-3 3h-1v5a4 4 0 01-4 4H8a4 4 0 01-4-4V4z" />
					<path d="M4 6h13" />
				</svg>
			</div>
			<h2 class="text-2xl font-bold mb-4">Support the Project</h2>
			<p class="text-slate-400 mb-6 leading-relaxed">
				Filegate is completely free and open source under the MIT License.
				If you find it useful and would like to support its development, consider buying me a coffee.
			</p>
			<a
				href="https://buymeacoffee.com/caybokotze"
				target="_blank"
				rel="noopener noreferrer"
				class="inline-flex items-center gap-2 bg-amber-500 hover:bg-amber-400 text-slate-900 font-semibold px-6 py-3 rounded-xl transition-all duration-300"
			>
				<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
					<path d="M4 4h13a3 3 0 013 3v1a3 3 0 01-3 3h-1v5a4 4 0 01-4 4H8a4 4 0 01-4-4V4zm2 2v10a2 2 0 002 2h4a2 2 0 002-2V6H6zm11 0v4h1a1 1 0 001-1V7a1 1 0 00-1-1h-1z"/>
				</svg>
				Buy me a coffee
			</a>
		</div>
	</section>

	<footer class="py-12 px-6 border-t border-slate-700/50 bg-slate-900/50">
		<div class="max-w-6xl mx-auto">
			<div class="flex flex-col md:flex-row items-center justify-between gap-6">
				<div class="flex items-center gap-2">
					<div class="w-8 h-8 bg-gradient-to-br from-emerald-400 to-cyan-500 rounded-lg flex items-center justify-center font-bold text-slate-900">F</div>
					<span class="text-xl font-semibold">filegate</span>
				</div>
				<div class="flex items-center gap-6 text-slate-400">
					<a href="https://github.com/caybokotze/filegate" target="_blank" rel="noopener noreferrer" class="hover:text-white transition-colors">GitHub</a>
					<a href="https://github.com/caybokotze/filegate/releases" target="_blank" rel="noopener noreferrer" class="hover:text-white transition-colors">Releases</a>
					<a href="https://github.com/caybokotze/filegate/issues" target="_blank" rel="noopener noreferrer" class="hover:text-white transition-colors">Issues</a>
				</div>
				<div class="text-slate-500 text-sm">Open source under MIT License</div>
			</div>
		</div>
	</footer>
</div>
