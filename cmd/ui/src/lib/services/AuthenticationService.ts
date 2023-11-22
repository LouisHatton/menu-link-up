import { env } from '$env/dynamic/public';
import { authStore } from '$lib/stores/authStore';
import { initializeApp, type FirebaseApp, type FirebaseOptions } from 'firebase/app';
import {
	getAuth,
	type Auth,
	signInWithEmailAndPassword,
	type UserCredential,
	type User,
	createUserWithEmailAndPassword,
	GoogleAuthProvider,
	signInWithRedirect,
	sendEmailVerification
} from 'firebase/auth';

const firebaseConfig: FirebaseOptions = {
	apiKey: env['PUBLIC_API_KEY'],
	authDomain: env['PUBLIC_AUTH_DOMAIN']
};

class AuthenticationService {
	app: FirebaseApp;
	auth: Auth;
	googleProvider: GoogleAuthProvider;
	public currentUser: User | null = null;

	constructor() {
		this.app = initializeApp(firebaseConfig);
		this.auth = getAuth(this.app);
		this.googleProvider = new GoogleAuthProvider();

		this.auth.onAuthStateChanged((user) => {
			this.currentUser = user;
			authStore.set({
				isLoggedIn: user != null,
				user,
				initialised: true
			});
		});
	}

	async signInWithPassword(email: string, password: string): Promise<UserCredential> {
		return signInWithEmailAndPassword(this.auth, email, password);
	}

	async registerWithPassword(email: string, password: string): Promise<UserCredential> {
		return createUserWithEmailAndPassword(this.auth, email, password);
	}

	signInWithGoogle() {
		return signInWithRedirect(this.auth, this.googleProvider);
	}

	async isLoggedIn() {
		return this.auth.currentUser != null;
	}

	async logOut() {
		localStorage.removeItem('sent-verify-email');
		return this.auth.signOut();
	}

	async getToken() {
		return this.auth.currentUser?.getIdToken();
	}

	async reloadUser() {
		await this.auth.currentUser?.reload();
		authStore.set({
			isLoggedIn: this.auth.currentUser != null,
			user: this.auth.currentUser,
			initialised: true
		});
	}

	async sendVerifyEmail(force: boolean) {
		if (!this.currentUser) throw new Error('not logged in!');

		if (!force) {
			let sent = localStorage.getItem('sent-verify-email');
			if (sent != null) {
				return;
			}
		}

		localStorage.setItem('sent-verify-email', 'true');
		await sendEmailVerification(this.currentUser);
	}
}

export default new AuthenticationService();
