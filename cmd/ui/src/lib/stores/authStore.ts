import type { DbUser } from '$lib/services/UserService';
import type { User } from 'firebase/auth';
import { writable } from 'svelte/store';

export const authStore = writable<{
	isLoggedIn: boolean;
	user?: User | null;
	dbUser?: DbUser | null;
	initialised: boolean;
}>({
	isLoggedIn: false,
	initialised: false
});
