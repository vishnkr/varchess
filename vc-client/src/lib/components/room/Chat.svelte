<script>
	let messages = [
		{ username: 'John', text: 'Hello there!' },
		{ username: 'Alice', text: 'Hi John! How are you?' },
		{ username: 'John', text: "I'm doing great, thanks! How about you?" },
		{ username: 'Alice', text: "I'm good too. Just enjoying the day." },
		{ username: 'John', text: "That's nice to hear!" },
		{ username: 'Alice', text: 'Indeed. So, any plans for the weekend?' },
		{ username: 'John', text: 'Not yet. Maybe catch a movie or relax at home.' },

		{ username: 'John', text: 'Not yet. Maybe catch a movie or relax at home.' }
	];
</script>

<div class="m-2">
	<div class="max-h-[500px] overflow-y-auto">
		<div class="max-w-md mx-auto bg-white rounded-md shadow-md">
			<div class="p-4">
				<p class="text-center font-semibold text-lg">Chat</p>
				<div class="mt-4">
					{#each messages as message}
						<div class="flex items-start space-x-2">
							<p class="font-semibold">{message.username}:</p>
							<p>{message.text}</p>
						</div>
					{/each}
				</div>
			</div>
		</div>
	</div>
	<div class="p-4">
		<input
			type="text"
			placeholder="Enter your message..."
			class="border border-gray-300 rounded-md px-3 py-2 w-full"
		/>
		<button class="bg-blue-500 hover:bg-blue-600 text-white rounded-md px-4 py-2 mt-2">Send</button>
	</div>
</div>
